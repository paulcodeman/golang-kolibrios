package runtimecheckdemo

import (
	"../../kos"
	"../../ui"
)

const (
	runtimeCheckButtonExit    kos.ButtonID = 1
	runtimeCheckButtonRecheck kos.ButtonID = 2

	runtimeCheckWindowX      = 280
	runtimeCheckWindowY      = 180
	runtimeCheckWindowWidth  = 720
	runtimeCheckWindowHeight = 260
	runtimeCheckWindowTitle  = "KolibriOS Runtime Check"
)

type sourceText interface {
	Text() string
}

type targetText interface {
	Text() string
}

type onText struct {
	text string
}

func (value onText) Text() string {
	return value.text
}

type offText struct {
	text string
}

func (value offText) Text() string {
	return value.text
}

type bridgeText struct {
	text string
}

func (value bridgeText) Text() string {
	return value.text
}

type App struct {
	enabled        bool
	stringsOK      bool
	slicesOK       bool
	ifaceOK        bool
	emptyIfaceOK   bool
	assertionsOK   bool
	summary        string
	stringsLine    string
	slicesLine     string
	ifaceLine      string
	emptyIfaceLine string
	assertionsLine string
	recheck        ui.Button
}

func NewApp() App {
	recheck := ui.NewButton(runtimeCheckButtonRecheck, "Recheck", 28, 206)
	recheck.Width = 132

	app := App{
		enabled: true,
		recheck: recheck,
	}
	app.runChecks()
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
	case runtimeCheckButtonRecheck:
		app.enabled = !app.enabled
		app.runChecks()
		app.Redraw()
	case runtimeCheckButtonExit:
		kos.Exit()
		return true
	}

	return false
}

func (app *App) Redraw() {
	exit := ui.NewButton(runtimeCheckButtonExit, "Exit", 182, 206)
	exit.Width = 96

	kos.BeginRedraw()
	kos.OpenWindow(runtimeCheckWindowX, runtimeCheckWindowY, runtimeCheckWindowWidth, runtimeCheckWindowHeight, runtimeCheckWindowTitle)
	kos.DrawText(28, 48, app.summaryColor(), app.summary)
	kos.DrawText(28, 74, ui.Silver, "recheck toggles the string/int assertion branch and reruns the runtime smoke set")
	kos.DrawText(28, 104, app.statusColor(app.stringsOK), app.stringsLine)
	kos.DrawText(28, 124, app.statusColor(app.slicesOK), app.slicesLine)
	kos.DrawText(28, 144, app.statusColor(app.ifaceOK), app.ifaceLine)
	kos.DrawText(28, 164, app.statusColor(app.emptyIfaceOK), app.emptyIfaceLine)
	kos.DrawText(28, 184, app.statusColor(app.assertionsOK), app.assertionsLine)
	app.recheck.Draw()
	exit.Draw()
	kos.EndRedraw()
}

func (app *App) runChecks() {
	app.stringsOK, app.stringsLine = checkStrings(app.enabled)
	app.slicesOK, app.slicesLine = checkSlices(app.enabled)
	app.ifaceOK, app.ifaceLine = checkInterfaces(app.enabled)
	app.emptyIfaceOK, app.emptyIfaceLine = checkEmptyInterface(app.enabled)
	app.assertionsOK, app.assertionsLine = checkAssertions(app.enabled)

	if app.allOK() {
		if app.enabled {
			app.summary = "runtime ok / string assertion branch"
		} else {
			app.summary = "runtime ok / int assertion branch"
		}
		return
	}

	if app.enabled {
		app.summary = "runtime failure / string assertion branch"
	} else {
		app.summary = "runtime failure / int assertion branch"
	}
}

func (app *App) allOK() bool {
	return app.stringsOK &&
		app.slicesOK &&
		app.ifaceOK &&
		app.emptyIfaceOK &&
		app.assertionsOK
}

func (app *App) summaryColor() kos.Color {
	return app.statusColor(app.allOK())
}

func (app *App) statusColor(ok bool) kos.Color {
	if ok {
		return ui.Lime
	}

	return ui.Red
}

