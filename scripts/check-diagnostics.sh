#!/usr/bin/env bash

set -euo pipefail

repo_root=$(cd "$(dirname "$0")/.." && pwd)
base_image=$(bash "$repo_root/scripts/download-kolibri-image.sh")
artifact_dir="$repo_root/.cache/qemu-diag-run"
runtime_dir=$(mktemp -d -p /tmp kolibri-qemu-diag-XXXXXX)
image_path="$artifact_dir/kolibri-diag.img"
debug_log_path="$artifact_dir/GODIAG.debug.log"
report_path="$artifact_dir/GODIAG.TXT"
normalized_report_path="$artifact_dir/GODIAG.normalized.txt"
headless_flag_path="$runtime_dir/GODIAG.AUTO"
pidfile="$runtime_dir/qemu.pid"
app_binary="$repo_root/apps/diag/diag.kex"
qemu_binary=${QEMU_SYSTEM_I386:-qemu-system-i386}
timeout_seconds=${KOLIBRI_DIAG_TIMEOUT_SECONDS:-45}

cleanup() {
  if [[ -f "$pidfile" ]]; then
    pid=$(cat "$pidfile" 2>/dev/null || true)
    if [[ -n "${pid:-}" ]] && kill -0 "$pid" 2>/dev/null; then
      if kill -0 "$pid" 2>/dev/null; then
        kill "$pid" 2>/dev/null || true
        wait "$pid" 2>/dev/null || true
      fi
    fi
  fi

  rm -rf "$runtime_dir"
}

trap cleanup EXIT

require_tool() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "missing required tool: $1" >&2
    exit 1
  fi
}

require_tool "$qemu_binary"
require_tool mcopy
require_tool mdel
require_tool mdeltree
require_tool curl
require_tool python3

mkdir -p "$artifact_dir"
rm -f "$image_path" "$debug_log_path" "$report_path" "$normalized_report_path"

make -C "$repo_root/apps/diag" clean all

cp "$base_image" "$image_path"
bash "$repo_root/scripts/prune-kolibri-image.sh" "$image_path"
mdel -i "$image_path" "::/GODIAG.TXT" >/dev/null 2>&1 || true
mdel -i "$image_path" "::/GODIAG.TMP" >/dev/null 2>&1 || true
mdel -i "$image_path" "::/GODIAG.AUTO" >/dev/null 2>&1 || true
printf 'headless diagnostics\n' > "$headless_flag_path"
mcopy -o -i "$image_path" "$headless_flag_path" ::/GODIAG.AUTO
mcopy -o -i "$image_path" "$app_binary" ::/@HA

"$qemu_binary" \
  -daemonize \
  -pidfile "$pidfile" \
  -display none \
  -nic none \
  -debugcon "file:$debug_log_path" \
  -global isa-debugcon.iobase=0x402 \
  -drive "file=$image_path,format=raw,if=floppy,index=0" \
  -boot a

if [[ ! -f "$pidfile" ]]; then
  echo "qemu did not create a pidfile" >&2
  exit 1
fi

pid=$(cat "$pidfile")

deadline=$((SECONDS + timeout_seconds))
while (( SECONDS < deadline )); do
  if ! kill -0 "$pid" 2>/dev/null; then
    break
  fi

  sleep 1
done

if kill -0 "$pid" 2>/dev/null; then
  echo "guest did not power off within ${timeout_seconds}s" >&2
  kill "$pid" 2>/dev/null || true
  wait "$pid" 2>/dev/null || true
  exit 1
fi

extract_debug_report() {
  python3 - "$debug_log_path" "$normalized_report_path" <<'PY'
from pathlib import Path
import sys

log_path = Path(sys.argv[1])
out_path = Path(sys.argv[2])

if not log_path.exists():
    raise SystemExit(1)

data = log_path.read_bytes().replace(b"\r", b"")
begin = b"[[GODIAG-BEGIN]]\n"
end = b"[[GODIAG-END]]\n"
start = data.rfind(begin)
if start < 0:
    raise SystemExit(1)
finish = data.find(end, start + len(begin))
if finish < 0:
    raise SystemExit(1)

report = data[start + len(begin):finish]
out_path.write_bytes(report)
PY
}

if ! extract_debug_report; then
  if ! mcopy -o -i "$image_path" "::/GODIAG.TXT" "$report_path" >/dev/null 2>&1; then
    echo "diagnostic report was not written to the debug log or the disk image" >&2
    exit 1
  fi

  tr -d '\r' < "$report_path" > "$normalized_report_path"
fi

cat "$normalized_report_path"

if ! grep -q '^RESULT: PASS$' "$normalized_report_path"; then
  echo "diagnostic report indicates failure" >&2
  exit 1
fi

echo "report saved to $normalized_report_path"
