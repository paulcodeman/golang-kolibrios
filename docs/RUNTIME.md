# Runtime Contract

This file documents the current bootstrap runtime surface used by the `gccgo`
KolibriOS path.

The goal is not full Go compatibility yet. The goal is a small, explicit
language subset that can build real `.kex` samples without guessing which
runtime symbols the compiler will emit.

## Current Helper Surface

The current runtime glue in `abi/runtime_gccgo.c` provides these key helpers:

- allocation: `runtime.newobject`, `runtime.makeslice`, `runtime.growslice`
- byte and string helpers: `runtime.concatstrings`, `runtime.strequal`,
  `runtime.slicebytetostring`, `runtime.stringtoslicebyte`
- interface helpers: `runtime.efaceeq`, `runtime.ifaceeq`,
  `runtime.interequal`, `runtime.assertitab`, `runtime.ifaceE2I2`,
  `runtime.ifaceE2T2`, `runtime.ifaceI2I2`, `runtime.nilinterequal`
- memory helpers: `runtime.memmove`, `runtime.memequal`,
  `runtime.memequal8`, `runtime.memequal16`, `runtime.memequal32`,
  `runtime.memequal64`, `runtime.typedmemmove`,
  internal `memcpy`/`memcmp`/zeroing helpers
- bootstrap GC stubs: `runtime.registerGCRoots`, `runtime.gcWriteBarrier`,
  `runtime.writeBarrier`
- panic helpers for bounds failures:
  `runtime.goPanicIndex*`, `runtime.goPanicSlice*`, `runtime.panicmem`,
  `runtime.panicdottype`

These helpers are intentionally minimal. GC is not implemented; the current
barrier and root-registration symbols only satisfy the linker for the supported
bootstrap subset.

## Current GCCGo ABI Notes

The current helper surface is based on local `gccgo -m32` probe builds.

- `runtime.makeslice` is currently treated as
  `void *runtime.makeslice(Type *elem, int len, int cap)`
- `runtime.growslice` uses aggregate-return ABI and is implemented as
  `GoSlice runtime.growslice(Type *elem, void *oldPtr, int oldLen, int oldCap, int newLen)`
- `runtime.typedmemmove` is currently treated as
  `void runtime.typedmemmove(Type *t, void *dst, void *src)`
- `runtime.slicebytetostring` uses aggregate-return ABI and is implemented as
  `GoString runtime.slicebytetostring(void *tmpbuf, const uint8 *ptr, int len)`
- `runtime.stringtoslicebyte` uses aggregate-return ABI and is implemented as
  `GoSlice runtime.stringtoslicebyte(void *tmpbuf, const char *ptr, int len)`
- `runtime.ifaceeq` is currently treated as
  `bool runtime.ifaceeq(IMT *leftTab, void *leftData, IMT *rightTab, void *rightData)`
- `runtime.efaceeq` is currently treated as
  `bool runtime.efaceeq(Type *leftType, void *leftData, Type *rightType, void *rightData)`
- `runtime.interequal` is currently treated as
  `bool runtime.interequal(Interface *left, Interface *right)`
- `runtime.ifaceE2T2` is currently treated as
  `bool runtime.ifaceE2T2(Type *wantType, Type *haveType, void *haveData, void *out)`
- `runtime.assertitab` is currently treated as
  `IMT *runtime.assertitab(Type *targetInterface, Type *sourceType)`
- `runtime.ifaceE2I2` uses aggregate-return ABI and is implemented as
  `InterfaceAssert runtime.ifaceE2I2(Type *targetInterface, Type *sourceType, void *sourceData)`
- `runtime.ifaceI2I2` uses aggregate-return ABI and is implemented as
  `InterfaceAssert runtime.ifaceI2I2(Type *targetInterface, IMT *sourceMethods, void *sourceData)`
- `runtime.nilinterequal` is currently treated as
  `bool runtime.nilinterequal(Eface *left, Eface *right)`
- `memmove` must exist as a plain symbol because `gccgo` may emit direct calls
  to `memmove`, not only to `runtime.memmove`
- `runtime.*..f` equality references are function-descriptor data symbols, not
  direct code labels
- `runtime.panicdottype` is treated as a no-return type assertion failure path
- `runtime.goPanicSlice*` helpers are treated as no-return failure paths
- runtime failure output currently uses function `63.1` from `sysfuncs.txt`
  for debug-board writes and function `-1` for termination

## Supported Language Envelope

Validated by current samples:

- strings
  - equality with `==`
  - concatenation with `+`
- fixed-size arrays
  - declaration with `[N]T`
  - literals with `[...]T{...}`
  - indexing and `len`
  - value assignment/copy
  - equality for comparable element types
