# Syscall Inventory

This file maps the currently exported raw bindings to `sysfuncs.txt`.

`sysfuncs.txt` is the source of truth for syscall numbers, register usage,
packed arguments, and return conventions.

## Covered Raw Bindings

| Sysfunc | Raw binding | Higher-level wrapper | Notes |
| --- | --- | --- | --- |
| `0` | `kos.Window` | `kos.OpenWindow` | Window creation and redraw |
| `2` | `kos.GetKey` | `kos.ReadKey` | Raw keyboard event decode |
| `3` | `kos.GetTime` | - | Returns packed BCD time |
| `4` | `kos.WriteText` | `kos.DrawText` | UTF-8 flag is set in the asm stub |
| `5` | `kos.Sleep` | - | Delay in 1/100 second |
| `8` | `kos.CreateButton` | `kos.DrawButton` | Button definition |
| `9` | `kos.GetThreadInfo` | `kos.ReadThreadInfo`, `kos.ReadCurrentThreadInfo` | Reads and decodes the 1-KiB thread info buffer |
| `10` | `kos.Event` | `kos.WaitEvent` | Wait for event |
| `11` | `kos.CheckEvent` | `kos.PollEvent` | Non-blocking event check |
| `12` | `kos.Redraw` | `kos.BeginRedraw`, `kos.EndRedraw` | Window redraw flow |
| `13` | `kos.DrawBar` | `kos.FillRect` | Rectangle fill |
| `14` | `kos.GetScreenSize` | `kos.ScreenSize` | Returns packed lower-right coords |
| `17` | `kos.GetButtonID` | `kos.CurrentButtonID` | Button ID extraction |
| `18.16` | `kos.GetFreeRAM` | `kos.FreeRAMKB` | Returns free RAM in kilobytes |
| `18.17` | `kos.GetTotalRAM` | `kos.TotalRAMKB` | Returns total RAM in kilobytes |
| `23` | `kos.WaitEventTimeout` | `kos.WaitEventFor` | Timed wait |
| `37.0` | `kos.GetMouseScreenPosition` | `kos.MouseScreenPosition` | Mouse coordinates on the screen |
| `37.1` | `kos.GetMouseWindowPosition` | `kos.MouseWindowPosition` | Mouse coordinates relative to the window |
| `37.2` | `kos.GetMouseButtonState` | `kos.MouseHeldButtons` | Held mouse button states |
| `37.3` | `kos.GetMouseButtonEventState` | `kos.MouseButtons` | Mouse button states plus edge events |
| `37.7` | `kos.GetMouseScrollData` | `kos.MouseScrollDelta` | Signed horizontal/vertical scroll deltas |
| `38` | `kos.DrawLine` | `kos.StrokeLine` | Line draw |
| `40` | `kos.SetEventMask` | `kos.SwapEventMask` | Event filtering |
| `63` | `kos.DebugOutHex`, `kos.DebugOutChar`, `kos.DebugOutStr` | `kos.DebugHex`, `kos.DebugChar`, `kos.DebugString` | Debug board helpers |
| `68.12` | `malloc` | runtime-only | Lazy heap init via `68.11` |
| `68.13` | `free` | runtime-only | Lazy heap init via `68.11` |
| `68.20` | `realloc` | runtime-only | Lazy heap init via `68.11` |
| `70` | `kos.FileSystem` | - | Raw long-name file system interface |
| `71.2` | `kos.SetCaption` | `kos.SetWindowTitle` | Uses explicit UTF-8 encoding |
| `80` | `kos.FileSystemEncoded` | `kos.ReadFile`, `kos.ReadDirectory`, `kos.CreateOrRewriteFile`, `kos.WriteFile`, `kos.SetFileSize`, `kos.GetPathInfo`, `kos.SetPathInfo`, `kos.StartApplication`, `kos.DeletePath`, `kos.CreateDirectory`, `kos.RenamePath` | UTF-8-friendly file system layer |
| `-1` | `kos.Exit` | - | Terminate thread/process |

## Priority Gaps For Phase 1

These are the nearest missing wrappers with high utility and relatively low ABI
risk:

- `71.1` caption update using the alternative encoding contract
- `72` window messaging
- `18.13` kernel version
- `48.4/48.5` skin and working-area queries

## Notes

- The current raw layer intentionally stays thin; friendly Go wrappers live in
  `kos/*.go`.
- Window, text, and caption helpers use Go strings, but the asm/runtime glue is
  responsible for adapting them to KolibriOS string conventions.
