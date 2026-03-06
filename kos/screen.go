package kos

func ScreenSize() (width int, height int) {
	packed := GetScreenSize()
	return int(packed>>16) + 1, int(packed&0xFFFF) + 1
}

func ScreenWorkingArea() Rect {
	var vertical uint32

	horizontal := GetScreenWorkingArea(&vertical)
	return Rect{
		Left:   int(horizontal >> 16),
		Top:    int(vertical >> 16),
		Right:  int(horizontal & 0xFFFF),
		Bottom: int(vertical & 0xFFFF),
	}
}

func SkinHeight() int {
	return GetSkinHeight()
}
