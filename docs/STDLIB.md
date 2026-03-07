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

## Build Contract

The shared app makefile now accepts an ordered `PACKAGE_DIRS` list.

This lets the bootstrap build precompile additional shared packages before the
app object itself, instead of hardcoding only `kos` and `ui`.

For ordinary import paths such as `import "errors"`, the current bootstrap shim
package lives in a same-named repository-root directory and is exposed through
the shared `-I$(ROOT)` include path.

Current example:

```make
PACKAGE_DIRS = kos errors ui
```

Order matters. Later packages may depend on earlier package export data.

## Compatibility Sample

`examples/files` is the first compatibility sample using an ordinary import path:

- `import "errors"`
- wrapped file-probe failures with `Unwrap`
- sentinel classification with `errors.Is`

The sample still uses the KolibriOS SDK for actual system interaction, but the
error-handling path now follows normal Go package structure instead of a
custom-only local helper.

## Not Yet Supported

The following roadmap packages are still pending bootstrap implementations:

- `bytes`
- `strings`
- `io`
- `time`
- `path`
- `os`
- `syscall`

Until they are explicitly documented here, they should be treated as
unsupported for the KolibriOS bootstrap target.
