package assertionsdemo

import (
	"../../kos"
	"../../ui"
)

const (
	assertionsButtonExit   kos.ButtonID = 1
	assertionsButtonToggle kos.ButtonID = 2

	assertionsWindowX      = 360
	assertionsWindowY      = 220
	assertionsWindowWidth  = 540
	assertionsWindowHeight = 188
	assertionsWindowTitle  = "KolibriOS Assertions"
)

type demoApp struct {
	enabled bool
	message string
	toggle  ui.Button
}

func newDemoApp() *demoApp {
	toggle := ui.NewButton(assertionsButtonToggle, "Toggle", 28, 126)
	toggle.Width = 118

	app := &demoApp{
		enabled: true,
		toggle:  toggle,
	}
	app.rebuildMessage()
	return app
}

func (app *demoApp) draw() {
	exit := ui.NewButton(assertionsButtonExit, "Exit", 150, 126)
	exit.Width = 96

	kos.BeginRedraw()
	kos.OpenWindow(assertionsWindowX, assertionsWindowY, assertionsWindowWidth, assertionsWindowHeight, assertionsWindowTitle)
	kos.DrawText(28, 48, ui.White, "interface{} assertions and type switch")
	kos.DrawText(28, 74, ui.Silver, app.message)
	app.toggle.Draw()
	exit.Draw()
	kos.EndRedraw()
}

func (app *demoApp) handleButton(id kos.ButtonID) bool {
	switch id {
	case assertionsButtonToggle:
		app.enabled = !app.enabled
		app.rebuildMessage()
		app.draw()
	case assertionsButtonExit:
		kos.Exit()
		return true
	}

	return false
}

func (app *demoApp) loop() {
	for {
		switch kos.WaitEvent() {
		case kos.EventRedraw:
			app.draw()
		case kos.EventButton:
			if app.handleButton(kos.CurrentButtonID()) {
				return
			}
		}
	}
}

func (app *demoApp) rebuildMessage() {
	var value interface{}

	if app.enabled {
		value = "assert string"
	} else {
		value = 2026
	}

	text, ok := value.(string)
	if ok {
		forced := value.(string)
		app.message = forced + " / ok / " + describeValue(value, text)
		return
	}

	app.message = "not string / miss / " + describeValue(value, text)
}

func describeValue(value interface{}, text string) string {
	switch value.(type) {
	case string:
		return "switch string / " + text
	case int:
		return "switch int / 2026"
	default:
		return "switch default"
	}
}

func Run() {
	app := newDemoApp()
	app.loop()
}
