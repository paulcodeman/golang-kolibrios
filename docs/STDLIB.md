# Bootstrap Stdlib Surface

This file tracks the current bootstrap-compatible standard-library-shaped
packages for the KolibriOS `gccgo` path.

The goal is not full standard library support yet. The goal is to let the
current SDK grow from custom-only imports toward ordinary Go package patterns
without pulling in the full native Go toolchain.

## Current Package Surface

### `errors`

Implemented locally in the repository as a bootstrap shim.

Supported API:

- `errors.New`
- `errors.Unwrap`
- `errors.Is`

Current behavior notes:

- `errors.New` returns a unique heap-allocated error value for each call.
- `errors.Is` currently follows a simple unwrap chain and checks equality
  against sentinel errors.
- method-based custom `Is(error) bool` matching is not implemented yet.
- `errors.As`, `errors.Join`, and formatted error construction are not
  implemented yet.

### `path`

Implemented locally in the repository as a bootstrap shim.

Supported API:

- `path.Base`
- `path.Clean`
- `path.Dir`
- `path.Ext`
- `path.IsAbs`
- `path.Join`
- `path.Split`

Current behavior notes:

- semantics are slash-based, matching KolibriOS path conventions
- repeated `/`, `.` segments, and `..` collapse are handled by `path.Clean`
- rooted paths clamp parent traversal at `/`
- relative paths preserve leading `..` segments
- globbing helpers such as `Match` are not implemented yet

### `strings`

Implemented locally in the repository as a bootstrap shim.

Supported API:

- `strings.Contains`
- `strings.Cut`
- `strings.HasPrefix`
- `strings.HasSuffix`
- `strings.Index`
- `strings.Join`
- `strings.LastIndex`
- `strings.TrimPrefix`
- `strings.TrimSuffix`

Current behavior notes:

- matching and indexing are byte-oriented
- `strings.Cut` follows first-separator semantics and treats an empty separator
  as found at byte index `0`
- helpers are intentionally ASCII/byte-focused for the current bootstrap stage
- higher-level Unicode-aware helpers and replacer/builder APIs are not
  implemented yet

### `bytes`

Implemented locally in the repository as a bootstrap shim.

Supported API:

- `bytes.Contains`
- `bytes.Cut`
- `bytes.Equal`
- `bytes.HasPrefix`
- `bytes.HasSuffix`
- `bytes.Index`
- `bytes.IndexByte`
- `bytes.Join`
- `bytes.TrimPrefix`
- `bytes.TrimSuffix`

Current behavior notes:

- matching and indexing are byte-oriented
- `bytes.Cut` returns slices into the original input on success and
  `(s, nil, false)` when the separator is absent
- `bytes.Join` always allocates a new output slice and returns an empty
  non-nil slice for empty input
- higher-level buffer, reader, and Unicode-aware helpers are not implemented
  yet

### `io`

Implemented locally in the repository as a bootstrap shim.

Supported API:

- `io.Reader`
- `io.Writer`
- `io.Closer`
- `io.ReadWriter`
- `io.ReadCloser`
- `io.WriteCloser`
- `io.StringWriter`
- `io.EOF`
- `io.ErrShortWrite`
- `io.Copy`
- `io.CopyBuffer`
- `io.ReadAll`
- `io.WriteString`

Current behavior notes:

- `io.ReadAll` consumes until `io.EOF` and returns the collected bytes with a
  `nil` error
- `io.Copy` and `io.CopyBuffer` use a simple reader/writer loop; `ReaderFrom`
  and `WriterTo` fast paths are not implemented yet
- `io.WriteString` always falls back to `Writer.Write([]byte(s))`; it does not
  use `io.StringWriter` optimization yet
- `io.EOF` and `io.ErrShortWrite` are local bootstrap sentinels implemented as
  concrete package values to avoid init-time cross-package calls

### `fmt`

Implemented locally in the repository as a bootstrap shim.

Supported API:

- `fmt.Stringer`
- `fmt.Sprint`
- `fmt.Sprintln`
- `fmt.Sprintf`
- `fmt.Fprint`
- `fmt.Fprintln`
- `fmt.Fprintf`
- `fmt.Fscan`
- `fmt.Fscanln`
- `fmt.Print`
- `fmt.Println`
- `fmt.Printf`
- `fmt.Scan`
- `fmt.Scanln`
- `fmt.Errorf`

Current behavior notes:

- current formatting support is intentionally narrow and covers `%s`, `%d`,
  `%x`, `%X`, `%t`, `%v`, `%c`, and `%%`
