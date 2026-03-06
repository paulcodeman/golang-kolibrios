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
- `mtools`
- `qemu-system-x86`

## Build Commands

From the repository root:

```sh
make all
```

Check the focused bootstrap runtime validation:

```sh
make check-runtime
```

Create a new app from the shared template:

```sh
bash ./scripts/new-app.sh demo "KolibriOS Demo"
```

This creates `cmd/demo` with a package-local `Main`, a minimal window loop, and
the shared `mk/kolibri-app.mk` build wiring.

Verify that the template itself still works:

```sh
make check-app-template
```

Run the emulator-backed smoke test:

```sh
make check-emulator-smoke
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
- `cmd/runtimecheck/runtimecheck.kex`
- `cmd/timeprobe/timeprobe.kex`
- `cmd/smokeapp/smokeapp.kex`
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
- `cmd/runtimecheck` - implemented
  - integrated runtime smoke panel
  - strings, slices, interfaces, empty interface equality
  - assertions and type switch in one `.kex`
- `cmd/timeprobe` - implemented
  - system time decode from syscall `3`
  - uptime counters from `26.9` and `26.10`
  - timed wait via function `23`
  - delay probe via function `5`
- `cmd/smokeapp` - implemented
  - headless autorun QEMU smoke for the documented bootstrap subset
  - runtime checks for strings, slices, interfaces, assertions, and timed wait
  - system geometry checks against a temporary pruned copy of the official `kolibri.img`
- `cmd/sysinfo` - implemented, including skin-margin, cursor, keyboard-layout, system-language, skin-switch, and active-window/focus probes
  - kernel version query
  - screen working-area query
  - skin height query
  - keyboard layout table and language probe via `21.2/26.2`
  - system language probe via `21.5/26.5`
  - cursor load/set/delete probe via `37.4/37.5/37.6`
  - default skin apply probes via `48.8/48.13`
  - active-window slot and focus probe via `18.3/18.7`
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
- New applications can be scaffolded from `templates/basic-app` via
  `scripts/new-app.sh`.
- The linker script is generated from `mk/static.lds.in`.
- The linker template emits separate RX and RW load segments, so bootstrap
  builds no longer carry the old RWX `LOAD` warning.
- The current bootstrap runtime subset is documented in `docs/RUNTIME.md`.
- Focused compiler/runtime symbol checks live in `tests/runtime` and run via
  `scripts/check-runtime-probes.sh`.
- Host-side runtime behavior checks live in `tests/runtime/behavior.c` and run
  via `scripts/check-runtime-behavior.sh`.
  On hosts that cannot execute a 32-bit ELF directly, this harness falls back
  to a native-host build while `scripts/check-runtime-probes.sh` still validates
  the bootstrap `gccgo -m32` symbol inventory.
- Template verification lives in `scripts/check-app-template.sh` and confirms
  that `scripts/new-app.sh` generates a buildable sample under `cmd/`.
- Emulator smoke verification lives in `scripts/check-emulator-smoke.sh`; it
  downloads the official `kolibri.img`, prunes non-system payload from a
  temporary copy to free space, replaces the existing `@HA` autorun slot with
  `cmd/smokeapp`, boots QEMU headless, and waits for the smoke app to power the
  guest off after runtime plus system self-checks pass.
