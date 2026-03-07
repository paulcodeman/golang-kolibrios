# Go for KolibriOS Roadmap

## Mission

Make KolibriOS a real build target for Go applications, starting from the
current `gccgo` prototype and ending with a reproducible SDK and, ideally, a
native `go build` target.

## Current Baseline

- `examples/window` already builds into `window.kex`.
- The low-level syscall ABI lives in `abi/syscalls_i386.asm`.
- Runtime glue currently lives in `abi/runtime_gccgo.c`.
- Exported low-level Go declarations live in `kos/raw.go`.
- Higher-level wrappers live in `kos/*.go` and `ui/*.go`.
- The syscall specification is already present in `sysfuncs.txt`.
- The project is still a prototype: small runtime surface, hand-curated
  toolchain, but the documented bootstrap subset now has host-side runtime
  checks and headless QEMU smoke validation.

## Strategic Direction

Recommended path:

1. Stabilize the current `gccgo` path so it can build small real applications.
2. Turn the prototype into a reusable KolibriOS Go SDK.
3. Use that SDK work as the bootstrap path toward a native Go port
   (`GOOS=kolibrios`, `GOARCH=386`).

Why this path:

- `gccgo` is useful for bootstrap and experimentation.
- A truly full Go application target will eventually require control over
  `runtime`, `syscall`, standard library behavior, linker integration, and
  `cmd/go`.
- Staying forever on ad-hoc runtime shims will keep the project in demo mode.

## Non-Negotiable Rules

- `sysfuncs.txt` is the source of truth for syscall numbers, register contracts,
  packed arguments, and return conventions.
- Every new syscall wrapper must be traceable to an entry in `sysfuncs.txt`.
- High-level wrappers must not hide unexplained magic constants if the value can
  be named and tied to the syscall specification.

## Phase 0 - Stabilize the Bootstrap Prototype

Goal: turn the current demo into a repeatable bootstrap environment.

Tasks:

- Pin the supported host environment: Linux or WSL Ubuntu with an exact package
  list.
- Replace per-example build assumptions with shared build scripts or shared
  `make` includes.
- Document the one command required to build `examples/window/window.kex`.
- Add a small example matrix:
  - hello world
  - window app
  - heap and string app

Exit criteria:

- A clean machine can build the example with one documented command.
- The example build does not require manual cleanup or manual object handling.
- Toolchain prerequisites are explicit and versioned.

## Phase 1 - Audit and Expand the Syscall Layer

Goal: make the ABI layer trustworthy and scalable.

Tasks:

- Build a syscall inventory from `sysfuncs.txt`.
- Map every exported function in `abi/syscalls_i386.asm` and `kos/raw.go` back
  to the spec.
- Add missing core wrappers for time, memory, file, process, IPC, and basic
  system services.
- Keep raw syscall bindings separate from friendly Go wrappers.
- Add comments or metadata that point to the exact syscall and subfunction
  source.

Exit criteria:

- Every exported raw binding has a spec reference.
- No wrapper signature is based on guesswork.
- New wrappers follow a consistent pattern.

Current bootstrap status:

- The documented raw syscall inventory in `docs/SYSCALLS.md` satisfies the
  original Phase 1 priority audit gaps for the current bootstrap subset.

## Phase 2 - Build a Minimal Runtime for Real Programs

Goal: replace one-off runtime shims with a documented runtime contract.

Tasks:

- Formalize startup, entrypoint, stack assumptions, and init order.
- Implement the required low-level helpers:
  - `memcpy`
  - `memmove`
  - `memcmp`
  - zeroing
  - allocation hooks
  - panic stubs
  - string helpers
- Inventory the runtime symbols emitted by `gccgo` for the supported language
  subset.
- Define the first supported language envelope:
  - strings
  - slices
  - structs
  - interfaces
  - `defer`
  - `panic` and `recover`
