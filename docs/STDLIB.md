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

### `path/filepath`

Implemented locally in the repository as a bootstrap shim.

Supported API:

- `filepath.Separator`
- `filepath.ListSeparator`
- `filepath.Abs`
- `filepath.Base`
- `filepath.Clean`
- `filepath.Dir`
- `filepath.Ext`
- `filepath.FromSlash`
- `filepath.IsAbs`
- `filepath.Join`
- `filepath.Split`
- `filepath.ToSlash`
- `filepath.VolumeName`

Current behavior notes:

- semantics are intentionally slash-first and map onto the current KolibriOS
  path model rather than Windows drive-letter rules
- `filepath.Clean`, `Join`, `Base`, `Dir`, `Split`, `Ext`, and `IsAbs` route
  through the local bootstrap `path` package after normalizing backslashes with
  `ToSlash`
- `filepath.Abs` resolves relative paths against `os.Getwd`
- `filepath.Separator` is `'/'`, `filepath.ListSeparator` is `':'`, and
  `filepath.VolumeName` always returns `""` for the current bootstrap target
- globbing, walking, symlink evaluation, volume-aware behavior, and the broader
  filepath surface are not implemented yet

### `strings`

Implemented locally in the repository as a bootstrap shim.

Supported API:

- `strings.Reader`
- `strings.NewReader`
- `strings.Builder`
- `(*strings.Reader).Len`
- `(*strings.Reader).Size`
- `(*strings.Reader).Reset`
- `(*strings.Reader).Read`
- `(*strings.Reader).ReadAt`
- `(*strings.Reader).ReadByte`
- `(*strings.Reader).UnreadByte`
- `(*strings.Reader).Seek`
- `(*strings.Reader).WriteTo`
- `(*strings.Builder).String`
- `(*strings.Builder).Len`
- `(*strings.Builder).Cap`
- `(*strings.Builder).Reset`
- `(*strings.Builder).Grow`
- `(*strings.Builder).Write`
- `(*strings.Builder).WriteByte`
- `(*strings.Builder).WriteString`
- `strings.Contains`
- `strings.Cut`
- `strings.HasPrefix`
- `strings.HasSuffix`
- `strings.Index`
- `strings.Join`
- `strings.LastIndex`
- `strings.Fields`
- `strings.ReplaceAll`
- `strings.Split`
- `strings.SplitN`
- `strings.TrimSpace`
- `strings.TrimPrefix`
- `strings.TrimSuffix`

Current behavior notes:

- matching and indexing are byte-oriented
- `strings.Builder` currently provides a narrow append-only text builder with
  `Write`, `WriteByte`, `WriteString`, `Len`, `Cap`, `String`, `Reset`, and
  `Grow`
- because the current bootstrap runtime still lacks a general `panic` path,
  `(*strings.Builder).Grow` treats non-positive values as a no-op instead of
  panicking and the standard library's copy-detection panic for copied builders
  is not implemented yet
- `strings.Cut` follows first-separator semantics and treats an empty separator
  as found at byte index `0`
- `strings.Split`/`SplitN` follow the current byte-oriented bootstrap model; an
  empty separator splits by single bytes rather than full UTF-8 runes
- `strings.TrimSpace` and `strings.Fields` currently use ASCII whitespace only
- `strings.Reader` is a narrow byte-oriented read/seek wrapper; `UnreadByte`
  only tracks the last successful `ReadByte`, and rune-aware reader APIs are
  not implemented yet
- helpers are intentionally ASCII/byte-focused for the current bootstrap stage
- higher-level Unicode-aware helpers and replacer APIs are not implemented yet

### `bytes`

Implemented locally in the repository as a bootstrap shim.

Supported API:

- `bytes.Reader`
- `bytes.Buffer`
- `bytes.NewReader`
- `bytes.NewBuffer`
- `bytes.NewBufferString`
- `(*bytes.Reader).Len`
- `(*bytes.Reader).Size`
- `(*bytes.Reader).Reset`
- `(*bytes.Reader).Read`
- `(*bytes.Reader).ReadAt`
- `(*bytes.Reader).ReadByte`
- `(*bytes.Reader).UnreadByte`
- `(*bytes.Reader).Seek`
- `(*bytes.Reader).WriteTo`
- `(*bytes.Buffer).Bytes`
- `(*bytes.Buffer).String`
- `(*bytes.Buffer).Len`
- `(*bytes.Buffer).Cap`
- `(*bytes.Buffer).Reset`
- `(*bytes.Buffer).Grow`
- `(*bytes.Buffer).Write`
- `(*bytes.Buffer).WriteByte`
- `(*bytes.Buffer).WriteString`
- `bytes.Contains`
- `bytes.Cut`
- `bytes.Equal`
- `bytes.HasPrefix`
- `bytes.HasSuffix`
- `bytes.Index`
- `bytes.IndexByte`
- `bytes.Join`
- `bytes.Fields`
- `bytes.ReplaceAll`
- `bytes.Split`
- `bytes.SplitN`
- `bytes.TrimSpace`
- `bytes.TrimPrefix`
- `bytes.TrimSuffix`

