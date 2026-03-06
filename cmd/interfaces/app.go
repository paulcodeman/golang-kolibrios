package interfacesdemo

import (
	"../../kos"
	"../../ui"
)

const (
	interfacesButtonExit kos.ButtonID = 1
	interfacesButtonToggle kos.ButtonID = 2

	interfacesWindowX = 360
	interfacesWindowY = 210
	interfacesWindowWidth = 500
	interfacesWindowHeight = 180
	interfacesWindowTitle = "KolibriOS Interfaces"
)

type MessageSource interface {
	Text() string
}

type OnSource struct {
	text string
}

func (source OnSource) Text() string {
	return source.text
}

type OffSource struct {
	text string
}

func (source OffSource) Text() string {
	return source.text
}

type App struct {
	enabled bool
	message string
	toggle  ui.Button
}

func NewApp() App {
	toggle := ui.NewButton(interfacesButtonToggle, "Toggle", 28, 118)
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
	case interfacesButtonToggle:
		app.enabled = !app.enabled
		app.rebuildMessage()
		app.Redraw()
	case interfacesButtonExit:
		kos.Exit()
		return true
	}

	return false
}

func (app *App) Redraw() {
	kos.BeginRedraw()
	kos.OpenWindow(interfacesWindowX, interfacesWindowY, interfacesWindowWidth, interfacesWindowHeight, interfacesWindowTitle)
	kos.DrawText(28, 48, ui.White, "Interface bootstrap sample")
	kos.DrawText(28, 74, ui.Silver, "non-empty interface dispatch and equality")
	kos.DrawText(28, 96, ui.Aqua, app.message)
	app.toggle.Draw()
	kos.EndRedraw()
}

func (app *App) rebuildMessage() {
	var source MessageSource
	var mirror MessageSource

	if app.enabled {
		source = OnSource{text: "interface mode on"}
		mirror = OnSource{text: "interface mode on"}
	} else {
		source = OffSource{text: "interface mode off"}
		mirror = OffSource{text: "interface mode off"}
	}

	if source == mirror {
		app.message = source.Text() + " / eq ok"
		return
	}

	app.message = "interface equality failed"
}
