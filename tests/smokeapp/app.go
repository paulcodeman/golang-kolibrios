package smokeapp

import (
	"errors"
	"os"

	"../../kos"
	"../../ui"
)

const (
	smokeButtonExit kos.ButtonID = 1

	smokeWindowTitle  = "KolibriOS Emulator Smoke"
	smokeWindowX      = 250
	smokeWindowY      = 150
	smokeWindowWidth  = 760
	smokeWindowHeight = 302
)

type sourceText interface {
	Text() string
}

type targetText interface {
	Text() string
}

type smokeText struct {
	text string
}

func (value smokeText) Text() string {
	return value.text
}

type smokeWrappedError struct {
	cause error
}

func (err smokeWrappedError) Error() string {
	return "wrapped"
}

func (err smokeWrappedError) Unwrap() error {
	return err.cause
}

type App struct {
	timeoutOK    bool
	timeOK       bool
	stringsOK    bool
	slicesOK     bool
	ifaceOK      bool
	assertionsOK bool
	errorsOK     bool
	systemOK     bool
	shutdownOK   bool
	summary      string
	clockLine    string
	timeoutLine  string
	stringsLine  string
	slicesLine   string
	ifaceLine    string
	assertLine   string
	errorsLine   string
	systemLine   string
	powerLine    string
}

func NewApp() App {
	app := App{}
	app.runChecks()
	return app
}

func (app *App) Run() {
	app.waitForTimeout()
	app.tryPowerOff()
	if app.shutdownOK {
		return
	}

	app.Redraw()
	app.eventLoop()
}

func (app *App) eventLoop() {
	for {
		switch kos.WaitEvent() {
		case kos.EventRedraw:
			app.Redraw()
		case kos.EventButton:
			if kos.CurrentButtonID() == smokeButtonExit {
				os.Exit(0)
				return
			}
		}
	}
}

func (app *App) Redraw() {
	exit := ui.NewButton(smokeButtonExit, "Exit", 632, 258)
	exit.Width = 96

	kos.BeginRedraw()
	kos.OpenWindow(smokeWindowX, smokeWindowY, smokeWindowWidth, smokeWindowHeight, smokeWindowTitle)
	kos.DrawText(28, 44, app.statusColor(app.allOK()), app.summary)
	kos.DrawText(28, 68, ui.Silver, "headless smoke uses guest poweroff as the host-side success signal")
	kos.DrawText(28, 94, app.statusColor(app.timeOK), app.clockLine)
	kos.DrawText(28, 114, app.statusColor(app.timeoutOK), app.timeoutLine)
	kos.DrawText(28, 134, app.statusColor(app.stringsOK), app.stringsLine)
	kos.DrawText(28, 154, app.statusColor(app.slicesOK), app.slicesLine)
	kos.DrawText(28, 174, app.statusColor(app.ifaceOK), app.ifaceLine)
	kos.DrawText(28, 194, app.statusColor(app.assertionsOK), app.assertLine)
	kos.DrawText(28, 214, app.statusColor(app.errorsOK), app.errorsLine)
	kos.DrawText(28, 234, app.statusColor(app.systemOK), app.systemLine)
	kos.DrawText(28, 256, app.statusColor(app.shutdownOK), app.powerLine)
	exit.Draw()
	kos.EndRedraw()
}

func (app *App) runChecks() {
	app.summary = "emulator smoke checks running"
	app.clockLine = "clock : pending"
	app.timeoutLine = "event : pending"
	app.stringsLine = "text  : pending"
	app.slicesLine = "slice : pending"
	app.ifaceLine = "iface : pending"
	app.assertLine = "assert: pending"
	app.errorsLine = "error : pending"
	app.systemLine = "sys   : pending"
	app.powerLine = "power : pending"

	app.timeOK, app.clockLine = checkClock()
	app.stringsOK, app.stringsLine = checkStrings()
	app.slicesOK, app.slicesLine = checkSlices()
	app.ifaceOK, app.ifaceLine = checkInterfaces()
	app.assertionsOK, app.assertLine = checkAssertions()
	app.errorsOK, app.errorsLine = checkErrors()
	app.systemOK, app.systemLine = checkSystemSurface()
	app.powerLine = "power : waiting for timeout gate"

	if app.allOK() {
		app.summary = "emulator smoke checks passed"
		return
	}

	app.summary = "emulator smoke checks failed"
}

func (app *App) waitForTimeout() {
	start := kos.UptimeCentiseconds()
	event := kos.WaitEventFor(5)
	end := kos.UptimeCentiseconds()
	if end < start {
		app.timeoutLine = "event : FAIL / timed wait regressed uptime"
		return
	}

	app.timeoutOK = true
	if event == kos.EventNone {
		app.timeoutLine = "event : PASS / WaitEventFor returned idle"
	} else {
		app.timeoutLine = "event : PASS / WaitEventFor observed event " + formatInt(int(event))
	}

	if app.allOK() {
		app.summary = "emulator smoke checks passed"
	} else {
		app.summary = "emulator smoke checks failed"
	}
}