Current behavior notes:

- matching and indexing are byte-oriented
- `bytes.Buffer` currently provides a narrow append-only write buffer with
  `Write`, `WriteByte`, `WriteString`, `Bytes`, `String`, `Len`, `Cap`,
  `Reset`, `Grow`, `NewBuffer`, and `NewBufferString`
- `bytes.NewBuffer` keeps the caller-provided slice as the initial live backing
  store, while `Buffer.Bytes()` returns the current live slice view
- because the current bootstrap runtime still lacks a general `panic` path,
  `(*bytes.Buffer).Grow` treats non-positive values as a no-op instead of
  panicking
- `bytes.Cut` returns slices into the original input on success and
  `(s, nil, false)` when the separator is absent
- `bytes.Join` always allocates a new output slice and returns an empty
  non-nil slice for empty input
- `bytes.Split`/`SplitN` are byte-oriented, and an empty separator splits into
  one-byte slices
- `bytes.TrimSpace` and `bytes.Fields` currently use ASCII whitespace only
- `bytes.Reader` is a narrow byte-oriented read/seek wrapper with
  `UnreadByte`, `ReadAt`, `Seek`, and `WriteTo`, but no rune-aware reader APIs
  yet
- higher-level read-side buffer helpers, rune-aware behavior, and the broader
  `bytes` surface are not implemented yet

### `io`

Implemented locally in the repository as a bootstrap shim.

Supported API:

- `io.ReaderAt`
- `io.WriterTo`
- `io.ReaderFrom`
- `io.ByteReader`
- `io.ByteScanner`
- `io.Seeker`
- `io.ReadSeeker`
- `io.Reader`
- `io.Writer`
- `io.Closer`
- `io.ReadWriter`
- `io.ReadCloser`
- `io.WriteCloser`
- `io.StringWriter`
- `io.SeekStart`
- `io.SeekCurrent`
- `io.SeekEnd`
- `io.EOF`
- `io.ErrShortWrite`
- `io.Copy`
- `io.CopyBuffer`
- `io.ReadAll`
- `io.WriteString`

Current behavior notes:

- `io.ReadAll` consumes until `io.EOF` and returns the collected bytes with a
  `nil` error
- `io.Copy` and `io.CopyBuffer` now honor narrow `WriterTo` and `ReaderFrom`
  fast paths before falling back to the buffered reader/writer loop
- `io.WriteString` always falls back to `Writer.Write([]byte(s))`; it does not
  use `io.StringWriter` optimization yet
- `io.EOF` and `io.ErrShortWrite` are local bootstrap sentinels implemented as
  concrete package values to avoid init-time cross-package calls

### `time`

Implemented locally in the repository as a bootstrap shim.

Supported API:

- `time.Duration`
- `time.Nanosecond`
- `time.Microsecond`
- `time.Millisecond`
- `time.Second`
- `time.Minute`
- `time.Hour`
- `time.Month`
- `time.January` through `time.December`
- `time.Time`
- `time.Now`
- `time.Unix`
- `time.Sleep`
- `time.Since`
- `time.Time.Add`
- `time.Time.Sub`
- `time.Time.Before`
- `time.Time.After`
- `time.Time.Equal`
- `time.Time.IsZero`
- `time.Time.Unix`
- `time.Time.Nanosecond`
- `time.Time.Second`
- `time.Time.Minute`
- `time.Time.Hour`
- `time.Time.Day`
- `time.Time.Month`
- `time.Time.Year`

Current behavior notes:

- `time.Now` assembles the wall clock from syscalls `29` and `3`, and carries a
  separate monotonic component from `26.10` so `time.Since` and `Time.Sub`
  remain useful even though the current wall clock itself is only second
  resolution
- the current bootstrap mapping for syscall `29` expands the documented
  two-digit year `YY` into `2000 + YY`; this is now part of the documented
  compatibility contract for the bootstrap path
