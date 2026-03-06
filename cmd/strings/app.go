package stringsdemo

import (
	"../../kos"
	"../../ui"
)

const (
	buttonExit kos.ButtonID = 1
	buttonToggle kos.ButtonID = 2

	windowX = 420
	windowY = 210
	windowWidth = 420
	windowHeight = 170
	windowTitle = "KolibriOS Strings"
)

type App struct {
	state string
	message string
	toggle ui.Button
}

func NewApp() App {
	toggle := ui.NewButton(buttonToggle, "Toggle", 28, 112)
	toggle.Width = 110

	app := App{
		state: "ready",
		toggle: toggle,
	}
	app.rebuildMessage()

	return app
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
	case buttonToggle:
		if app.state == "ready" {
			app.state = "updated"
		} else {
			app.state = "ready"
		}
		app.rebuildMessage()
		app.Redraw()
	case buttonExit:
		kos.Exit()
		return true
	}

	return false
}

func (app *App) Redraw() {
	kos.BeginRedraw()
	kos.OpenWindow(windowX, windowY, windowWidth, windowHeight, windowTitle)
	kos.DrawText(28, 52, ui.White, "String bootstrap sample")
	kos.DrawText(28, 78, ui.Aqua, app.message)
	app.toggle.Draw()
	kos.EndRedraw()
}

func (app *App) rebuildMessage() {
	app.message = "KolibriOS " + "Go " + app.state + " message"
}
