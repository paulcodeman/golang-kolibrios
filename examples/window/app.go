package windowdemo

import (
	"os"

	"../../kos"
	"../../ui"
)

const (
	buttonExit      kos.ButtonID = 1
	buttonMoveLeft  kos.ButtonID = 2
	buttonMoveRight kos.ButtonID = 3

	windowX      = 500
	windowY      = 250
	windowWidth  = 420
	windowHeight = 200
	windowTitle  = "KolibriOS Window Demo"

	guideLineLeft  = 32
	guideLineRight = 150
	guideLineY     = 80

	barWidth    = 100
	barHeight   = 30
	barY        = 90
	barStep     = 32
	initialBarX = 160
	minBarX     = 32
	maxBarX     = 288
)

type App struct {
	barX      int
	moveLeft  ui.Button
	moveRight ui.Button
}

func NewApp() App {
	leftButton := ui.NewButton(buttonMoveLeft, "<-", 32, 128)
	leftButton.Width = 60

	rightButton := ui.NewButton(buttonMoveRight, "->", 310, 128)
	rightButton.Width = 60

	return App{
		barX:      initialBarX,
		moveLeft:  leftButton,
		moveRight: rightButton,
	}
}

func (app *App) Run() {
	for {
		switch kos.WaitEvent() {
		case kos.EventRedraw:
			app.Redraw()
		case kos.EventButton:
			if app.handleButton(kos.CurrentButtonID()) {
				return
			}
		}
	}
}

func (app *App) handleButton(id kos.ButtonID) bool {
	switch id {
	case buttonMoveLeft:
		app.moveBar(-barStep)
		app.Redraw()
	case buttonMoveRight:
		app.moveBar(barStep)
		app.Redraw()
	case buttonExit:
		os.Exit(0)
		return true
	}

	return false
}

func (app *App) Redraw() {
	kos.BeginRedraw()
	kos.OpenWindow(windowX, windowY, windowWidth, windowHeight, windowTitle)
	kos.StrokeLine(guideLineLeft, guideLineY, guideLineRight, guideLineY, ui.Green)
	kos.FillRect(app.barX, barY, barWidth, barHeight, ui.Red)
	app.moveLeft.Draw()
	app.moveRight.Draw()
	kos.EndRedraw()
}

func (app *App) moveBar(delta int) {
	next := app.barX + delta
	if next < minBarX {
		next = minBarX
	}
	if next > maxBarX {
		next = maxBarX
	}

	app.barX = next
}
