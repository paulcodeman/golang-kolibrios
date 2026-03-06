# golang-kolibrios

[![Build Samples](https://github.com/paulcodeman/golang-kolibrios/actions/workflows/build-example.yml/badge.svg)](https://github.com/paulcodeman/golang-kolibrios/actions/workflows/build-example.yml)

Experimental Go bootstrap for building KolibriOS applications.

This repository currently provides:

- low-level KolibriOS syscall entrypoints in assembly
- Go declarations and small typed wrappers for those syscalls
- a minimal `gccgo`-based runtime glue layer
- a working example application that builds into a `.kex` binary

The project is still in prototype stage. Right now the practical path is
`gccgo` + custom ABI/runtime glue, not native `go build`.

## Current Status

- `cmd/example` builds successfully into `cmd/example/example.kex`
- the build flow targets 32-bit KolibriOS binaries
- syscall bindings are being aligned with the official KolibriOS API spec
- a longer-term plan is tracked in `ROADMAP.md`

## Repository Layout

- `abi/` - syscall assembly stubs and runtime glue used during linking
- `docs/` - bootstrap and build documentation
- `kos/` - raw Go bindings and small higher-level wrappers
- `mk/` - shared bootstrap make logic and linker templates
- `scripts/` - helper scripts for supported host environments
- `ui/` - minimal UI helpers built on top of `kos`
- `cmd/example/` - demo KolibriOS application and linker/build files
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

## Build The Example

From the repository root:

```sh
make all
```

Output:

- `cmd/example/example.kex`
- `cmd/hello/hello.kex`
- `cmd/strings/strings.kex`

The current `Makefile` removes intermediate `.o` and `.gox` files after a
successful build, so only the final `.kex` artifact remains.

For full bootstrap instructions, see `docs/BUILD.md`.

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