- strings, byte slices, booleans, signed and unsigned integers, `error`, and
  `fmt.Stringer` values are supported by the current formatter
- `fmt.Errorf` currently formats through the local bootstrap formatter and then
  returns `errors.New(...)`; `%w` is rendered like `%v` and does not yet create
  an unwrap chain
- `fmt.Print`, `fmt.Printf`, and `fmt.Println` now route through `os.Stdout`,
  so ordinary stdout-style Go code can be exercised either by redirecting
  `os.Stdout` to a pipe-backed `*os.File` or by opening an active
  `CONSOLE.OBJ` backend through `kos.OpenConsole`
- `fmt.Fscan`, `fmt.Fscanln`, `fmt.Scan`, and `fmt.Scanln` currently support
  narrow token parsing for `string`, `bool`, and signed or unsigned integer
  pointer targets
- `fmt.Scan` and `fmt.Scanln` route through `os.Stdin`; when an active
  `CONSOLE.OBJ` instance has been opened through `kos.OpenConsole`, the
  bootstrap `os.Stdin` path reads one cooked line at a time through
  `con_gets`
- `Scanf`, width/precision directives, floating-point formatting or scanning,
  maps, structs, custom scanner interfaces, and the broader printing/scanning
  surface are not implemented yet

### `os`

Implemented locally in the repository as a bootstrap shim.

Supported API:

- `os.File`
- `os.FileInfo`
- `os.FileMode`
- `os.ModeDir`
- `os.FileMode.IsDir`
- `os.O_RDONLY`
- `os.O_WRONLY`
- `os.O_RDWR`
- `os.O_CREATE`
- `os.O_TRUNC`
- `os.O_APPEND`
- `os.ErrInvalid`
- `os.ErrPermission`
- `os.ErrExist`
- `os.ErrNotExist`
- `os.ErrClosed`
- `os.Stdin`
- `os.Stdout`
- `os.Stderr`
- `os.DefaultStdin`
- `os.DefaultStdout`
- `os.DefaultStderr`
- `os.PathError`
- `os.LinkError`
- `os.Getwd`
- `os.Stat`
- `os.Open`
- `os.Create`
- `os.OpenFile`
- `(*os.File).Stat`
- `os.Pipe`
- `os.IsNotExist`
- `os.ReadFile`
- `os.WriteFile`
- `os.Mkdir`
- `os.Remove`
- `os.Rename`

Current behavior notes:

- current files are name-plus-offset wrappers over KolibriOS path-based file
  syscalls, but the package now also has a narrow fd-backed mode for
  `os.Stdin`, `os.Stdout`, `os.Stderr`, and `os.Pipe`
- `os.Stat` and `(*os.File).Stat` expose a narrow bootstrap `os.FileInfo`
  surface with `Name`, `Size`, `Mode`, `IsDir`, and `Sys`; `Mode` currently
  exposes only `ModeDir`, and `Sys()` returns the underlying `kos.FileInfo`
  record for callers that still need KolibriOS-specific metadata such as raw
  file attributes
- the bootstrap `os.FileInfo` type intentionally does not expose `ModTime`
  yet; that part waits on the local `time` shim
- `OpenFile` currently supports the narrow bootstrap flag set documented above;
  descriptor duplication, permissions, and sync semantics are not implemented
  yet
- the fd-backed path currently follows the documented `77.10/77.11/77.13`
  contracts, which are currently specified for pipe descriptors; the default
  `os.Stdout` and `os.Stderr` handles additionally route through the active
  `CONSOLE.OBJ` instance when one has been opened through the bootstrap `kos`
  console wrapper, and `os.Stdin` switches to the active console line-input
  bridge on the same condition
- the bootstrap runtime currently re-materializes default stdio handles through
  `os.DefaultStdin`, `os.DefaultStdout`, and `os.DefaultStderr`; `fmt.Print*`
  uses that path internally so ordinary stdout-style Go code stays usable even
  when imported package globals arrive zeroed
- `os.IsNotExist` currently follows the unwrap chain for bootstrap `os.PathError`
  and `os.LinkError` values and checks against the local `os.ErrNotExist`
  sentinel
- `(*os.File).Close` currently marks bootstrap fd-backed files closed locally,
  but it does not yet invoke a documented kernel close-handle syscall
- `(*os.File).Stat` currently works only for path-backed files; fd-backed files
  such as pipes and stdio handles still return `ErrInvalid`
- `Rename` resolves ordinary Go-style relative and absolute paths into the
  special KolibriOS `80.10` target-path contract and currently supports only
  same-volume rename or move operations
- `WriteFile`, `Mkdir`, and `OpenFile` ignore Unix permission bits for now and
  keep only the narrow bootstrap behavior required by the compatibility sample