- Explicitly defer goroutines and channels until runtime support is ready.
- Add focused tests for allocation, concatenation, interface dispatch, and nil
  checks.

Exit criteria:

- Small apps using heap allocation, strings, slices, and interfaces link
- Runtime symbol requirements are documented and covered by tests.
- Unsupported language features fail in a known, documented way.

## Phase 3 - Provide a Usable KolibriOS Go SDK

Goal: let application authors write Go code instead of linker archaeology.

Tasks:

- Turn `kos` into a layered SDK:
  - raw syscall layer
  - typed Go-friendly wrapper layer
  - optional `ui` helpers for simple GUI work
- Add packages for:
  - windows
  - events
  - drawing
  - time
  - files
  - memory
  - debug output
- Introduce a reusable top-level build script or template for apps.
- Decide how apps declare entrypoints, resources, linker settings, and output
  layout.

Exit criteria:

- A new app can be created from a template without copying build internals.
- Common app patterns are covered by stable packages.
- Multiple independent sample apps build on the same SDK unchanged.

## Phase 4 - Add Emulator-Backed Verification and CI

Goal: catch regressions before they land in real KolibriOS images.

Tasks:

- Choose and script an emulator-based run path.
- Create smoke tests that boot an app and validate observable behavior:
  - process starts
  - window appears
  - event loop handles input
  - timer or file calls behave as expected
- Run build and smoke tests in CI on a Linux or WSL-compatible environment.

Exit criteria:

- Every push can build the SDK and boot at least one smoke-test app.
- Syscall and runtime regressions are detected automatically.
- Manual desktop testing is no longer required for every change.

Current bootstrap status:

- For the documented `gccgo` bootstrap subset, milestones `M0-M4` are in place.
- Remaining roadmap work now starts at standard-library growth and the native
  `GOOS=kolibrios GOARCH=386` path.

## Phase 5 - Grow Standard Library Compatibility

Goal: move from demos to normal Go application structure.

Tasks:

- Define the first supported standard library surface:
  - `errors`
  - `fmt`
  - `bytes`
  - `strings`
  - `io`
  - `time`
  - `path`
  - `os`
  - `syscall`
- Specify where KolibriOS semantics differ from Unix-like systems.
- Decide how file descriptors, paths, process state, environment, and clocks map
  to KolibriOS behavior.
- Add compatibility samples that use ordinary Go package patterns rather than
  custom SDK-only code.

Exit criteria:

- Real apps can import a documented subset of the Go standard library.
- Non-Unix behavior is specified instead of left implicit.
- The supported package set is tested and versioned.

Current bootstrap status:

- The first bootstrap stdlib-compatible shim is now in place with local
  support for `errors.New`, `errors.Unwrap`, and `errors.Is`.
- `examples/files` is the first compatibility sample that imports `errors`
  through the ordinary Go import path.
- Local support for `path.Clean`, `path.Join`, `path.Dir`, `path.Base`,
  `path.Ext`, `path.Split`, and `path.IsAbs` is now in place.
- `examples/path` is the second compatibility sample, using ordinary `import "path"`.
- Local support for `strings.Contains`, `strings.Cut`, `strings.HasPrefix`,
  `strings.HasSuffix`, `strings.Index`, `strings.Join`, `strings.LastIndex`,
  `strings.TrimPrefix`, and `strings.TrimSuffix` is now in place.
- `examples/strings` is the third compatibility sample, using ordinary
  `import "strings"`.
- Local support for `bytes.Contains`, `bytes.Cut`, `bytes.Equal`,
  `bytes.HasPrefix`, `bytes.HasSuffix`, `bytes.Index`, `bytes.IndexByte`,
  `bytes.Join`, `bytes.TrimPrefix`, and `bytes.TrimSuffix` is now in place.
- `examples/bytes` is the fourth compatibility sample, using ordinary
  `import "bytes"`.
