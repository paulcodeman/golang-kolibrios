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
  `runtime.interequal`
- memory helpers: `runtime.memmove`, `runtime.memequal`,
  `runtime.memequal8`, `runtime.memequal32`, `runtime.typedmemmove`,
  internal `memcpy`/`memcmp`/zeroing helpers
- bootstrap GC stubs: `runtime.registerGCRoots`, `runtime.gcWriteBarrier`,
  `runtime.writeBarrier`
- panic helpers for bounds failures:
  `runtime.goPanicIndex*`, `runtime.goPanicSlice*`, `runtime.panicmem`

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
- `memmove` must exist as a plain symbol because `gccgo` may emit direct calls
  to `memmove`, not only to `runtime.memmove`
- `runtime.*..f` equality references are function-descriptor data symbols, not
  direct code labels
- `runtime.goPanicSlice*` helpers are treated as no-return failure paths
- runtime failure output currently uses function `63.1` from `sysfuncs.txt`
  for debug-board writes and function `-1` for termination

## Supported Language Envelope

Validated by current samples:

- strings
  - equality with `==`
  - concatenation with `+`
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
- empty interfaces
  - assignment
  - equality for matching comparable concrete types
- basic Go control flow
  - `if`
  - `for`
  - `switch`
  - methods on structs

Sample coverage:

- `cmd/strings` validates string equality and concatenation
- `cmd/slices` validates `make([]byte, n)`, `[]byte(string)`, and
  `string([]byte)`
- `cmd/interfaces` validates non-empty interface assignment, dispatch, and
  equality
- `cmd/emptyiface` validates empty interface assignment and equality for
  matching comparable concrete values
- `cmd/ipc` validates that small real apps can stay within the current runtime
  envelope while using the syscall/UI layers

Focused runtime probe coverage:

- `tests/runtime/strings.go` validates the emitted string concat/equality
  symbol path
- `tests/runtime/slices.go` validates the emitted byte-slice/string conversion
  and growth symbol path
- `tests/runtime/interfaces.go` validates the emitted non-empty interface
  dispatch/equality symbol path
- `tests/runtime/emptyiface.go` validates the emitted empty interface equality
  symbol path
- `scripts/check-runtime-probes.sh` compiles these probes and checks that
  `abi/runtime_gccgo.c` exports the required symbol set

## Not Yet Supported

These features are not yet a supported part of the bootstrap contract:

- general slice growth beyond the validated bootstrap byte-slice paths
- maps
- type assertions and type switches
- general interface conversions beyond the validated empty/non-empty equality
  paths
- `defer`
- `panic` and `recover` as a normal language path
- goroutines
- channels

Some of these may partially compile in isolated cases, but they are not yet
documented or guaranteed.

## Failure Behavior

- bounds-check and bootstrap runtime panic helpers now write a short diagnostic
  to the debug board and terminate the process
- `growslice` currently allocates a fresh backing array and leaves the old one
  untouched, so repeated append-heavy code leaks memory under the current
  no-GC bootstrap runtime
- unsupported runtime paths should be treated as prototype limitations, not as
  stable behavior

## Next Runtime Targets

The next runtime milestones after this slice/string subset are:

- document the emitted `gccgo` runtime symbol inventory more formally
- add empty-interface conversions beyond equality
- make runtime failure reporting richer than the current short debug-board text
- grow the runtime probes beyond symbol inventory into behavior checks
