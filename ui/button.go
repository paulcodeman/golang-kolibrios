package ui

import "../kos"

const (
	defaultButtonWidthPadding = 20
	defaultButtonHeight = 30
	defaultCharWidth = 15
	defaultTextPaddingX = 10
	defaultTextPaddingY = 8
)

type Button struct {
	ID kos.ButtonID
	Label string
	X int
	Y int
	Width int
	Height int
	Background kos.Color
	Foreground kos.Color
	TextPaddingX int
	TextPaddingY int
}

func NewButton(id kos.ButtonID, label string, x int, y int) Button {
	return Button{
		ID: id,
		Label: label,
		X: x,
		Y: y,
		Width: len(label)*defaultCharWidth + defaultButtonWidthPadding,
		Height: defaultButtonHeight,
		Background: Blue,
		Foreground: White,
		TextPaddingX: defaultTextPaddingX,
		TextPaddingY: defaultTextPaddingY,
	}
}

func (button Button) Draw() {
	kos.DrawButton(
		button.X,
		button.Y,
		button.resolvedWidth(),
		button.resolvedHeight(),
		button.ID,
		button.Background,
	)
	kos.DrawText(
		button.X+button.resolvedTextPaddingX(),
		button.Y+button.resolvedTextPaddingY(),
		button.Foreground,
		button.Label,
	)
}

func (button Button) resolvedWidth() int {
	if button.Width > 0 {
		return button.Width
	}

	return len(button.Label)*defaultCharWidth + defaultButtonWidthPadding
}

func (button Button) resolvedHeight() int {
	if button.Height > 0 {
		return button.Height
	}

	return defaultButtonHeight
}

func (button Button) resolvedTextPaddingX() int {
	if button.TextPaddingX > 0 {
		return button.TextPaddingX
	}

	return defaultTextPaddingX
}

func (button Button) resolvedTextPaddingY() int {
	if button.TextPaddingY > 0 {
		return button.TextPaddingY
	}

	return defaultTextPaddingY
}
