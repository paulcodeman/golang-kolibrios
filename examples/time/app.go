package timedemo

import (
	"os"
	"time"

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
	now          time.Time
	uptime       uint32
	uptimeNS     uint64
	timeoutTicks uint32
	sleepDelta   time.Duration
	unixStable   bool
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
		before := time.Now()
		time.Sleep(500 * time.Millisecond)
		app.sleepDelta = time.Since(before)
		app.refreshTimeState()
		app.lastEvent = "sleep delta / " + formatDurationMilliseconds(app.sleepDelta)
		app.Redraw()
	case timeprobeButtonRefresh:
		app.refreshTimeState()
		app.lastEvent = "manual refresh"
		app.Redraw()
	case timeprobeButtonExit:
		os.Exit(0)
		return true
	}

	return false
}

func (app *App) Redraw() {
	exit := ui.NewButton(timeprobeButtonExit, "Exit", 322, 191)
	exit.Width = 96

	kos.BeginRedraw()
	kos.OpenWindow(timeprobeWindowX, timeprobeWindowY, timeprobeWindowWidth, timeprobeWindowHeight, timeprobeWindowTitle)
	kos.DrawText(28, 44, ui.White, "time.Now(): "+formatTimeStamp(app.now))
	kos.DrawText(28, 64, ui.Silver, "Unix seconds: "+formatInt64(app.now.Unix()))
	kos.DrawText(28, 84, ui.Aqua, "Unix roundtrip: "+formatBoolWord(app.unixStable))
	kos.DrawText(28, 104, ui.Lime, "Sleep 0.5s delta: "+formatDurationMilliseconds(app.sleepDelta))
	kos.DrawText(28, 124, ui.Yellow, "Uptime: "+formatUint32(app.uptime)+" cs / "+formatCentisecondsAsSeconds(app.uptime))
	kos.DrawText(28, 144, ui.White, "High precision uptime: "+formatHex64(app.uptimeNS))
	kos.DrawText(28, 164, ui.Silver, "Wall clock source: syscalls 29 + 3 / YY => 2000+YY")
	kos.DrawText(28, 184, ui.Aqua, "Monotonic source: syscall 26.10 for Since/Sub")
	kos.DrawText(28, 204, ui.Lime, "Wait timeouts: "+formatUint32(app.timeoutTicks)+" / last "+app.lastEvent)
	app.sleep.Draw()
	app.refresh.Draw()
	exit.Draw()
	kos.EndRedraw()
}

func (app *App) refreshTimeState() {
	app.now = time.Now()
	app.uptime = kos.UptimeCentiseconds()
	app.uptimeNS = kos.UptimeNanoseconds()
	app.unixStable = time.Unix(app.now.Unix(), int64(app.now.Nanosecond())).Equal(app.now)
}