- `time.Unix` returns wall-clock-only values without a monotonic component
- `Time.Before`, `Time.After`, `Time.Equal`, and `Time.Sub` prefer the
  monotonic component when both operands were produced by `time.Now` or derived
  from such values through `Time.Add`; otherwise they fall back to wall-clock
  comparison
- `time.Sleep` rounds positive durations up to the documented 1/100-second
  kernel delay granularity from syscall `5`
- time zones, locations, parsing, formatting, timers, tickers, and the broader
  calendar surface are not implemented yet

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

### `bufio`

Implemented locally in the repository as a bootstrap shim.

Supported API:

- `bufio.Reader`
- `bufio.NewReader`
- `(*bufio.Reader).Read`
- `(*bufio.Reader).ReadByte`
- `(*bufio.Reader).UnreadByte`
- `(*bufio.Reader).ReadBytes`
- `(*bufio.Reader).ReadString`
- `bufio.Writer`
- `bufio.NewWriter`
- `(*bufio.Writer).Write`
- `(*bufio.Writer).WriteByte`
- `(*bufio.Writer).WriteString`
- `(*bufio.Writer).Flush`
- `bufio.Scanner`
- `bufio.NewScanner`
- `(*bufio.Scanner).Scan`
- `(*bufio.Scanner).Text`
- `(*bufio.Scanner).Bytes`
- `(*bufio.Scanner).Err`
- `(*bufio.Scanner).Buffer`
- `(*bufio.Scanner).Split`
- `bufio.SplitFunc`
- `bufio.ScanLines`
- `bufio.ScanWords`
- `bufio.ScanBytes`
- `bufio.ErrInvalidUnreadByte`
- `bufio.ErrTooLong`
- `bufio.MaxScanTokenSize`

Current behavior notes:

- the current bootstrap reader and scanner surfaces are byte-oriented and
  intentionally ASCII-focused
- `bufio.Reader` currently supports sequential buffered reads plus
  single-byte unread through `UnreadByte`
- `bufio.Writer` provides buffered writes and explicit `Flush`; although the
  kernel still has no documented close-handle syscall for `77` pipe handles,
  the bootstrap `os.Pipe` wrappers now emulate in-process close semantics, so
  closing the writer end produces `io.EOF` on the peer after buffered data is
  drained and closing the reader end makes subsequent writer calls fail with
  `syscall.EPIPE`
- `bufio.Scanner` currently supports `ScanLines`, `ScanWords`, and `ScanBytes`
  with a configurable maximum token size via `Buffer`
- split functions are expected to follow the normal scanner contract, but the
  broader standard package surface such as `ReadSlice`, `Peek`, `Discard`,
  `ReadRune`, `ReadLine`, `NewReadWriter`, `NewWriterSize`, and
  `Scanner.Bytes()` lifetime guarantees beyond the next scan are not
  implemented yet

### `strconv`

Implemented locally in the repository as a bootstrap shim.

Supported API:

- `strconv.IntSize`
- `strconv.ErrRange`
- `strconv.ErrSyntax`
- `strconv.NumError`
- `(*strconv.NumError).Error`
- `(*strconv.NumError).Unwrap`
- `strconv.FormatBool`
- `strconv.AppendBool`
- `strconv.ParseBool`
- `strconv.Itoa`
- `strconv.Atoi`
- `strconv.FormatInt`
- `strconv.FormatUint`
- `strconv.AppendInt`
- `strconv.AppendUint`
- `strconv.ParseInt`
- `strconv.ParseUint`

Current behavior notes:

- integer formatting and parsing currently support bases `2` through `36`,
  plus the usual `base == 0` auto-prefix handling for `ParseInt` and
  `ParseUint` with `0x`, `0b`, `0o`, and legacy leading-zero octal input
- `ParseBool` accepts the narrow standard tokens `1`, `0`, `t`, `f`, `T`,
  `F`, `true`, and `false` with ASCII case-folding for the full-word forms
- `Atoi` and `ParseInt` follow the local bootstrap `IntSize == 32` contract
- `NumError.Unwrap` exposes the underlying `ErrRange` or `ErrSyntax` sentinel,
  so ordinary `errors.Is` matching works
- because the current bootstrap runtime still lacks a general `panic` path and
  64-bit division helpers, invalid formatting bases for `FormatInt`,
  `FormatUint`, `AppendInt`, and `AppendUint` currently coerce to base `10`
  instead of panicking like the full standard library
