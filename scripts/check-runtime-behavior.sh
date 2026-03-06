#!/usr/bin/env bash

set -euo pipefail

repo_root=$(cd "$(dirname "$0")/.." && pwd)
tmp_dir=$(mktemp -d)

cleanup() {
  rm -rf "$tmp_dir"
}

trap cleanup EXIT

gcc_flags=(
  -m32
  -std=c99
  -fno-pic
  -fno-pie
  -fno-stack-protector
)

native_gcc_flags=(
  -std=c99
  -fno-pic
  -fno-pie
  -fno-stack-protector
)

try_compile_and_run() {
  local binary_path=$1
  local stderr_path=$2
  local status=0
  shift 2

  if ! gcc "$@" \
    "$repo_root/abi/runtime_gccgo.c" \
    "$repo_root/tests/runtime/behavior.c" \
    -no-pie \
    -o "$binary_path"; then
    return $?
  fi

  set +e
  "$binary_path" 2>"$stderr_path"
  status=$?
  set -e
  if [[ $status -ne 0 && $status -ne 126 && $status -ne 127 ]]; then
    cat "$stderr_path" >&2
  fi

  return "$status"
}

main() {
  local status=0

  if try_compile_and_run "$tmp_dir/runtime-behavior-32" "$tmp_dir/runtime-behavior-32.stderr" "${gcc_flags[@]}"; then
    return 0
  else
    status=$?
  fi
  if [[ $status -ne 126 && $status -ne 127 ]]; then
    exit "$status"
  fi

  printf 'runtime behavior checks: falling back to native host execution\n'
  try_compile_and_run "$tmp_dir/runtime-behavior-native" "$tmp_dir/runtime-behavior-native.stderr" "${native_gcc_flags[@]}"
}

main "$@"