- heap allocation
  - `new(T)`
  - compiler-generated object temporaries
- byte slices
  - `make([]byte, n)`
  - indexing and `len`
  - `append(dst, src...)`
  - `append(dst, b1, b2, ...)`
  - `copy(dst, src)`
  - `[]byte(string)`
  - `string([]byte)`
- non-empty interfaces
  - concrete-to-interface assignment
  - method dispatch
  - equality for matching comparable concrete types
  - assertion to a matching interface type
  - comma-ok assertion to a matching interface type
- empty interfaces
  - assignment
  - equality for matching comparable concrete types
  - assertion to a matching concrete type
  - comma-ok assertion to a concrete type
  - assertion to a matching interface type
  - comma-ok assertion to a matching interface type
  - simple type switches over validated concrete cases
- basic Go control flow
  - `if`
  - `for`
  - `switch`
  - methods on structs

Sample coverage:

- `examples/runtime` validates string equality/concatenation, fixed-array
  equality and value-copy, byte-slice growth and conversion, non-empty
  interface dispatch/equality, empty interface equality, assertions, comma-ok
  assertions, and a simple type switch inside one interactive KolibriOS app
- `examples/ipc` validates that a small real app can stay within the current
  runtime envelope while using the syscall/UI layers
- `tests/smokeapp` validates a headless runtime subset under the emulator smoke
  path

Focused runtime check coverage:

- `tests/runtime/strings.go` validates the emitted string concat/equality
  symbol path
- `tests/runtime/arrays.go` validates the emitted fixed-array equality path
- `tests/runtime/slices.go` validates the emitted byte-slice/string conversion
  and growth symbol path
- `tests/runtime/interfaces.go` validates the emitted non-empty interface
  dispatch/equality symbol path
- `tests/runtime/emptyiface.go` validates the emitted empty interface equality
  symbol path
- `tests/runtime/assertions.go` validates the emitted empty-interface assertion
  and comma-ok symbol path
- `tests/runtime/assert_iface.go` validates the emitted empty-interface to
  interface assertion symbol path
- `tests/runtime/iface_to_iface.go` validates the emitted non-empty interface
  assertion symbol path
- `tests/runtime/type_switch.go` validates the emitted type-switch symbol path
- `tests/runtime/gcbarrier.go` validates the emitted heap-allocation plus
  pointer-write barrier symbol path for the current malloc-based runtime
- `tests/runtime/behavior.c` validates host-side runtime behavior for
  allocation, fixed-array equality helpers, byte-slice helpers,
  write-barrier stubs, empty-interface equality, and the validated interface
  assertion/dispatch helpers
- `scripts/check-runtime-probes.sh` compiles these probes and checks that
  `abi/runtime_gccgo.c` exports the required symbol set
- `scripts/check-runtime-behavior.sh` compiles and runs the host-side behavior
  harness against `abi/runtime_gccgo.c`; it prefers a 32-bit host build and
  falls back to native-host execution when the environment cannot run 32-bit
  ELF binaries directly

## Not Yet Supported

These features are not yet a supported part of the bootstrap contract:

- general slice growth beyond the validated bootstrap byte-slice paths
- maps
- general type assertions and type switches beyond the validated concrete and
  interface assertion paths
- general interface conversions beyond the validated assertion and equality
  paths
- `defer`
- `panic` and `recover` as a normal language path
- goroutines
- channels

Some of these may partially compile in isolated cases, but they are not yet
documented or guaranteed.

## Garbage Collector Status

- the current bootstrap runtime does not implement tracing or collection
- heap allocation is backed directly by `malloc`/`realloc`/`free` inside
  `abi/runtime_gccgo.c`
- `runtime.registerGCRoots`, `runtime.gcWriteBarrier`, and
  `runtime.writeBarrier` only provide the symbol surface needed by `gccgo`
- pointer-write paths can link and execute for the validated subset, but no
  heap memory is reclaimed automatically

## Failure Behavior

- bounds-check and bootstrap runtime panic helpers now write a short diagnostic
  to the debug board and terminate the process
- `growslice` currently allocates a fresh backing array and leaves the old one
  untouched, so repeated append-heavy code leaks memory under the current
  no-GC bootstrap runtime
- interface assertions currently allocate fresh IMT tables on demand, so
  assertion-heavy code also leaks memory under the current no-GC bootstrap
  runtime
- unsupported runtime paths should be treated as prototype limitations, not as
  stable behavior

## Next Runtime Targets

The next runtime milestones after this slice/string subset are:

- document the emitted `gccgo` runtime symbol inventory more formally
- add broader interface conversions beyond the validated assertion paths
- make runtime failure reporting richer than the current short debug-board text
- broaden the current host-side behavior checks toward failure paths and
  emulator-backed execution