func (app *App) tryPowerOff() {
	if !app.allOK() {
		return
	}

	kos.SleepCentiseconds(2)
	app.shutdownOK = kos.PowerOff()
	if app.shutdownOK {
		return
	}

	app.powerLine = "power : FAIL / system shutdown rejected"
	app.summary = "emulator smoke checks passed but poweroff failed"
	kos.DebugString("smoke: poweroff rejected\r\n")
}

func (app *App) allOK() bool {
	return app.timeoutOK &&
		app.timeOK &&
		app.stringsOK &&
		app.slicesOK &&
		app.ifaceOK &&
		app.assertionsOK &&
		app.errorsOK &&
		app.systemOK
}

func (app *App) statusColor(ok bool) kos.Color {
	if ok {
		return ui.Lime
	}

	return ui.Red
}

func checkClock() (bool, string) {
	start := kos.UptimeCentiseconds()
	kos.SleepCentiseconds(1)
	end := kos.UptimeCentiseconds()
	clock := formatClock(kos.SystemTime())
	if end >= start {
		return true, "clock : PASS / " + clock + " / uptime " + formatUint32(start) + " -> " + formatUint32(end)
	}

	return false, "clock : FAIL / uptime regressed around " + clock
}

func checkStrings() (bool, string) {
	message := "smoke" + " strings"
	if message == "smoke strings" {
		return true, "text  : PASS / string concat and equality"
	}

	return false, "text  : FAIL / string runtime mismatch"
}

func checkSlices() (bool, string) {
	buf := make([]byte, 0, 4)
	buf = append(buf, []byte("go")...)
	buf = append(buf, '4')
	out := make([]byte, len(buf))
	copied := copy(out, buf)
	if copied == 3 && string(out) == "go4" {
		return true, "slice : PASS / append copy convert"
	}

	return false, "slice : FAIL / append or copy mismatch"
}

func checkInterfaces() (bool, string) {
	var left sourceText = smokeText{text: "iface smoke"}
	var right sourceText = smokeText{text: "iface smoke"}
	if left == right && left.Text() == "iface smoke" {
		return true, "iface : PASS / dispatch and equality"
	}

	return false, "iface : FAIL / interface runtime mismatch"
}

func checkAssertions() (bool, string) {
	var any interface{} = smokeText{text: "assert smoke"}
	var iface sourceText = smokeText{text: "bridge smoke"}

	direct, okDirect := any.(smokeText)
	bridge, okAny := any.(targetText)
	converted, okIface := iface.(targetText)
	switch describeAssertionValue(any) {
	case "switch smoke":
		if okDirect &&
			okAny &&
			okIface &&
			direct.Text() == "assert smoke" &&
			bridge.Text() == "assert smoke" &&
			converted.Text() == "bridge smoke" {
			return true, "assert: PASS / e2t e2i i2i switch"
		}
	}

	return false, "assert: FAIL / assertion runtime mismatch"
}

func checkErrors() (bool, string) {
	sentinel := errors.New("smoke sentinel")
	wrapped := smokeWrappedError{cause: sentinel}

	if !errors.Is(sentinel, sentinel) {
		return false, "error : FAIL / sentinel self-match"
	}

	if errors.Unwrap(wrapped) != sentinel {
		return false, "error : FAIL / Unwrap lost wrapped cause"
	}

	if !errors.Is(wrapped, sentinel) {
		return false, "error : FAIL / Is failed through Unwrap"
	}

	return true, "error : PASS / ordinary import errors"
}

func checkSystemSurface() (bool, string) {
	version := kos.KernelVersion()
	screenWidth, screenHeight := kos.ScreenSize()
	workArea := kos.ScreenWorkingArea()
	skinHeight := kos.SkinHeight()
	margins := kos.WindowSkinMargins()

	if screenWidth <= 0 || screenHeight <= 0 {
		return false, "sys   : FAIL / invalid screen size"
	}

	if workArea.Width() <= 0 || workArea.Height() <= 0 {
		return false, "sys   : FAIL / invalid work area"
	}

	if workArea.Left < 0 || workArea.Top < 0 || workArea.Right >= screenWidth || workArea.Bottom >= screenHeight {
		return false, "sys   : FAIL / work area outside screen"
	}

	if skinHeight < 0 || margins.Left < 0 || margins.Top < 0 || margins.Right < 0 || margins.Bottom < 0 {
		return false, "sys   : FAIL / invalid skin geometry"
	}

	return true, "sys   : PASS / v " + formatKernelVersion(version) +
		" / " + formatInt(screenWidth) + "x" + formatInt(screenHeight) +
		" / work " + formatInt(workArea.Width()) + "x" + formatInt(workArea.Height()) +
		" / skin " + formatInt(skinHeight)
}

func describeAssertionValue(value interface{}) string {
	switch value.(type) {
	case smokeText:
		return "switch smoke"
	default:
		return "switch default"
	}
}

func Run() {
	app := NewApp()
	app.Run()
}
