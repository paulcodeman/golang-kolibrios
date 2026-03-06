# Syscall Inventory

This file maps the currently exported raw bindings to `sysfuncs.txt`.

`sysfuncs.txt` is the source of truth for syscall numbers, register usage,
packed arguments, and return conventions.

## Covered Raw Bindings

| Sysfunc | Raw binding | Higher-level wrapper | Notes |
| --- | --- | --- | --- |
| `0` | `kos.Window` | `kos.OpenWindow` | Window creation and redraw |
| `2` | `kos.GetKey` | `kos.ReadKey` | Raw keyboard event decode |
| `3` | `kos.GetTime` | `kos.SystemTime` | Returns packed BCD time |
| `4` | `kos.WriteText` | `kos.DrawText` | UTF-8 flag is set in the asm stub |
| `5` | `kos.Sleep` | `kos.SleepCentiseconds`, `kos.SleepMilliseconds`, `kos.SleepSeconds` | Delay in 1/100 second |
| `8` | `kos.CreateButton` | `kos.DrawButton` | Button definition |
| `9` | `kos.GetThreadInfo` | `kos.ReadThreadInfo`, `kos.ReadCurrentThreadInfo` | Reads and decodes the 1-KiB thread info buffer |
| `10` | `kos.Event` | `kos.WaitEvent` | Wait for event |
| `11` | `kos.CheckEvent` | `kos.PollEvent` | Non-blocking event check |
| `12` | `kos.Redraw` | `kos.BeginRedraw`, `kos.EndRedraw` | Window redraw flow |
| `13` | `kos.DrawBar` | `kos.FillRect` | Rectangle fill |
| `14` | `kos.GetScreenSize` | `kos.ScreenSize` | Returns packed lower-right coords |
| `17` | `kos.GetButtonID` | `kos.CurrentButtonID` | Button ID extraction |
| `18.3` | `kos.FocusWindowBySlot` | `kos.FocusWindowSlot` | Focus the window of a thread slot |
| `18.7` | `kos.GetActiveWindowSlotRaw` | `kos.ActiveWindowSlot` | Return the active window thread slot |
| `18.9` | `kos.SystemShutdown` | `kos.Shutdown`, `kos.PowerOff`, `kos.Reboot`, `kos.RestartKernel` | System shutdown and reboot helpers |
| `18.13` | `kos.GetKernelVersion` | `kos.KernelVersion` | Reads the kernel version/ABI metadata block |
| `18.16` | `kos.GetFreeRAM` | `kos.FreeRAMKB` | Returns free RAM in kilobytes |
| `18.17` | `kos.GetTotalRAM` | `kos.TotalRAMKB` | Returns total RAM in kilobytes |
| `21.2` | `kos.SetKeyboardLayoutRaw`, `kos.SetKeyboardLanguageRaw` | `kos.SetKeyboardLayoutTable`, `kos.SetKeyboardLayoutLanguage` | Per-layout keyboard table upload and global layout language id |
| `21.5` | `kos.SetSystemLanguageRaw` | `kos.SetSystemLanguage` | Writes the global system language id used by higher-level apps such as `@taskbar` |
| `23` | `kos.WaitEventTimeout` | `kos.WaitEventFor` | Timed wait |
| `26.2` | `kos.GetKeyboardLayoutRaw`, `kos.GetKeyboardLanguageRaw` | `kos.ReadKeyboardLayoutTable`, `kos.KeyboardLayoutLanguage` | Per-layout keyboard table dump and current layout language id |
| `26.5` | `kos.GetSystemLanguageRaw` | `kos.SystemLanguage` | Reads the global system language id |
| `26.9` | `kos.GetTimeCounter` | `kos.UptimeCentiseconds` | Returns uptime in 1/100 second |
| `26.10` | `kos.GetTimeCounterPro` | `kos.UptimeNanoseconds` | Returns uptime in nanoseconds |
| `37.0` | `kos.GetMouseScreenPosition` | `kos.MouseScreenPosition` | Mouse coordinates on the screen |
| `37.1` | `kos.GetMouseWindowPosition` | `kos.MouseWindowPosition` | Mouse coordinates relative to the window |
| `37.2` | `kos.GetMouseButtonState` | `kos.MouseHeldButtons` | Held mouse button states |
| `37.3` | `kos.GetMouseButtonEventState` | `kos.MouseButtons` | Mouse button states plus edge events |
| `37.4` | `kos.LoadCursorRaw` | `kos.LoadCursorCURData`, `kos.LoadCursorARGB` | Generic cursor loader for in-memory `.cur` data or indirect ARGB cursor images |
| `37.5` | `kos.SetCursorRaw` | `kos.SetCursor`, `kos.RestoreDefaultCursor` | Set or restore the current thread window cursor |
| `37.6` | `kos.DeleteCursorRaw` | `kos.DeleteCursor` | Delete a cursor previously loaded by the thread |
| `37.7` | `kos.GetMouseScrollData` | `kos.MouseScrollDelta` | Signed horizontal/vertical scroll deltas |
| `37.8` | `kos.LoadCursorWithEncoding` | `kos.LoadCursorFile`, `kos.LoadCursorFileWithEncoding` | Cursor file load with explicit path encoding |
| `38` | `kos.DrawLine` | `kos.StrokeLine` | Line draw |
| `40` | `kos.SetEventMask` | `kos.SwapEventMask` | Event filtering |
| `48.4` | `kos.GetSkinHeight` | `kos.SkinHeight` | Height of the skinned-window header |
| `48.5` | `kos.GetScreenWorkingArea` | `kos.ScreenWorkingArea` | Returns inclusive working-area bounds |
| `48.7` | `kos.GetSkinMarginsRaw` | `kos.WindowSkinMargins` | Returns header text margins for skinned windows |
| `48.8` | `kos.SetSkin` | `kos.SetSystemSkinLegacy` | Apply a system skin path using the legacy default-encoding contract |
| `48.13` | `kos.SetSkinWithEncoding` | `kos.SetSystemSkin`, `kos.SetSystemSkinWithEncoding` | Apply a system skin path with explicit string encoding |
| `60.1` | `kos.SetIPCArea` | `kos.RegisterIPCBuffer` | Registers the process receive buffer for event `7` delivery |
| `60.2` | `kos.SendIPCMessage` | `kos.SendIPCRaw`, `kos.SendIPC`, `kos.InspectIPCBuffer` | Non-allocating IPC helpers keep the sample within the current bootstrap runtime surface |
| `63` | `kos.DebugOutHex`, `kos.DebugOutChar`, `kos.DebugOutStr` | `kos.DebugHex`, `kos.DebugChar`, `kos.DebugString` | Debug board helpers |
| `68.12` | `malloc` | runtime-only | Lazy heap init via `68.11` |
| `68.13` | `free` | runtime-only | Lazy heap init via `68.11` |
| `68.20` | `realloc` | runtime-only | Lazy heap init via `68.11` |
| `70` | `kos.FileSystem` | - | Raw long-name file system interface |
| `71.1` | `kos.SetCaptionWithPrefix` | `kos.SetWindowTitleWithEncodingPrefix` | Uses the caption-prefix encoding contract |
| `71.2` | `kos.SetCaption` | `kos.SetWindowTitle` | Uses explicit UTF-8 encoding |
| `72` | `kos.SendMessage` | `kos.SendActiveWindowMessage`, `kos.SendActiveWindowKey`, `kos.SendActiveWindowButton` | Active-window key/button injection |
| `80` | `kos.FileSystemEncoded` | `kos.ReadFile`, `kos.ReadDirectory`, `kos.CreateOrRewriteFile`, `kos.WriteFile`, `kos.SetFileSize`, `kos.GetPathInfo`, `kos.SetPathInfo`, `kos.StartApplication`, `kos.DeletePath`, `kos.CreateDirectory`, `kos.RenamePath` | UTF-8-friendly file system layer |
| `-1` | `kos.Exit` | - | Terminate thread/process |

## Phase 1 Status

The bootstrap repo currently has no remaining priority gaps from the original
Phase 1 syscall inventory. Further syscall additions are now SDK expansion work
rather than initial audit cleanup.

## Notes

- The current raw layer intentionally stays thin; friendly Go wrappers live in
  `kos/*.go`.
- Window, text, and caption helpers use Go strings, but the asm/runtime glue is
  responsible for adapting them to KolibriOS string conventions.