- floating-point, quoting, rune-classification, and the broader `strconv`
  surface are not implemented yet

### `os`

Implemented locally in the repository as a bootstrap shim.

Supported API:

- `os.File`
- `os.FileInfo`
- `os.FileMode`
- `os.ModeDir`
- `os.FileMode.IsDir`
- `os.PathSeparator`
- `os.PathListSeparator`
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
- `os.Args`
- `os.PathError`
- `os.LinkError`
- `os.Getwd`
- `os.Getpid`
- `os.Getppid`
- `os.Exit`
- `os.Getenv`
- `os.LookupEnv`
- `os.Setenv`
- `os.Unsetenv`
- `os.Clearenv`
- `os.Environ`
- `os.Stat`
- `os.Open`
- `os.Create`
- `os.OpenFile`
- `(*os.File).Seek`
- `(*os.File).ReadAt`
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
  surface with `Name`, `Size`, `Mode`, `ModTime`, `IsDir`, and `Sys`; `Mode`
  currently exposes only `ModeDir`, `ModTime` is assembled from the raw
  `kos.FileInfo` date/time fields through the local bootstrap `time` package,
  and `Sys()` returns the underlying `kos.FileInfo` record for callers that
  still need KolibriOS-specific metadata such as raw file attributes
- `OpenFile` currently supports the narrow bootstrap flag set documented above;
  descriptor duplication, permissions, and sync semantics are not implemented
  yet
- `(*os.File).Seek` and `(*os.File).ReadAt` currently work only for path-backed
  files; fd-backed files such as pipes, stdio handles, and the active-console
  bridge still return `ErrInvalid` for random-access operations
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
- `os.Args` currently defaults to a single empty program slot `[]string{""}`;
  the bootstrap path does not yet parse a loader command line into ordinary Go
  argv semantics
- `os.Getpid` currently reports the current KolibriOS thread id through the
  documented thread-info path, `os.Getppid` currently returns `0`, and
  `os.Exit(code)` currently ignores the numeric code and exits through
  `kos.Exit()`
- `os.Getenv`, `os.LookupEnv`, `os.Setenv`, `os.Unsetenv`, `os.Clearenv`, and
  `os.Environ` are implemented as a process-local userspace store; there is no
  inherited kernel environment block yet, so values live only inside the
  running bootstrap process
- `os.IsNotExist` currently follows the unwrap chain for bootstrap `os.PathError`
  and `os.LinkError` values and checks against the local `os.ErrNotExist`
  sentinel
- `(*os.File).Close` still does not invoke a documented kernel close-handle
  syscall, because the current `sysfuncs.txt` contract for function `77`
  exposes `Read`, `Write`, and `Pipe`, but no close operation; for
  `os.Pipe`-created wrapper pairs the bootstrap layer now emulates in-process
  close semantics so the reader sees `io.EOF` after the last local writer is
  closed and the writer sees `syscall.EPIPE` after the last local reader is
  closed
- `(*os.File).Stat` currently works only for path-backed files; fd-backed files
  such as pipes and stdio handles still return `ErrInvalid`
- `Rename` resolves ordinary Go-style relative and absolute paths into the
  special KolibriOS `80.10` target-path contract and currently supports only
  same-volume rename or move operations
- `WriteFile`, `Mkdir`, and `OpenFile` ignore Unix permission bits for now and
  keep only the narrow bootstrap behavior required by the compatibility sample
- directory iteration, process spawning, command execution, and the broader
  `os` surface are not implemented yet

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

For ordinary import paths such as `import "errors"` or
`import "path/filepath"`, the current bootstrap shim sources live under
`stdlib/<package>`. Top-level compiled export data is still exposed through the
shared `-I$(ROOT)` include path, while nested import paths are emitted under the
shared `-I$(ROOT)/.pkg` include root so apps keep the ordinary Go import path
even though the repository layout is now cleaner.

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
  - metadata probe through ordinary `os.Stat`
- `examples/filepath`
  - `import "path/filepath"`
  - slash and backslash normalization with `Clean`, `ToSlash`, and `FromSlash`
  - component extraction with `Split`, `Base`, `Dir`, `Ext`, and `VolumeName`
  - relative-path resolution with `Abs`
  - metadata probe through ordinary `os.Stat`
