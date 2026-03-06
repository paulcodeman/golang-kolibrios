# Build Guide

## Supported Bootstrap Environment

The current bootstrap flow is supported on:

- Ubuntu 24.04
- WSL Ubuntu 24.04 on Windows

The repository is not currently set up for native PowerShell-only builds.

## Toolchain Installation

Use the shared bootstrap script:

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

## Build Commands

From the repository root:

```sh
make all
```

Check the focused bootstrap runtime probes:

```sh
make check-runtime
```

Clean the sample builds:

```sh
make clean
```

Rebuild from scratch:

```sh
make rebuild-all
```

## Current Output

Successful build output:

- `cmd/example/example.kex`
- `cmd/hello/hello.kex`
- `cmd/strings/strings.kex`
- `cmd/slices/slices.kex`
- `cmd/interfaces/interfaces.kex`
- `cmd/emptyiface/emptyiface.kex`
- `cmd/assertions/assertions.kex`
- `cmd/sysinfo/sysinfo.kex`
- `cmd/message/message.kex`
- `cmd/ipc/ipc.kex`

Intermediate `.o`, `.gox`, and generated linker files are deleted after a
successful build.

## Current Sample Matrix

- `cmd/example` - implemented
  - window creation
  - redraw loop
  - button input
  - primitive drawing
- `cmd/hello` - implemented
  - minimal window loop
  - text drawing
- `cmd/strings` - implemented
  - string concatenation
  - string equality
  - button-triggered redraw
- `cmd/slices` - implemented
  - `make([]byte, n)`
  - `append(dst, src...)`
  - `append(dst, b1, b2, ...)`
  - `copy(dst, src)`
  - `[]byte(string)`
  - `string([]byte)`
  - slice indexing and `len`
- `cmd/interfaces` - implemented
  - concrete-to-interface assignment
  - non-empty interface method dispatch
  - interface equality for matching comparable concrete types
- `cmd/emptyiface` - implemented
  - empty interface assignment
  - empty interface equality for matching comparable concrete types
- `cmd/assertions` - implemented
  - empty interface to concrete assertion
  - empty interface comma-ok assertion
  - empty interface to interface assertion
  - non-empty interface to interface assertion
  - empty interface type switch
- `cmd/sysinfo` - implemented
  - kernel version query
  - screen working-area query
  - skin height query
  - caption update via function `71.1`
- `cmd/message` - implemented
  - active-window button injection via function `72`
  - active-window key injection via function `72`
  - key event decoding via function `2`
- `cmd/ipc` - implemented
  - IPC buffer registration via function `60.1`
  - self-targeted IPC send via function `60.2`
  - IPC event `7` handling and buffer drain

## Notes

- The syscall reference for all new bindings is `sysfuncs.txt`.
- Shared bootstrap build logic now lives in `mk/kolibri-app.mk`.
- The linker script is generated from `mk/static.lds.in`.
- The current bootstrap runtime subset is documented in `docs/RUNTIME.md`.
- Focused compiler/runtime symbol checks live in `tests/runtime` and run via
  `scripts/check-runtime-probes.sh`.
