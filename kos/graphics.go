package kos

const (
	redrawBegin = 1
	redrawEnd = 2
)

func BeginRedraw() {
	Redraw(redrawBegin)
}

func EndRedraw() {
	Redraw(redrawEnd)
}

func OpenWindow(x int, y int, width int, height int, title string) {
	Window(x, y, width, height, title)
}

func SetWindowTitle(title string) {
	SetCaption(title)
}

func SetWindowTitleWithEncodingPrefix(encoding StringEncoding, title string) {
	if encoding == EncodingDefault {
		SetCaption(title)
		return
	}

	SetCaptionWithPrefix(encoding, title)
}

func DrawText(x int, y int, color Color, text string) {
	WriteText(x, y, uint32(color), text)
}

func FillRect(x int, y int, width int, height int, color Color) {
	DrawBar(x, y, width, height, uint32(color))
}

func StrokeLine(x1 int, y1 int, x2 int, y2 int, color Color) {
	DrawLine(x1, y1, x2, y2, uint32(color))
}

func DrawButton(x int, y int, width int, height int, id ButtonID, color Color) {
	CreateButton(x, y, width, height, int(id), uint32(color))
}
