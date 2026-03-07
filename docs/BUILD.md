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

This creates `examples/demo` with a package-local `Main`, a minimal window
loop, and the shared `mk/kolibri-app.mk` build wiring.

Verify that the template itself still works:

```sh
make check-app-template
```

Run the emulator-backed smoke test:

```sh
make check-emulator-smoke
```

Run the headless diagnostics utility:

```sh
make check-diagnostics
```

Clean the example builds:

```sh
make clean
```

Rebuild from scratch:

```sh
make rebuild-all
```

## Current Output

Successful build output:

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
- `examples/io/io.kex`
- `examples/os/os.kex`
- `examples/fmt/fmt.kex`
- `examples/console/console.kex`
- `tests/smokeapp/smokeapp.kex`

Intermediate `.o`, `.gox`, and generated linker files are deleted after a
successful build.

## Current Example Matrix

- `examples/window` - implemented
  - window creation
  - redraw loop
  - button input
  - primitive drawing
- `examples/runtime` - implemented
  - integrated runtime smoke panel
  - strings, fixed arrays, slices, interfaces, empty interface equality
  - assertions and type switch in one `.kex`
- `examples/time` - implemented
  - bootstrap-compatible `import "time"` sample with `Now`, `Unix`, `Sleep`, and `Since`
  - wall clock assembled from syscalls `29` and `3`
  - monotonic duration path backed by `26.10`
  - timed wait and redraw loop still driven by the KolibriOS event layer
- `examples/system` - implemented, including skin-margin, cursor, keyboard-layout, system-language, skin-switch, and active-window/focus probes
  - kernel version query
  - screen working-area query
  - skin height query
  - keyboard layout table and language probe via `21.2/26.2`
  - system language probe via `21.5/26.5`
  - cursor load/set/delete probe via `37.4/37.5/37.6`
  - default skin apply probes via `48.8/48.13`
  - active-window slot and focus probe via `18.3/18.7`
  - caption update via function `71.1`
- `examples/input` - implemented
  - active-window button injection via function `72`
  - active-window key injection via function `72`
  - key event decoding via function `2`
- `examples/ipc` - implemented
  - IPC buffer registration via function `60.1`
  - self-targeted IPC send via function `60.2`
  - IPC event `7` handling and buffer drain
- `examples/files` - implemented
  - metadata probe via ordinary `os.Stat` with raw `kos.FileInfo` available through `FileInfo.Sys()`
  - file head read via ordinary `os.Open` / `Read` / `Close`
  - bootstrap-compatible `import "errors"`, `import "io"`, and `import "os"` sample with wrapped sentinel checks
- `examples/path` - implemented
  - bootstrap-compatible `import "path"` sample with `Clean`, `Join`, `Dir`, `Base`, `Ext`, and `IsAbs`
  - slash-based path normalization against a real KolibriOS file probe
- `examples/strings` - implemented
  - bootstrap-compatible `import "strings"` sample with `Join`, `Contains`, `HasPrefix`, `HasSuffix`, `Index`, `LastIndex`, `Cut`, `TrimPrefix`, and `TrimSuffix`
  - string helper checks tied to a real KolibriOS file path and current-folder probe
- `examples/bytes` - implemented
  - bootstrap-compatible `import "bytes"` sample with `Join`, `Equal`, `Contains`, `HasPrefix`, `HasSuffix`, `Index`, `IndexByte`, `Cut`, `TrimPrefix`, and `TrimSuffix`
  - byte-slice helper checks tied to a real KolibriOS file path and current-folder probe
- `examples/io` - implemented
  - bootstrap-compatible `import "io"` sample with `Reader`, `Writer`, `ReadAll`, `Copy`, and `WriteString`
  - chunked stream checks tied to a real KolibriOS file path and current-folder probe
- `examples/os` - implemented
  - bootstrap-compatible `import "os"` sample with `Getwd`, `Stat`, `(*os.File).Stat`, `Create`, `Open`, `OpenFile`, `ReadFile`, `Mkdir`, `Rename`, `Remove`, `IsNotExist`, `Getpid`, `Getppid`, and process-local environment helpers
  - file lifecycle checks against a real writable KolibriOS path with append, `ModTime`, rename, environment, and cleanup validation
