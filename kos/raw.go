package kos

// Low-level bindings are kept exported to preserve the original assembly ABI.
func Sleep(uint32)
func GetTime() uint32
func Event() int
func GetButtonID() int
func CreateButton(x int, y int, width int, height int, id int, color uint32)
func Exit()
func Redraw(mode int)
func Window(x int, y int, width int, height int, title string)
func WriteText(x int, y int, color uint32, text string)
func DrawLine(x1 int, y1 int, x2 int, y2 int, color uint32)
func DrawBar(x int, y int, width int, height int, color uint32)
func DebugOutHex(uint32)
func DebugOutChar(byte)
func DebugOutStr(string)

func Pointer2byteSlice(ptr uint32) *[]byte __asm__("__unsafe_get_addr")
