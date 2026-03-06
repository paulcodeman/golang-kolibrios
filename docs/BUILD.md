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
make example
```

Clean the example build:

```sh
make clean
```

Rebuild from scratch:

```sh
make rebuild-example
```

## Current Output

Successful build output:

- `cmd/example/example.kex`

Intermediate `.o`, `.gox`, and generated linker files are deleted after a
successful build.

## Current Sample Matrix

- `cmd/example` - implemented
  - window creation
  - redraw loop
  - button input
  - primitive drawing
- `cmd/hello` - planned
- `cmd/strings` - planned

## Notes

- The syscall reference for all new bindings is `sysfuncs.txt`.
- Shared bootstrap build logic now lives in `mk/kolibri-app.mk`.
- The linker script is generated from `mk/static.lds.in`.
