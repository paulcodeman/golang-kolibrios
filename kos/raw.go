package kos

// Low-level bindings are kept exported to preserve the original assembly ABI.
// sysfuncs.txt is the source of truth for the function numbers and contracts.

// Function 5 - delay.
func Sleep(uint32)

// Function 3 - get system time.
func GetTime() uint32

// Function 29 - get system date.
func GetDate() uint32

// Function 26, subfunction 9 - get uptime counter in 1/100 second.
func GetTimeCounter() uint32

// Function 26, subfunction 10 - get high precision uptime counter in nanoseconds.
func GetTimeCounterPro() uint64

// Function 10 - wait for event.
func Event() int

// Function 2 - get the code of the pressed key.
func GetKey() int

// Function 21, subfunction 2 - set one of the keyboard layout tables.
func SetKeyboardLayoutRaw(which int, table *byte) int

// Function 21, subfunction 2 - set the global keyboard layout language id.
func SetKeyboardLanguageRaw(language int) int

// Function 21, subfunction 5 - set the global system language id.
func SetSystemLanguageRaw(language int) int

// Function 26, subfunction 2 - get one of the keyboard layout tables.
func GetKeyboardLayoutRaw(which int, buffer *byte) int

// Function 26, subfunction 2 - get the global keyboard layout language id.
func GetKeyboardLanguageRaw() int

// Function 26, subfunction 5 - get the global system language id.
func GetSystemLanguageRaw() int

// Function 11 - check for event, no wait.
func CheckEvent() int

// Function 9 - information on execution thread.
func GetThreadInfo(buffer *byte, slot int) int

// Function 23 - wait for event with timeout.
func WaitEventTimeout(uint32) int

// Function 40 - set the mask for expected events.
func SetEventMask(uint32) uint32

// Function 46 - reserve/free a group of I/O ports.
func SetPortsRaw(mode int, start uint32, end uint32) int

// Function 60, subfunction 1 - register the IPC receive area.
func SetIPCArea(buffer *byte, size uint32) uint32

// Function 60, subfunction 2 - send an IPC message to a PID/TID.
func SendIPCMessage(pid uint32, data *byte, size uint32) uint32

// Function 18, subfunction 3 - make active the window of the given thread slot.
func FocusWindowBySlot(int)

// Function 18, subfunction 7 - get the slot number of the active window.
func GetActiveWindowSlotRaw() int

// Function 48, subfunction 4 - get skinned-window header height.
func GetSkinHeight() int

// Function 48, subfunction 7 - get skin margins for header text layout.
func GetSkinMarginsRaw(vertical *uint32) uint32

// Function 48, subfunction 8 - set the current skin using the default encoding path contract.
func SetSkin(path string) uint32

// Function 48, subfunction 13 - set the current skin using an explicit path encoding.
func SetSkinWithEncoding(encoding StringEncoding, path string) uint32

// Function 48, subfunction 5 - get packed screen working-area coordinates.
func GetScreenWorkingArea(bottom *uint32) uint32

// Function 18, subfunction 13 - get kernel version metadata.
func GetKernelVersion(buffer *byte)

// Function 18, subfunction 9 - system shutdown with a mode parameter.
func SystemShutdown(uint32) uint32

// Function 18, subfunction 16 - get size of free RAM in kilobytes.
func GetFreeRAM() uint32

// Function 18, subfunction 17 - get total RAM in kilobytes.
func GetTotalRAM() uint32

// Function 68, subfunction 18 - load DLL with explicit path encoding.
func LoadDLLWithEncoding(encoding StringEncoding, path string) uint32

// Function 68, subfunction 19 - load DLL using the legacy/default path contract.
func LoadDLL(path string) uint32

// Runtime helper - resolve a function pointer from a DLL export table.
func LookupDLLExportRaw(table uint32, name *byte) uint32 __asm__("runtime_kos_lookup_dll_export")

// Runtime helper - invoke a stdcall function pointer with 0 arguments.
func CallStdcall0Raw(proc uint32) uint32 __asm__("runtime_kos_call_stdcall0")

// Runtime helper - invoke a stdcall function pointer with 1 argument.
func CallStdcall1Raw(proc uint32, arg0 uint32) uint32 __asm__("runtime_kos_call_stdcall1")

// Runtime helper - invoke a stdcall function pointer with 2 arguments.
func CallStdcall2Raw(proc uint32, arg0 uint32, arg1 uint32) uint32 __asm__("runtime_kos_call_stdcall2")

// Runtime helper - invoke a stdcall function pointer with 1 argument and no return value.
func CallStdcall1VoidRaw(proc uint32, arg0 uint32) __asm__("runtime_kos_call_stdcall1_void")

// Runtime helper - invoke a stdcall function pointer with 2 arguments and no return value.
func CallStdcall2VoidRaw(proc uint32, arg0 uint32, arg1 uint32) __asm__("runtime_kos_call_stdcall2_void")

