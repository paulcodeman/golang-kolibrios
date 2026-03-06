#!/usr/bin/env bash

set -euo pipefail

repo_root=$(cd "$(dirname "$0")/.." && pwd)
base_image=$(bash "$repo_root/scripts/download-kolibri-image.sh")
work_dir=$(mktemp -d -p /tmp kolibri-qemu-smoke-XXXXXX)
image_path="$work_dir/kolibri-smoke.img"
pidfile="$work_dir/qemu.pid"
app_binary="$repo_root/cmd/smokeapp/smokeapp.kex"
qemu_binary=${QEMU_SYSTEM_I386:-qemu-system-i386}
timeout_seconds=${KOLIBRI_SMOKE_TIMEOUT_SECONDS:-45}

cleanup() {
  if [[ -f "$pidfile" ]]; then
    pid=$(cat "$pidfile" 2>/dev/null || true)
    if [[ -n "${pid:-}" ]] && kill -0 "$pid" 2>/dev/null; then
      kill "$pid" 2>/dev/null || true
      wait "$pid" 2>/dev/null || true
    fi
  fi
  if [[ ${KOLIBRI_SMOKE_KEEP:-0} != 1 ]]; then
    rm -rf "$work_dir"
  fi
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
require_tool mdeltree
require_tool mdel
require_tool curl

make -C "$repo_root/cmd/smokeapp" clean all

cp "$base_image" "$image_path"
bash "$repo_root/scripts/prune-kolibri-image.sh" "$image_path"
# Reuse an existing autorun slot from the stock image instead of depending on
# path-resolution quirks in AUTORUN.DAT for new entries.
mcopy -o -i "$image_path" "$app_binary" ::/@HA

"$qemu_binary" \
  -daemonize \
  -pidfile "$pidfile" \
  -display none \
  -nic none \
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
    echo "guest powered off after emulator smoke"
    exit 0
  fi

  sleep 1
done

echo "guest did not power off within ${timeout_seconds}s" >&2
if kill -0 "$pid" 2>/dev/null; then
  kill "$pid" 2>/dev/null || true
  wait "$pid" 2>/dev/null || true
fi
exit 1