- `examples/fmt` - implemented
  - bootstrap-compatible `import "fmt"` sample with `Sprintf`, `Sprintln`, `Fprintf`, `Print`, `Printf`, `Println`, `Fscanln`, `Scanln`, and `Errorf`
  - formatted text and narrow scanning checks tied to a real KolibriOS file path, current-folder probe, and pipe-backed stdio capture
- `examples/console` - implemented
  - bootstrap `CONSOLE.OBJ` wrapper sample with DLL load, export lookup, `con_init`, `con_write_string`, `con_getch`, and `con_exit`
  - ordinary `fmt.Print`, `fmt.Printf`, and `fmt.Println` output through the active console-backed `os.Stdout`, plus direct `fmt.Scanln` on the same active console-backed `os.Stdin`
  - waits for `Esc` after the line-input prompt so the sample exercises both cooked console input and direct key reads
- `apps/diag` - implemented
  - fuller GUI diagnostics utility outside the public examples tree
  - runtime, file, narrow `syscall`, `os`, `fmt`, `time`, DLL-load, real `CONSOLE.OBJ` init/write/exit, stdout-console bridge, pipe-backed scan checks, `os.Stat`, `FileInfo.ModTime`, process/env semantics, and system probes in one reusable tool
  - headless QEMU diagnostics capture via debug console with `/FD/1` report fallback
- `tests/smokeapp` - implemented
  - headless autorun QEMU smoke for the documented bootstrap subset
  - runtime checks for strings, slices, interfaces, assertions, and timed wait
  - system geometry checks against a temporary pruned copy of the official `kolibri.img`

## Notes

- The syscall reference for all new bindings is `sysfuncs.txt`.
- Shared bootstrap build logic now lives in `mk/kolibri-app.mk`.
- `mk/kolibri-app.mk` now accepts an ordered `PACKAGE_DIRS` list so apps can
  precompile extra shared packages before the final app object.
- Bootstrap stdlib shim sources now live under `stdlib/`, while the compiled
  `.gox` export data remains available through the repository-root include path
  for ordinary imports such as `import "errors"` and `import "bytes"`.
- New applications can be scaffolded from `templates/basic-app` via
  `scripts/new-app.sh`.
- The linker script is generated from `mk/static.lds.in`.
- The linker template emits separate RX and RW load segments, so bootstrap
  builds no longer carry the old RWX `LOAD` warning.
- The linker template now derives the `MENUET01` memory reservation and stack
  top from the final linked image size plus `APP_STACK_RESERVE` (default
  `0x10000`), so larger bootstrap apps remain executable without manual header
  tuning.
- The current bootstrap runtime subset is documented in `docs/RUNTIME.md`.
- Focused compiler/runtime symbol checks live in `tests/runtime` and run via
  `scripts/check-runtime-probes.sh`.
- Host-side runtime behavior checks live in `tests/runtime/behavior.c` and run
  via `scripts/check-runtime-behavior.sh`.
  On hosts that cannot execute a 32-bit ELF directly, this harness falls back
  to a native-host build while `scripts/check-runtime-probes.sh` still validates
  the bootstrap `gccgo -m32` symbol inventory.
- The bootstrap runtime still has no real GC; current heap paths are
  `malloc`-based and the GC/barrier symbols only satisfy the validated subset.
- The current bootstrap stdlib-compatible package surface is tracked in
  `docs/STDLIB.md`.
- Template verification lives in `scripts/check-app-template.sh` and confirms
  that `scripts/new-app.sh` generates a buildable example under `examples/`.
- Emulator smoke verification lives in `scripts/check-emulator-smoke.sh`; it
  downloads the official `kolibri.img`, prunes non-system payload from a
  temporary copy to free space, replaces the existing `@HA` autorun slot with
  `tests/smokeapp`, boots QEMU headless, and waits for the smoke app to power
  the guest off after runtime plus system self-checks pass.
- Diagnostics verification lives in `scripts/check-diagnostics.sh`; it boots
  `apps/diag` headless, captures the emitted report from the QEMU debug console,
  and falls back to `/FD/1/GODIAG.TXT` only if debug-console capture is absent.