- Local support for `io.Reader`, `io.Writer`, `io.Closer`, `io.ReadWriter`,
  `io.ReadCloser`, `io.WriteCloser`, `io.EOF`, `io.ErrShortWrite`,
  `io.ReadAll`, `io.Copy`, `io.CopyBuffer`, and `io.WriteString` is now in
  place.
- `examples/io` is the fifth compatibility sample, using ordinary
  `import "io"`.
- Local support for `os.Getwd`, `os.Open`, `os.Create`, `os.OpenFile`,
  `os.ReadFile`, `os.WriteFile`, `os.Mkdir`, `os.Remove`, `os.Rename`,
  `os.Pipe`, narrow `os.Err*` sentinels, `os.PathError`, `os.LinkError`,
  and fd-backed `os.Stdin`/`os.Stdout`/`os.Stderr` are now in place.
- `examples/os` is the sixth compatibility sample, using ordinary
  `import "os"`.
- Local support for `fmt.Sprint`, `fmt.Sprintln`, `fmt.Sprintf`,
  `fmt.Fprint`, `fmt.Fprintln`, `fmt.Fprintf`, `fmt.Print`, `fmt.Printf`,
  `fmt.Println`, narrow `%s/%d/%x/%X/%t/%v/%c/%%` formatting, and `fmt.Errorf`
  is now in place.
- `examples/fmt` is the seventh compatibility sample, using ordinary
  `import "fmt"`.
- Local support for `syscall.Errno`, `syscall.Read`, `syscall.Write`,
  `syscall.Pipe`, and `syscall.Pipe2` is now in place through the documented
  `77.10`, `77.11`, and `77.13` contracts.
- `apps/diag` now validates the bootstrap `syscall` pipe layer and `fmt.Print*`
  through pipe-backed `os.Stdout` capture in headless QEMU.
- `apps/diag` is the first fuller utility outside `examples/`, giving the
  project a reusable diagnostics app plus a headless emulator validation path.
- `kos` now has a first bootstrap wrapper for `/sys/lib/console.obj`,
  including DLL export lookup plus narrow `con_init`, `con_write_string`,
  `con_set_title`, `con_getch`, and `con_exit` coverage.
- `examples/console` is the first non-window SDK sample built on that console
  wrapper, and `apps/diag` headless validation now exercises the real
  `CONSOLE.OBJ` init/write/exit path instead of DLL load alone.
- Broader package coverage (`time`)
  remains pending.

## Phase 6 - Port the Native Go Toolchain

Goal: reach a real `go build` story instead of relying on `gccgo` forever.

Tasks:

- Define the target tuple and platform contracts:
  - `GOOS=kolibrios`
  - `GOARCH=386`
- Port the `runtime`, `syscall`, and linker pieces needed by the main Go
  toolchain.
- Teach `cmd/dist`, `cmd/link`, and `cmd/go` about the new target.
- Decide the scheduler model:
  - single-threaded first
  - multitasking and goroutines after runtime stability
- Bring up progressively harder programs:
  - hello world
  - allocations
  - maps
  - interfaces
  - goroutines
  - channels

Exit criteria:

- `GOOS=kolibrios GOARCH=386 go build` works for the supported package set.
- The custom `gccgo` bootstrap path is optional instead of mandatory.
- The target is no longer limited to hand-tuned demos.

## Suggested Milestone Order

- M0: bootstrap build is reproducible
- M1: syscall layer is audited against `sysfuncs.txt`
- M2: minimal runtime supports small real apps
- M3: SDK and app template are usable
- M4: emulator smoke tests and CI are in place
- M5: documented stdlib subset works
- M6: native `go build` target exists

## Definition of Success

The project can be called a real Go target for KolibriOS when:

- a clean environment can build and run apps reproducibly
- syscall bindings are spec-driven, not memory-driven
- the runtime contract is explicit and tested
- app authors use stable packages instead of manual ABI knowledge
- CI boots real artifacts
- the project moves toward or reaches native `go build` support