// Runtime helper - invoke a stdcall function pointer with 5 arguments and no return value.
func CallStdcall5VoidRaw(proc uint32, arg0 uint32, arg1 uint32, arg2 uint32, arg3 uint32, arg4 uint32) __asm__("runtime_kos_call_stdcall5_void")

// Runtime helper - check whether a shared active console bridge is registered.
func ConsoleBridgeReadyRaw() uint32 __asm__("runtime_console_bridge_ready")

// Runtime helper - register active console write/exit procedures in shared runtime state.
func ConsoleBridgeSetRaw(table uint32, writeProc uint32, exitProc uint32, getsProc uint32) __asm__("runtime_console_bridge_set")

// Runtime helper - clear the shared active console bridge if the table matches.
func ConsoleBridgeClearRaw(table uint32) __asm__("runtime_console_bridge_clear")

// Runtime helper - write directly through the shared active console bridge.
func ConsoleBridgeWriteRaw(data uint32, size uint32) uint32 __asm__("runtime_console_bridge_write")

// Runtime helper - read one line through the shared active console bridge.
func ConsoleBridgeReadLineRaw(data uint32, size uint32) uint32 __asm__("runtime_console_bridge_read_line")

// Runtime helper - close the shared active console bridge and clear it.
func ConsoleBridgeCloseRaw(closeWindow uint32) __asm__("runtime_console_bridge_close")

// Function 70 - file system interface with long names support.
func FileSystem(request *FileSystemRequest, secondary *uint32) int

// Function 80 - file system interface with parameter of encoding.
func FileSystemEncoded(request *EncodedFileSystemRequest, secondary *uint32) int

// Function 77, subfunction 10 - read from a file handle.
// The current kernel contract documents pipe descriptors on this path.
func PosixReadRaw(fd uint32, buffer *byte, size uint32) int

// Function 77, subfunction 11 - write to a file handle.
// The current kernel contract documents pipe descriptors on this path.
func PosixWriteRaw(fd uint32, buffer *byte, size uint32) int

// Function 77, subfunction 13 - create a pipe and return two file handles.
func PosixPipe2Raw(pipefd *uint32, flags uint32) int

// Function 30, subfunction 5 - get current folder with explicit encoding.
func GetCurrentFolderRaw(buffer *byte, size uint32, encoding StringEncoding) int

// Function 17 - get the identifier of the pressed button.
func GetButtonID() int

// Function 8 - define/delete button.
func CreateButton(x int, y int, width int, height int, id int, color uint32)

// Function -1 - terminate thread/process.
func ExitRaw()

// Function 12 - begin/end window redraw.
func Redraw(mode int)

// Function 0 - define and draw the window.
func Window(x int, y int, width int, height int, title string)

// Function 71, subfunction 2 - set window caption with explicit encoding.
func SetCaption(title string)

// Function 71, subfunction 1 - set window caption using an inline encoding prefix.
func SetCaptionWithPrefix(encoding StringEncoding, title string)

// Function 72, subfunction 1 - send a key or button event to the active window.
func SendMessage(event int, param uint32) int

// Function 37, subfunction 0 - get screen coordinates of the mouse.
func GetMouseScreenPosition() uint32

// Function 37, subfunction 1 - get mouse coordinates relative to the window.
func GetMouseWindowPosition() uint32

// Function 37, subfunction 2 - get states of the mouse buttons.
func GetMouseButtonState() uint32

// Function 37, subfunction 3 - get states and events of the mouse buttons.
func GetMouseButtonEventState() uint32

// Function 37, subfunction 4 - load a cursor from memory/file descriptor arguments.
func LoadCursorRaw(data uint32, descriptor uint32) uint32

// Function 37, subfunction 5 - set the current thread window cursor.
func SetCursorRaw(handle uint32) uint32

// Function 37, subfunction 6 - delete a cursor previously loaded by the thread.
func DeleteCursorRaw(handle uint32)

// Function 37, subfunction 7 - get scroll data.
func GetMouseScrollData() uint32

// Function 37, subfunction 8 - load a cursor from a path with explicit string encoding.
func LoadCursorWithEncoding(encoding StringEncoding, path string) uint32

// Function 4 - draw text string.
func WriteText(x int, y int, color uint32, text string)

// Function 38 - draw line.
func DrawLine(x1 int, y1 int, x2 int, y2 int, color uint32)

// Function 13 - draw rectangle.
func DrawBar(x int, y int, width int, height int, color uint32)

// Function 14 - get screen size.
func GetScreenSize() uint32

// Function 63 - work with the debug board, write byte helper.
func DebugOutHex(uint32)

// Function 63 - work with the debug board, write byte helper.
func DebugOutChar(byte)

// Function 63 - work with the debug board, write string helper.
func DebugOutStr(string)

// Direct OUT instruction helper for previously reserved ports.
func PortWriteByteRaw(port uint32, value byte)

func Pointer2byteSlice(ptr uint32) *[]byte __asm__("__unsafe_get_addr")
