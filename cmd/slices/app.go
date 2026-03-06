package slicesdemo

import (
	"../../kos"
	"../../ui"
)

const (
	slicesButtonExit kos.ButtonID = 1
	slicesButtonToggle kos.ButtonID = 2

	slicesWindowX = 380
	slicesWindowY = 210
	slicesWindowWidth = 460
	slicesWindowHeight = 180
	slicesWindowTitle = "KolibriOS Slices"
)

type App struct {
	enabled bool
	message string
	toggle  ui.Button
}

func NewApp() App {
	toggle := ui.NewButton(slicesButtonToggle, "Toggle", 28, 118)
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
	case slicesButtonToggle:
		app.enabled = !app.enabled
		app.rebuildMessage()
		app.Redraw()
	case slicesButtonExit:
		kos.Exit()
		return true
	}

	return false
}

func (app *App) Redraw() {
	kos.BeginRedraw()
	kos.OpenWindow(slicesWindowX, slicesWindowY, slicesWindowWidth, slicesWindowHeight, slicesWindowTitle)
	kos.DrawText(28, 48, ui.White, "Slice bootstrap sample")
	kos.DrawText(28, 74, ui.Silver, "make([]byte), append, copy, []byte/string")
	kos.DrawText(28, 96, ui.Aqua, app.message)
	app.toggle.Draw()
	kos.EndRedraw()
}

func (app *App) rebuildMessage() {
	var base string

	if app.enabled {
		base = "slice mode on"
	} else {
		base = "slice mode off"
	}

	src := []byte(base)
	buf := make([]byte, 0, 2)
	buf = append(buf, src...)
	buf = append(buf, '!', '!')
	dst := make([]byte, len(buf))
	copy(dst, buf)

	app.message = string(dst)
}
