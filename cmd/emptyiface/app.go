package emptyifacedemo

import (
	"../../kos"
	"../../ui"
)

const (
	emptyIfaceButtonExit kos.ButtonID = 1
	emptyIfaceButtonToggle kos.ButtonID = 2

	emptyIfaceWindowX = 340
	emptyIfaceWindowY = 210
	emptyIfaceWindowWidth = 520
	emptyIfaceWindowHeight = 180
	emptyIfaceWindowTitle = "KolibriOS Empty Interface"
)

type App struct {
	enabled bool
	message string
	toggle  ui.Button
}

func NewApp() App {
	toggle := ui.NewButton(emptyIfaceButtonToggle, "Toggle", 28, 118)
	toggle.Width = 118

	app := App{
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
	case emptyIfaceButtonToggle:
		app.enabled = !app.enabled
		app.rebuildMessage()
		app.Redraw()
	case emptyIfaceButtonExit:
		kos.Exit()
		return true
	}

	return false
}

func (app *App) Redraw() {
	kos.BeginRedraw()
	kos.OpenWindow(emptyIfaceWindowX, emptyIfaceWindowY, emptyIfaceWindowWidth, emptyIfaceWindowHeight, emptyIfaceWindowTitle)
	kos.DrawText(28, 48, ui.White, "Empty interface bootstrap sample")
	kos.DrawText(28, 74, ui.Silver, "interface{} assignment and equality")
	kos.DrawText(28, 96, ui.Aqua, app.message)
	app.toggle.Draw()
	kos.EndRedraw()
}

func (app *App) rebuildMessage() {
	var left interface{}
	var right interface{}

	if app.enabled {
		left = "empty iface on"
		right = "empty iface on"
	} else {
		left = "empty iface off"
		right = "empty iface off"
	}

	if left == right {
		if app.enabled {
			app.message = "empty iface on / eq ok"
		} else {
			app.message = "empty iface off / eq ok"
		}
		return
	}

	app.message = "empty interface equality failed"
}
