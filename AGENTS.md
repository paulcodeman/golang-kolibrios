# Repository Instructions

## KolibriOS API Source Of Truth

- When adding or changing Go bindings for KolibriOS system calls, always consult `sysfuncs.txt` in the repository root first.
- Do not invent syscall numbers, register layouts, subfunction codes, packed arguments, or return conventions from memory.
- Treat `sysfuncs.txt` as the source of truth for:
  - the function number placed into `eax`
  - input/output register usage
  - packed argument formats such as `x * 65536 + width`
  - preserved registers and return-value behavior

## Where To Apply Changes

- Keep low-level syscall entrypoints aligned with `abi/syscalls_i386.asm`.
- Keep exported Go declarations aligned with `kos/raw.go`.
- Keep higher-level Go wrappers and types aligned with the low-level ABI in `kos/*.go` and `ui/*.go` when relevant.

## Implementation Rule

- If a new Go function wraps a KolibriOS API call, verify the exact calling convention in `sysfuncs.txt` before writing the Go signature or the assembly stub.