func checkStrings(enabled bool) (bool, string) {
	base := "runtime off"
	if enabled {
		base = "runtime on"
	}

	message := base + " / strings"
	expected := base + " / strings"
	if message == expected {
		return true, "strings: PASS / " + message
	}

	return false, "strings: FAIL / equality mismatch"
}

func checkSlices(enabled bool) (bool, string) {
	base := "slice off"
	if enabled {
		base = "slice on"
	}

	src := []byte(base)
	buf := make([]byte, 0, 2)
	buf = append(buf, src...)
	if enabled {
		buf = append(buf, '!')
	} else {
		buf = append(buf, '.')
	}

	out := make([]byte, len(buf))
	copied := copy(out, buf)
	text := string(out)
	expected := base + "."
	if enabled {
		expected = base + "!"
	}

	if copied == len(buf) && text == expected {
		return true, "slices : PASS / " + text
	}

	return false, "slices : FAIL / copy or growth mismatch"
}

func checkInterfaces(enabled bool) (bool, string) {
	var source sourceText
	var mirror sourceText

	if enabled {
		source = onText{text: "iface on"}
		mirror = onText{text: "iface on"}
	} else {
		source = offText{text: "iface off"}
		mirror = offText{text: "iface off"}
	}

	if source == mirror {
		return true, "iface  : PASS / " + source.Text() + " / eq ok"
	}

	return false, "iface  : FAIL / dispatch or equality mismatch"
}

func checkEmptyInterface(enabled bool) (bool, string) {
	var left interface{}
	var right interface{}

	if enabled {
		left = "empty on"
		right = "empty on"
		if left == right {
			return true, "empty  : PASS / string eq ok"
		}
		return false, "empty  : FAIL / string equality mismatch"
	}

	left = 2026
	right = 2026
	if left == right {
		return true, "empty  : PASS / int eq ok"
	}

	return false, "empty  : FAIL / int equality mismatch"
}

func checkAssertions(enabled bool) (bool, string) {
	var directValue interface{}
	var candidate interface{}
	var anySource interface{}
	var ifaceSource sourceText
	var candidateText string
	var okString bool
	var switchPart string

	directValue = "direct string"
	anySource = bridgeText{text: "empty->iface"}
	ifaceSource = bridgeText{text: "iface->iface"}

	forcedString := directValue.(string)
	forcedAny := anySource.(targetText)
	anyTarget, okAny := anySource.(targetText)
	forcedIface := ifaceSource.(targetText)
	ifaceTarget, okIface := ifaceSource.(targetText)

	if enabled {
		candidate = "switch string"
		candidateText, okString = candidate.(string)
		switchPart = describeAssertionValue(candidate)
		if forcedString == "direct string" &&
			okString &&
			candidateText == "switch string" &&
			forcedAny.Text() == "empty->iface" &&
			okAny &&
			anyTarget.Text() == "empty->iface" &&
			forcedIface.Text() == "iface->iface" &&
			okIface &&
			ifaceTarget.Text() == "iface->iface" &&
			switchPart == "switch string" {
			return true, "assert : PASS / direct / e2i / i2i / string ok / " + switchPart
		}

		return false, "assert : FAIL / string assertion branch mismatch"
	}

	candidate = 2026
	_, okString = candidate.(string)
	switchPart = describeAssertionValue(candidate)
	_, missIface := candidate.(targetText)
	if forcedString == "direct string" &&
		!okString &&
		!missIface &&
		forcedAny.Text() == "empty->iface" &&
		okAny &&
		anyTarget.Text() == "empty->iface" &&
		forcedIface.Text() == "iface->iface" &&
		okIface &&
		ifaceTarget.Text() == "iface->iface" &&
		switchPart == "switch int" {
		return true, "assert : PASS / direct / e2i / i2i / string miss / " + switchPart
	}

	return false, "assert : FAIL / int assertion branch mismatch"
}

func describeAssertionValue(value interface{}) string {
	switch value.(type) {
	case string:
		return "switch string"
	case int:
		return "switch int"
	default:
		return "switch default"
	}
}

func Run() {
	app := NewApp()
	app.Run()
}
