package kos

func ScreenSize() (width int, height int) {
	packed := GetScreenSize()
	return int(packed>>16) + 1, int(packed&0xFFFF) + 1
}
