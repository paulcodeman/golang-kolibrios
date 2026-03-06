#!/usr/bin/env bash

set -euo pipefail

repo_root=$(cd "$(dirname "$0")/.." && pwd)
tmp_dir=$(mktemp -d)

cleanup() {
  rm -rf "$tmp_dir"
}

trap cleanup EXIT

gccgo_flags=(
  -m32
  -c
  -nostdlib
  -nostdinc
  -fno-stack-protector
  -fno-split-stack
  -static
  -fno-leading-underscore
  -fno-common
  -fno-pie
  -ffunction-sections
  -fdata-sections
)

gcc_flags=(
  -m32
  -c
  -ffunction-sections
  -fdata-sections
  -fno-pic
  -fno-pie
  -fno-stack-protector
)

compile_runtime() {
  gcc "${gcc_flags[@]}" \
    "$repo_root/abi/runtime_gccgo.c" \
    -o "$tmp_dir/runtime_gccgo.o"
}

runtime_defined_symbols() {
  nm -g --defined-only "$tmp_dir/runtime_gccgo.o" | awk '{print $3}'
}

compile_probe() {
  local probe_name=$1

  gccgo "${gccgo_flags[@]}" \
    -o "$tmp_dir/$probe_name.o" \
    "$repo_root/tests/runtime/$probe_name.go"
}

probe_unresolved_symbols() {
  local probe_name=$1

  nm -u "$tmp_dir/$probe_name.o" | awk '{print $2}' | sort -u
}

require_list_symbols() {
  local list_name=$1
  local list_data=$2
  shift 2
  local symbol

  for symbol in "$@"; do
    if ! grep -Fxq "$symbol" <<<"$list_data"; then
      printf 'missing symbol in %s: %s\n' "$list_name" "$symbol" >&2
      exit 1
    fi
  done
}

main() {
  local runtime_symbols
  local probe_symbols

  compile_runtime
  runtime_symbols=$(runtime_defined_symbols)

  require_list_symbols "runtime_gccgo.o exports" "$runtime_symbols" \
    "memcmp" \
    "memmove" \
    "runtime.concatstrings" \
    "runtime.growslice" \
    "runtime.makeslice" \
    "runtime.slicebytetostring" \
    "runtime.stringtoslicebyte" \
    "runtime.ifaceeq" \
    "runtime.efaceeq" \
    "runtime.interequal" \
    "runtime.interequal..f" \
    "runtime.memequal32..f" \
    "runtime.memequal8..f" \
    "runtime.newobject" \
    "runtime.panicmem" \
    "runtime.strequal..f" \
    "runtime.typedmemmove" \
    "runtime.writeBarrier"

  compile_probe "strings"
  probe_symbols=$(probe_unresolved_symbols "strings")
  require_list_symbols "strings probe" "$probe_symbols" \
    "memcmp" \
    "runtime.concatstrings"

  compile_probe "slices"
  probe_symbols=$(probe_unresolved_symbols "slices")
  require_list_symbols "slices probe" "$probe_symbols" \
    "memmove" \
    "runtime.growslice" \
    "runtime.makeslice" \
    "runtime.memequal32..f" \
    "runtime.memequal8..f" \
    "runtime.slicebytetostring" \
    "runtime.stringtoslicebyte"

  compile_probe "interfaces"
  probe_symbols=$(probe_unresolved_symbols "interfaces")
  require_list_symbols "interfaces probe" "$probe_symbols" \
    "memcmp" \
    "runtime.ifaceeq" \
    "runtime.interequal..f" \
    "runtime.memequal32..f" \
    "runtime.newobject" \
    "runtime.panicmem" \
    "runtime.strequal..f" \
    "runtime.typedmemmove" \
    "runtime.writeBarrier"

  compile_probe "emptyiface"
  probe_symbols=$(probe_unresolved_symbols "emptyiface")
  require_list_symbols "emptyiface probe" "$probe_symbols" \
    "runtime.efaceeq" \
    "runtime.memequal32..f" \
    "runtime.strequal..f"

  printf 'runtime probes passed\n'
}

main "$@"