- `examples/strings`
  - `import "strings"`
  - path assembly with `Join`
  - text assembly through `strings.Builder`, `Write`, `WriteByte`,
    `WriteString`, `String`, `Len`, `Cap`, `Grow`, and `Reset`
  - byte-oriented matching via `Contains`, `HasPrefix`, `HasSuffix`, `Index`, and `LastIndex`
  - delimiter and suffix trimming with `Cut`, `TrimPrefix`, and `TrimSuffix`
  - current-folder and metadata probes through ordinary `os.Getwd` and `os.Stat`
- `examples/bytes`
  - `import "bytes"`
  - byte-slice path assembly with `Join`
  - write-buffer assembly through `bytes.Buffer`, `NewBuffer`,
    `NewBufferString`, `Write`, `WriteByte`, `WriteString`, `Bytes`, `String`,
    `Len`, `Cap`, `Grow`, and `Reset`
  - byte-oriented matching via `Equal`, `Contains`, `HasPrefix`, `HasSuffix`, `Index`, and `IndexByte`
  - delimiter and suffix trimming with `Cut`, `TrimPrefix`, and `TrimSuffix`
  - current-folder and metadata probes through ordinary `os.Getwd` and `os.Stat`
- `examples/io`
  - `import "io"`
  - chunked stream reads with `ReadAll`
  - byte transfer through `Copy`
  - string-to-writer bridge through `WriteString`
  - file and current-folder probes through ordinary `os.ReadFile`, `os.Stat`, and `os.Getwd`
- `examples/time`
  - `import "time"`
  - wall clock access through `Now`, `Unix`, `Year`, `Month`, `Day`, `Hour`, `Minute`, and `Second`
  - monotonic duration path through `Sleep`, `Since`, and `Sub`
  - documented bootstrap year expansion `YY => 2000+YY`
- `examples/os`
  - `import "os"`
  - current-folder lookup through `Getwd`
  - metadata lookup through `Stat`, `(*os.File).Stat`, and `FileInfo.ModTime`
  - file create, append, read, and copy flow through `Create`, `OpenFile`, `ReadFile`, and `Open`
  - file rename and cleanup flow through `Rename`, `Remove`, and `IsNotExist`
  - process and environment flow through `Getpid`, `Getppid`, `Args`,
    `Setenv`, `LookupEnv`, `Environ`, and `Unsetenv`
- `examples/fmt`
  - `import "fmt"`
  - formatted strings via `Sprintf` and `Sprintln`
  - writer formatting via `Fprintf`
  - stdout-style formatting via `Print`, `Printf`, and `Println` redirected
    through a temporary `os.Pipe`
  - formatted error construction via `Errorf`
  - ordinary `os.Stdout` reassignment for bootstrap stdout capture
  - ordinary `os.Getwd`, `os.Stat`, and `os.ReadFile` for the file/cwd probe
- `examples/bufio`
  - `import "bufio"`
  - buffered pipe writes through `Writer`, `WriteString`, `WriteByte`, and `Flush`
  - buffered reads through `Reader`, `ReadByte`, `UnreadByte`, `ReadString`, and `ReadBytes`
  - token scanning through `Scanner`, `ScanLines`, `ScanWords`, and `ScanBytes`
  - ordinary `os.Pipe`, `os.Getwd`, and `os.Stat` for the runtime probe, plus
    EOF-after-close and broken-pipe validation through `io.EOF` and `syscall.EPIPE`
- `examples/strconv`
  - `import "strconv"`
  - bool and integer formatting through `FormatBool`, `FormatInt`, `FormatUint`, and `Itoa`
  - narrow parsing through `ParseBool`, `ParseInt`, `ParseUint`, and `Atoi`
  - append helpers through `AppendBool`, `AppendInt`, and `AppendUint`
  - wrapped `ErrRange` and `ErrSyntax` classification through ordinary `errors.Is`
  - ordinary `os.Getwd` and `os.Stat` for the cwd and file probe
- `apps/diag`
  - headless regression coverage for `strings.Builder` and `bytes.Buffer`
  - headless `bufio` regression coverage for reader, writer, scanner,
    EOF-after-close, and broken-pipe behavior on pipe-backed stdio
  - headless `strconv` regression coverage for bool/int format, parse, append, and `NumError` sentinel matching

The samples still use the KolibriOS SDK for actual system interaction, but the
stdlib-shaped path, filepath, string, byte-slice, io, os, fmt, strconv, time,
and error logic now follows ordinary Go package structure instead of
custom-only local helpers.

## Not Yet Supported

Any package surface not explicitly documented above should still be treated as
unsupported for the KolibriOS bootstrap target.
