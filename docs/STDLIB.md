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

## Compatibility Sample

Compatibility samples using ordinary import paths:

- `examples/files`
  - `import "errors"`
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

The samples still use the KolibriOS SDK for actual system interaction, but the
stdlib-shaped path, string, byte-slice, io, and error logic now follows ordinary Go package structure
instead of custom-only local helpers.

## Not Yet Supported

The following roadmap packages are still pending bootstrap implementations:

- `time`
- `os`
- `syscall`

Until they are explicitly documented here, they should be treated as
unsupported for the KolibriOS bootstrap target.
