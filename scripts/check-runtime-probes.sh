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

compile_package() {
  local package_name=$1
  shift
  local source_dir
  local include_dirs=("-I$repo_root")
  local dependency

  source_dir=$(package_source_dir "$package_name")

  for dependency in "$@"; do
    compile_package_export "$dependency"
    include_dirs+=("-I$tmp_dir")
  done

  gccgo "${gccgo_flags[@]}" \
    "${include_dirs[@]}" \
    -o "$tmp_dir/$package_name.gccgo.o" \
    "$source_dir"/*.go
}

compile_package_export() {
  local package_name=$1
  local source_dir

  source_dir=$(package_source_dir "$package_name")

  gccgo "${gccgo_flags[@]}" \
    -I"$repo_root" \
    -I"$tmp_dir" \
    -o "$tmp_dir/$package_name.export.o" \
    "$source_dir"/*.go

  objcopy -j .go_export "$tmp_dir/$package_name.export.o" "$tmp_dir/$package_name.gox"
}

package_source_dir() {
  local package_name=$1

  if [[ -d "$repo_root/$package_name" ]]; then
    printf '%s\n' "$repo_root/$package_name"
    return
  fi

  if [[ -d "$repo_root/stdlib/$package_name" ]]; then
    printf '%s\n' "$repo_root/stdlib/$package_name"
    return
  fi

  printf 'package source dir not found: %s\n' "$package_name" >&2
  exit 1
}

probe_unresolved_symbols() {
  local probe_name=$1

  nm -u "$tmp_dir/$probe_name.o" | sed -E 's/^[[:space:]]*U[[:space:]]+//' | sort -u
}

package_unresolved_symbols() {
  local package_name=$1

  nm -u "$tmp_dir/$package_name.gccgo.o" | sed -E 's/^[[:space:]]*U[[:space:]]+//' | sort -u
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
    "runtime.gcWriteBarrier" \
    "runtime.growslice" \
    "runtime.makeslice" \
    "runtime.assertitab" \
    "runtime.ifaceE2I2" \
    "runtime.ifaceE2T2" \
    "runtime.ifaceI2I2" \
    "runtime.memequal" \
    "runtime.slicebytetostring" \
    "runtime.stringtoslicebyte" \
    "runtime.ifaceeq" \
    "runtime.efaceeq" \
    "runtime.interequal" \
    "runtime.interequal..f" \
    "runtime.memequal32..f" \
    "runtime.memequal64..f" \
    "runtime.memequal8..f" \
    "runtime.nilinterequal" \
    "runtime.nilinterequal..f" \
    "runtime.newobject" \
    "runtime.panicdottype" \
    "runtime.goPanicIndex" \
    "runtime.goPanicIndexU" \
    "runtime.goPanicSliceAcap" \
    "runtime.goPanicSliceAlen" \
    "runtime.goPanicSliceB" \
    "runtime.panicmem" \
    "runtime.registerGCRoots" \
    "runtime.strequal..f" \
    "runtime.typedmemmove" \
    "runtime.writeBarrier"

  compile_package "kos"
  probe_symbols=$(package_unresolved_symbols "kos")
  require_list_symbols "kos package" "$probe_symbols" \
    "go_0kos.LoadDLL" \
    "go_0kos.LoadDLLWithEncoding" \
    "memcmp" \
    "runtime.concatstrings" \
    "runtime.gcWriteBarrier" \
    "runtime.growslice" \
    "runtime.makeslice" \
    "runtime.memequal" \
    "runtime_kos_call_stdcall0" \
    "runtime_kos_call_stdcall1" \
    "runtime_kos_call_stdcall1_void" \
    "runtime_kos_call_stdcall2" \
    "runtime_kos_call_stdcall2_void" \
    "runtime_kos_call_stdcall5_void" \
    "runtime_kos_lookup_dll_export" \
    "runtime_alloc_cstring" \
    "runtime_free_cstring" \
    "runtime_pointer_value" \
    "runtime.goPanicIndex" \
    "runtime.goPanicSliceAcap" \
    "runtime.goPanicSliceAcapU" \
    "runtime.goPanicSliceAlen" \
    "runtime.goPanicSliceB" \
    "runtime.memequal32..f" \
    "runtime.newobject" \
    "runtime.panicmem" \
    "runtime.slicebytetostring" \
    "runtime.strequal..f" \
    "runtime.typedmemmove" \
    "runtime.writeBarrier"

  compile_package "path"
  probe_symbols=$(package_unresolved_symbols "path")
  require_list_symbols "path package" "$probe_symbols" \
    "memcmp" \
    "runtime.concatstrings" \
    "runtime.gcWriteBarrier" \
    "runtime.goPanicIndex" \
    "runtime.goPanicSliceAcap" \
    "runtime.goPanicSliceAlen" \
    "runtime.goPanicSliceB" \
    "runtime.growslice" \
    "runtime.memequal32..f" \
    "runtime.strequal..f" \
    "runtime.writeBarrier"

  compile_package "strings"
  probe_symbols=$(package_unresolved_symbols "strings")
  require_list_symbols "strings package" "$probe_symbols" \
    "runtime.concatstrings" \
    "runtime.goPanicIndex" \
    "runtime.goPanicSliceAlen" \
    "runtime.goPanicSliceB"

  compile_package "bytes"
  probe_symbols=$(package_unresolved_symbols "bytes")
  require_list_symbols "bytes package" "$probe_symbols" \
    "memmove" \
    "runtime.goPanicIndex" \
    "runtime.goPanicSliceAcap" \
    "runtime.goPanicSliceB" \
    "runtime.growslice" \
    "runtime.makeslice" \
    "runtime.memequal32..f" \
    "runtime.memequal8..f" \
    "runtime.newobject"

  compile_package "io"
  probe_symbols=$(package_unresolved_symbols "io")
  require_list_symbols "io package" "$probe_symbols" \
    "memcmp" \
    "memmove" \
    "runtime.gcWriteBarrier" \
    "runtime.goPanicSliceAcap" \
    "runtime.goPanicSliceB" \
    "runtime.growslice" \
    "runtime.ifacevaleq" \
    "runtime.interequal..f" \
    "runtime.makeslice" \
    "runtime.memequal32..f" \
    "runtime.memequal8..f" \
    "runtime.newobject" \
    "runtime.registerGCRoots" \
    "runtime.strequal..f" \
    "runtime.stringtoslicebyte" \
    "runtime.typedmemmove" \
    "runtime.writeBarrier"

  compile_package "syscall" "kos"
  probe_symbols=$(package_unresolved_symbols "syscall")
  require_list_symbols "syscall package" "$probe_symbols" \
    "go_0kos.CreatePipe" \
    "go_0kos.FDRead" \
    "go_0kos.FDWrite" \
    "memcmp" \
    "runtime.concatstrings" \
    "runtime.goPanicIndex" \
    "runtime.goPanicIndexU" \
    "runtime.memequal" \
    "runtime.memequal32..f" \
    "runtime.newobject" \
    "runtime.panicmem" \
    "runtime.registerGCRoots" \
    "runtime.strequal..f"

  compile_package "os" "kos" "io" "syscall"
  probe_symbols=$(package_unresolved_symbols "os")
  require_list_symbols "os package" "$probe_symbols" \
    "go_0io.EOF" \
    "go_0io.ErrShortWrite" \
    "go_0io.ioError..p" \
    "go_0io.ioError.Error" \
    "go_0kos.CreateDirectory" \
    "go_0kos.CreateOrRewriteFile" \
    "go_0kos.CurrentFolder" \
    "go_0kos.DeletePath" \
    "go_0kos.FileSystemStatus..d" \
    "go_0kos.GetPathInfo" \
    "go_0kos.ReadAllFile" \
    "go_0kos.ReadFile" \
    "go_0kos.RenamePath" \
    "go_0kos.WriteFile" \
    "go_0syscall.Pipe" \
    "go_0syscall.Read" \
    "go_0syscall.Write" \
    "memcmp" \
    "runtime.concatstrings" \
    "runtime.gcWriteBarrier" \
    "runtime.goPanicIndex" \
    "runtime.goPanicIndexU" \
    "runtime.ifaceeq" \
    "runtime.interequal..f" \
    "runtime.memequal" \
    "runtime.memequal32..f" \
    "runtime.memequal64..f" \
    "runtime.memequal8..f" \
    "runtime.newobject" \
    "runtime.registerGCRoots" \
    "runtime.strequal..f" \
    "runtime.typedmemmove" \
    "runtime.writeBarrier"

  compile_package "fmt" "errors" "io" "kos" "syscall" "os"
  probe_symbols=$(package_unresolved_symbols "fmt")
  require_list_symbols "fmt package" "$probe_symbols" \
    "go_0errors.New" \
    "go_0io.ErrShortWrite" \
    "go_0io.WriteString" \
    "go_0io.Writer..d" \
    "go_0io.ioError..p" \
    "go_0io.ioError.Error" \
    "go_0os.File..p" \
    "go_0os.File.Write" \
    "go_0os.Stdout" \
    "memcmp" \
    "memmove" \
    "runtime.concatstrings" \
    "runtime.gcWriteBarrier" \
    "runtime.goPanicIndex" \
    "runtime.goPanicIndexU" \
    "runtime.goPanicSliceAlen" \
    "runtime.goPanicSliceB" \
    "runtime.growslice" \
    "runtime.ifaceE2I2" \
    "runtime.ifaceeq" \
    "runtime.interequal..f" \
    "runtime.memequal" \
    "runtime.memequal16..f" \
    "runtime.memequal32..f" \
    "runtime.memequal64..f" \
    "runtime.memequal8..f" \
    "runtime.newobject" \
    "runtime.nilinterequal..f" \
    "runtime.panicdottype" \
    "runtime.registerGCRoots" \
    "runtime.slicebytetostring" \
    "runtime.strequal..f" \
    "runtime.stringtoslicebyte" \
    "runtime.typedmemmove" \
    "runtime.writeBarrier"

  compile_probe "arrays"
  probe_symbols=$(probe_unresolved_symbols "arrays")
  require_list_symbols "arrays probe" "$probe_symbols" \
    "memcmp" \
    "runtime.memequal"

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

  compile_probe "assertions"
  probe_symbols=$(probe_unresolved_symbols "assertions")
  require_list_symbols "assertions probe" "$probe_symbols" \
    "runtime.ifaceE2T2" \
    "runtime.memequal32..f" \
    "runtime.nilinterequal..f" \
    "runtime.panicdottype" \
    "runtime.strequal..f"

  compile_probe "assert_iface"
  probe_symbols=$(probe_unresolved_symbols "assert_iface")
  require_list_symbols "assert_iface probe" "$probe_symbols" \
    "memcmp" \
    "runtime.assertitab" \
    "runtime.ifaceE2I2" \
    "runtime.interequal..f" \
    "runtime.memequal32..f" \
    "runtime.panicmem" \
    "runtime.strequal..f"

  compile_probe "iface_to_iface"
  probe_symbols=$(probe_unresolved_symbols "iface_to_iface")
  require_list_symbols "iface_to_iface probe" "$probe_symbols" \
    "runtime.assertitab" \
    "runtime.ifaceI2I2" \
    "runtime.interequal..f" \
    "runtime.memequal32..f" \
    "runtime.strequal..f"

  compile_probe "errors_unwrap"
  probe_symbols=$(probe_unresolved_symbols "errors_unwrap")
  require_list_symbols "errors_unwrap probe" "$probe_symbols" \
    "runtime.ifaceE2I2" \
    "runtime.interequal..f" \
    "runtime.memequal32..f" \
    "runtime.strequal..f"

  compile_probe "errors_is"
  probe_symbols=$(probe_unresolved_symbols "errors_is")
  require_list_symbols "errors_is probe" "$probe_symbols" \
    "runtime.ifaceE2I2" \
    "runtime.ifaceeq" \
    "runtime.interequal..f" \
    "runtime.memequal32..f" \
    "runtime.strequal..f"

  compile_probe "type_switch"
  probe_symbols=$(probe_unresolved_symbols "type_switch")
  require_list_symbols "type_switch probe" "$probe_symbols" \
    "runtime.memequal32..f" \
    "runtime.strequal..f"

  compile_probe "gcbarrier"
  probe_symbols=$(probe_unresolved_symbols "gcbarrier")
  require_list_symbols "gcbarrier probe" "$probe_symbols" \
    "runtime.gcWriteBarrier" \
    "runtime.memequal" \
    "runtime.memequal32..f" \
    "runtime.newobject" \
    "runtime.typedmemmove" \
    "runtime.writeBarrier"

  printf 'runtime probes passed\n'
}

main "$@"
