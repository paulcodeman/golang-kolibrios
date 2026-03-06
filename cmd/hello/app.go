package hello

import (
	"../../kos"
	"../../ui"
)

const (
	helloExitButton kos.ButtonID = 1

	helloWindowX = 440
	helloWindowY = 220
	helloWindowWidth = 320
	helloWindowHeight = 140
	helloWindowTitle = "KolibriOS Hello"
)

type App struct{}

func NewApp() App {
	return App{}
}

func (app *App) Run() {
	for {
		switch kos.WaitEvent() {
		case kos.EventRedraw:
			app.Redraw()
		case kos.EventButton:
			if kos.CurrentButtonID() == helloExitButton {
				kos.Exit()
				return
			}
		}
	}
}

func (app *App) Redraw() {
	kos.BeginRedraw()
	kos.OpenWindow(helloWindowX, helloWindowY, helloWindowWidth, helloWindowHeight, helloWindowTitle)
	kos.DrawText(28, 52, ui.White, "Hello from Go on KolibriOS")
	kos.DrawText(28, 78, ui.Silver, "Bootstrap sample application")
	kos.EndRedraw()
}
