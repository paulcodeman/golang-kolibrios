package kos

// Low-level bindings are kept exported to preserve the original assembly ABI.
// sysfuncs.txt is the source of truth for the function numbers and contracts.

// Function 5 - delay.
func Sleep(uint32)

// Function 3 - get system time.
func GetTime() uint32

// Function 10 - wait for event.
func Event() int

// Function 2 - get the code of the pressed key.
func GetKey() int

// Function 11 - check for event, no wait.
func CheckEvent() int

// Function 9 - information on execution thread.
func GetThreadInfo(buffer *byte, slot int) int

// Function 23 - wait for event with timeout.
func WaitEventTimeout(uint32) int

// Function 40 - set the mask for expected events.
func SetEventMask(uint32) uint32

// Function 18, subfunction 16 - get size of free RAM in kilobytes.
func GetFreeRAM() uint32

// Function 18, subfunction 17 - get total RAM in kilobytes.
func GetTotalRAM() uint32

// Function 17 - get the identifier of the pressed button.
func GetButtonID() int

// Function 8 - define/delete button.
func CreateButton(x int, y int, width int, height int, id int, color uint32)

// Function -1 - terminate thread/process.
func Exit()

// Function 12 - begin/end window redraw.
func Redraw(mode int)

// Function 0 - define and draw the window.
func Window(x int, y int, width int, height int, title string)

// Function 71, subfunction 2 - set window caption with explicit encoding.
func SetCaption(title string)

// Function 37, subfunction 0 - get screen coordinates of the mouse.
func GetMouseScreenPosition() uint32

// Function 37, subfunction 1 - get mouse coordinates relative to the window.
func GetMouseWindowPosition() uint32

// Function 37, subfunction 2 - get states of the mouse buttons.
func GetMouseButtonState() uint32

// Function 37, subfunction 3 - get states and events of the mouse buttons.
func GetMouseButtonEventState() uint32

// Function 37, subfunction 7 - get scroll data.
func GetMouseScrollData() uint32

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

func Pointer2byteSlice(ptr uint32) *[]byte __asm__("__unsafe_get_addr")
