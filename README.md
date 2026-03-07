# golang-kolibrios

[![Build Examples](https://github.com/paulcodeman/golang-kolibrios/actions/workflows/build-example.yml/badge.svg)](https://github.com/paulcodeman/golang-kolibrios/actions/workflows/build-example.yml)

Experimental Go bootstrap for building KolibriOS applications.

This repository currently provides:

- low-level KolibriOS syscall entrypoints in assembly
- Go declarations and small typed wrappers for those syscalls
- a minimal `gccgo`-based runtime glue layer
- working example applications that build into `.kex` binaries
- a separate `apps/` area for fuller utilities beyond the public examples

The project is still in prototype stage. Right now the practical path is
`gccgo` + custom ABI/runtime glue, not native `go build`.

## Current Status

- `examples/window` builds successfully into `examples/window/window.kex`
- the build flow targets 32-bit KolibriOS binaries
- the documented `gccgo` bootstrap line now covers `M0-M4`: reproducible build, audited syscall/runtime subset, reusable app template, and headless QEMU smoke
- Phase 5 bootstrap work now includes local `errors`, `path`, `strings`, and `bytes` shims plus compatibility samples that import those packages through ordinary Go import paths
- the shared linker script emits separate RX/RW load segments, so example builds no longer trigger the old RWX warning
- public demos now live under `examples/`, fuller utilities live under `apps/`, and internal smoke/test programs live under `tests/`
- `apps/diag` provides a reusable KolibriOS diagnostics utility plus a headless QEMU check path that prefers debug-console report capture and falls back to `/FD/1/GODIAG.TXT`
- a longer-term plan is tracked in `ROADMAP.md`

## Repository Layout

- `abi/` - syscall assembly stubs and runtime glue used during linking
- `apps/` - fuller KolibriOS utilities built on the same bootstrap SDK
- `docs/` - bootstrap and build documentation
- `examples/` - curated public KolibriOS demo applications
- `kos/` - raw Go bindings and small higher-level wrappers
- `mk/` - shared bootstrap make logic and linker templates
- `scripts/` - helper scripts for supported host environments
- `stdlib/` - bootstrap-compatible stdlib shim sources such as `errors`, `path`, `strings`, and `bytes`
- `tests/` - focused bootstrap runtime probes and internal smoke apps
- `ui/` - minimal UI helpers built on top of `kos`
- `sysfuncs.txt` - KolibriOS system function specification
- `AGENTS.md` - repository instructions for future agent work
- `ROADMAP.md` - staged plan for turning this into a fuller Go target

## KolibriOS API Source Of Truth

When adding or changing Go wrappers for KolibriOS APIs, use `sysfuncs.txt` from
the repository root as the source of truth.

Do not guess:

- syscall numbers
- register conventions
- packed argument formats
- return-value behavior
- subfunction codes

The low-level ABI in `abi/syscalls_i386.asm` and the exported declarations in
`kos/raw.go` should stay aligned with that specification.

## Build Requirements

The current build is intended for Linux or WSL. The supported bootstrap host is
Ubuntu 24.04 or WSL Ubuntu 24.04.

Install the toolchain with the shared script:

```sh
bash ./scripts/install-ubuntu-toolchain.sh
```

This installs:

- `gcc`
- `gccgo`
- `gcc-multilib`
- `gccgo-multilib`
- `make`
- `nasm`
- `binutils`
- `mtools`
- `qemu-system-x86`

## Build Examples

From the repository root:

```sh
make all
```

Focused runtime checks:

```sh
make check-runtime
```

Create a new scaffolded app:

```sh
bash ./scripts/new-app.sh demo "KolibriOS Demo"
```

Verify that the shared app template still generates a buildable example:

```sh
make check-app-template
```

Run the headless emulator smoke check:

```sh
make check-emulator-smoke
```

Run the headless diagnostics utility check:

```sh
make check-diagnostics
```

Output:

- `apps/diag/diag.kex`
- `examples/window/window.kex`
- `examples/runtime/runtime.kex`
- `examples/time/time.kex`
- `examples/system/system.kex`
- `examples/input/input.kex`
- `examples/ipc/ipc.kex`
- `examples/files/files.kex`
- `examples/path/path.kex`
- `examples/strings/strings.kex`
- `examples/bytes/bytes.kex`
- `tests/smokeapp/smokeapp.kex`

The current `Makefile` removes intermediate `.o` and `.gox` files after a
successful build, so only the final `.kex` artifact remains. `make check-runtime`
now runs both the unresolved-symbol inventory probes and a host-side C behavior
harness for the documented bootstrap subset. On hosts that cannot execute a
32-bit ELF directly, the behavior harness falls back to native host execution
while the probe inventory still validates the `gccgo -m32` symbol path.
New applications can be scaffolded from `templates/basic-app` via
`scripts/new-app.sh` into `examples/<name>`.
The shared app makefile now accepts an ordered `PACKAGE_DIRS` list so bootstrap
apps can precompile additional shared packages beyond `kos` and `ui`.
Shim sources for ordinary stdlib imports now live under `stdlib/<name>`, while
their compiled export data is still exposed through the repository-root include
path so existing `import "errors"` / `import "path"` style code keeps working.
The first emulator-backed smoke path is available through
`scripts/check-emulator-smoke.sh`; it boots a pruned temporary copy of the
official KolibriOS image in QEMU, replaces the existing `@HA` autorun slot with
`tests/smokeapp`, and expects the smoke app to power the guest off after runtime
and system self-checks pass.
The diagnostics runner is available through `scripts/check-diagnostics.sh`; it
boots the same pruned image with `apps/diag`, requests headless mode through a
small marker file, captures the report primarily from the QEMU debug console,
and only falls back to `/FD/1/GODIAG.TXT` if debug-console capture is unavailable.

For full bootstrap instructions, see `docs/BUILD.md`.
For the current raw syscall coverage map, see `docs/SYSCALLS.md`.
For the current bootstrap runtime contract, see `docs/RUNTIME.md`.
For the current bootstrap stdlib-compatible package surface, see
`docs/STDLIB.md`.

## Window Example

The window demo:

- opens a KolibriOS window
- draws a red bar and a guide line
- creates left and right buttons
- moves the bar in response to button events

Main sources:

- `examples/window/app.go`
- `examples/window/main.go`

## Example Matrix

- `examples/window` - basic window loop, redraw handling, buttons, and primitive drawing
- `examples/runtime` - integrated runtime smoke panel for strings, fixed arrays, slices, interfaces, empty interfaces, assertions, and type switches
- `examples/time` - system time, uptime counters, wait timeout, and sleep probe
- `examples/system` - kernel/style/title/skin/cursor/keyboard-layout/system-language/active-window probes
- `examples/input` - function `72` button/key injection and input event probe
- `examples/ipc` - function `60` self-IPC event and buffer probe
- `examples/files` - file info/read probe plus ordinary `import "errors"` compatibility sample
- `examples/path` - path normalization and split probe plus ordinary `import "path"` compatibility sample
- `examples/strings` - ordinary `import "strings"` compatibility sample for join, match, cut, index, and trim helpers
- `examples/bytes` - ordinary `import "bytes"` compatibility sample for byte-slice join, match, cut, equality, and trim helpers
- `apps/diag` - fuller diagnostic utility with GUI summary, report export, and headless QEMU diagnostics capture
- `tests/smokeapp` - internal headless QEMU autorun smoke for the runtime and system bootstrap subset

## Development Notes

- This is not yet a complete Go port for KolibriOS.
- Runtime support is intentionally minimal and currently tuned for the example
  and small bootstrap programs.
- The current bootstrap runtime does not implement a real GC yet; heap objects
  are still allocated through a malloc-based shim with barrier/root stubs.
- The project should evolve toward a reusable SDK first, and eventually toward
  native `GOOS=kolibrios GOARCH=386` support.

## Roadmap

See `ROADMAP.md` for the staged plan:

- stabilize the bootstrap toolchain
- audit and expand the syscall layer
- build a better runtime contract
- provide an SDK for apps
- add emulator-backed testing
- move toward a native Go target

## Repository URL

- https://github.com/paulcodeman/golang-kolibrios
