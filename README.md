# golang-kolibrios

[![Build Samples](https://github.com/paulcodeman/golang-kolibrios/actions/workflows/build-example.yml/badge.svg)](https://github.com/paulcodeman/golang-kolibrios/actions/workflows/build-example.yml)

Experimental Go bootstrap for building KolibriOS applications.

This repository currently provides:

- low-level KolibriOS syscall entrypoints in assembly
- Go declarations and small typed wrappers for those syscalls
- a minimal `gccgo`-based runtime glue layer
- working sample applications that build into `.kex` binaries

The project is still in prototype stage. Right now the practical path is
`gccgo` + custom ABI/runtime glue, not native `go build`.

## Current Status

- `cmd/example` builds successfully into `cmd/example/example.kex`
- the build flow targets 32-bit KolibriOS binaries
- the Phase 1 syscall inventory gaps are closed for the documented bootstrap subset
- a longer-term plan is tracked in `ROADMAP.md`

## Repository Layout

- `abi/` - syscall assembly stubs and runtime glue used during linking
- `docs/` - bootstrap and build documentation
- `kos/` - raw Go bindings and small higher-level wrappers
- `mk/` - shared bootstrap make logic and linker templates
- `scripts/` - helper scripts for supported host environments
- `tests/` - focused bootstrap runtime probes
- `ui/` - minimal UI helpers built on top of `kos`
- `cmd/` - sample KolibriOS applications and linker/build files
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

## Build Samples

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

Verify that the shared app template still generates a buildable sample:

```sh
make check-app-template
```

Run the headless emulator smoke check:

```sh
make check-emulator-smoke
```

Output:

- `cmd/example/example.kex`
- `cmd/hello/hello.kex`
- `cmd/strings/strings.kex`
- `cmd/slices/slices.kex`
- `cmd/interfaces/interfaces.kex`
- `cmd/emptyiface/emptyiface.kex`
- `cmd/assertions/assertions.kex`
- `cmd/runtimecheck/runtimecheck.kex`
- `cmd/timeprobe/timeprobe.kex`
- `cmd/sysinfo/sysinfo.kex`
- `cmd/message/message.kex`
- `cmd/ipc/ipc.kex`

The current `Makefile` removes intermediate `.o` and `.gox` files after a
successful build, so only the final `.kex` artifact remains. `make check-runtime`
now runs both the unresolved-symbol inventory probes and a host-side C behavior
harness for the documented bootstrap subset. On hosts that cannot execute a
32-bit ELF directly, the behavior harness falls back to native host execution
while the probe inventory still validates the `gccgo -m32` symbol path.
New applications can be scaffolded from `templates/basic-app` via
`scripts/new-app.sh`.
The first emulator-backed smoke path is available through
`scripts/check-emulator-smoke.sh`; it boots the official KolibriOS image in
QEMU and expects the autorun smoke app to power the guest off after its
self-checks pass.

For full bootstrap instructions, see `docs/BUILD.md`.
For the current raw syscall coverage map, see `docs/SYSCALLS.md`.
For the current bootstrap runtime contract, see `docs/RUNTIME.md`.

## Example Application

The example app:

- opens a KolibriOS window
- draws a red bar and a guide line
- creates left and right buttons
- moves the bar in response to button events

Main sources:

- `cmd/example/app.go`
- `cmd/example/main.go`

## Sample Matrix

- `cmd/example` - implemented
- `cmd/hello` - implemented
- `cmd/strings` - implemented
- `cmd/slices` - byte-slice growth, copy, and string conversion probe
- `cmd/interfaces` - non-empty interface dispatch and equality probe
- `cmd/emptyiface` - empty interface equality probe for comparable values
- `cmd/assertions` - empty/non-empty interface assertions and type switch probe
- `cmd/runtimecheck` - integrated in-app runtime smoke panel for the current subset
- `cmd/timeprobe` - system time, uptime counter, wait timeout, and sleep probe
- `cmd/sysinfo` - kernel/style/title/skin/cursor/keyboard-layout/system-language/active-window probes
- `cmd/message` - function `72` message injection probe
- `cmd/ipc` - function `60` self-IPC event and buffer probe

## Development Notes

- This is not yet a complete Go port for KolibriOS.
- Runtime support is intentionally minimal and currently tuned for the example
  and small bootstrap programs.
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
