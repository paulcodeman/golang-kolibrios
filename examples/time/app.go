package timedemo

import (
	"../../kos"
	"../../ui"
)

const (
	timeprobeButtonExit    kos.ButtonID = 1
	timeprobeButtonSleep   kos.ButtonID = 2
	timeprobeButtonRefresh kos.ButtonID = 3

	timeprobeWindowX      = 320
	timeprobeWindowY      = 180
	timeprobeWindowWidth  = 650
	timeprobeWindowHeight = 245
	timeprobeWindowTitle  = "KolibriOS Time Demo"
)

type App struct {
	clock        kos.ClockTime
	uptime       uint32
	uptimeNS     uint64
	timeoutTicks uint32
	sleepDelta   uint32
	lastEvent    string
	sleep        ui.Button
	refresh      ui.Button
}

func NewApp() App {
	sleep := ui.NewButton(timeprobeButtonSleep, "Sleep 0.5s", 28, 191)
	sleep.Width = 142

	refresh := ui.NewButton(timeprobeButtonRefresh, "Refresh", 190, 191)
	refresh.Width = 112

	app := App{
		sleep:   sleep,
		refresh: refresh,
	}
	app.refreshTimeState()
	app.lastEvent = "startup refresh"

	return app
}

func (app *App) Run() {
	for {
		switch kos.WaitEventFor(50) {
		case kos.EventNone:
			app.timeoutTicks++
			app.refreshTimeState()
			app.lastEvent = "wait timeout / auto refresh"
			app.Redraw()
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
	case timeprobeButtonSleep:
		before := kos.UptimeCentiseconds()
		kos.SleepMilliseconds(500)
		after := kos.UptimeCentiseconds()
		app.sleepDelta = after - before
		app.refreshTimeState()
		app.lastEvent = "sleep delta / " + formatUint32(app.sleepDelta) + " cs"
		app.Redraw()
	case timeprobeButtonRefresh:
		app.refreshTimeState()
		app.lastEvent = "manual refresh"
		app.Redraw()
	case timeprobeButtonExit:
		kos.Exit()
		return true
	}

	return false
}

func (app *App) Redraw() {
	exit := ui.NewButton(timeprobeButtonExit, "Exit", 322, 191)
	exit.Width = 96

	kos.BeginRedraw()
	kos.OpenWindow(timeprobeWindowX, timeprobeWindowY, timeprobeWindowWidth, timeprobeWindowHeight, timeprobeWindowTitle)
	kos.DrawText(28, 44, ui.White, "System time: "+formatClock(app.clock))
	kos.DrawText(28, 64, ui.Silver, "Uptime: "+formatUint32(app.uptime)+" cs / "+formatCentisecondsAsSeconds(app.uptime))
	kos.DrawText(28, 84, ui.Aqua, "High precision uptime: "+formatHex64(app.uptimeNS))
	kos.DrawText(28, 104, ui.Lime, "Wait timeouts: "+formatUint32(app.timeoutTicks))
	kos.DrawText(28, 124, ui.Yellow, "Sleep 0.5s delta: "+formatUint32(app.sleepDelta)+" cs")
	kos.DrawText(28, 144, ui.White, "System clock source: syscall 3 / packed BCD")
	kos.DrawText(28, 164, ui.Silver, "Uptime sources: syscall 26.9 and 26.10")
	kos.DrawText(28, 184, ui.Aqua, "Last event: "+app.lastEvent)
	app.sleep.Draw()
	app.refresh.Draw()
	exit.Draw()
	kos.EndRedraw()
}

func (app *App) refreshTimeState() {
	app.clock = kos.SystemTime()
	app.uptime = kos.UptimeCentiseconds()
	app.uptimeNS = kos.UptimeNanoseconds()
}