- `Stat`, directory iteration, environment handling, process spawning, and the
  broader `os` surface are not implemented yet

### `syscall`

Implemented locally in the repository as a bootstrap shim.

Supported API:

- `syscall.Errno`
- `syscall.EBADF`
- `syscall.EINVAL`
- `syscall.EFAULT`
- `syscall.ENFILE`
- `syscall.EMFILE`
- `syscall.EPIPE`
- `syscall.O_CLOEXEC`
- `syscall.Read`
- `syscall.Write`
- `syscall.Pipe`
- `syscall.Pipe2`

Current behavior notes:

- `syscall.Read`, `syscall.Write`, `syscall.Pipe`, and `syscall.Pipe2` are
  backed directly by the documented `77.10`, `77.11`, and `77.13` contracts in
  `sysfuncs.txt`
- the current kernel documentation marks `77.10` and `77.11` as pipe-descriptor
  paths, so the bootstrap fd compatibility layer is currently validated through
  pipes rather than full Unix-style file descriptor coverage
- `syscall.Errno.Error()` currently formats only a narrow set of known bootstrap
  errno values
- close, dup, poll/select, signals, and the broader syscall surface are not
  implemented yet

## Build Contract

The shared app makefile now accepts an ordered `PACKAGE_DIRS` list.

This lets the bootstrap build precompile additional shared packages before the
app object itself, instead of hardcoding only `kos` and `ui`.

For ordinary import paths such as `import "errors"`, the current bootstrap shim
sources live under `stdlib/<package>`. Their compiled export data is still
exposed through the shared `-I$(ROOT)` include path, so apps keep the ordinary
Go import path even though the repository layout is now cleaner.

Current example:

```make
PACKAGE_DIRS = kos io ui
```

Order matters. Later packages may depend on earlier package export data.
The generated linker header now derives the app memory reservation from the
final linked image size plus `APP_STACK_RESERVE` (default `0x10000`), so larger
bootstrap apps remain executable without hand-tuned `MENUET01` header values.

## Compatibility Sample

Compatibility samples using ordinary import paths:

- `examples/files`
  - `import "errors"`, `import "io"`, `import "os"`
  - metadata probe through `os.Stat` with raw KolibriOS attributes available via `FileInfo.Sys()`
  - ordinary `os.Open` / `Read` / `Close` for the preview read path
  - wrapped file-probe failures with `Unwrap`
  - sentinel classification with `errors.Is`
- `examples/path`
  - `import "path"`
  - slash normalization with `Clean` and `Join`
  - component extraction with `Split`, `Base`, `Dir`, and `Ext`
- `examples/strings`
  - `import "strings"`
  - path assembly with `Join`
  - byte-oriented matching via `Contains`, `HasPrefix`, `HasSuffix`, `Index`, and `LastIndex`
  - delimiter and suffix trimming with `Cut`, `TrimPrefix`, and `TrimSuffix`
- `examples/bytes`
  - `import "bytes"`
  - byte-slice path assembly with `Join`
  - byte-oriented matching via `Equal`, `Contains`, `HasPrefix`, `HasSuffix`, `Index`, and `IndexByte`
  - delimiter and suffix trimming with `Cut`, `TrimPrefix`, and `TrimSuffix`
- `examples/io`
  - `import "io"`
  - chunked stream reads with `ReadAll`
  - byte transfer through `Copy`
  - string-to-writer bridge through `WriteString`
- `examples/os`
  - `import "os"`
  - current-folder lookup through `Getwd`
  - metadata lookup through `Stat` and `(*os.File).Stat`
  - file create, append, read, and copy flow through `Create`, `OpenFile`, `ReadFile`, and `Open`
  - file rename and cleanup flow through `Rename`, `Remove`, and `IsNotExist`
- `examples/fmt`
  - `import "fmt"`
  - formatted strings via `Sprintf` and `Sprintln`
  - writer formatting via `Fprintf`
  - stdout-style formatting via `Print`, `Printf`, and `Println` redirected
    through a temporary `os.Pipe`
  - formatted error construction via `Errorf`
  - ordinary `os.Stdout` reassignment for bootstrap stdout capture

The samples still use the KolibriOS SDK for actual system interaction, but the
stdlib-shaped path, string, byte-slice, io, os, and error logic now follows ordinary Go package structure
instead of custom-only local helpers.

## Not Yet Supported

The following roadmap packages are still pending bootstrap implementations:

- `time`

Until they are explicitly documented here, they should be treated as
unsupported for the KolibriOS bootstrap target.
