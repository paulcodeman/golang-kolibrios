package pathdemo

import (
	"path"

	"../../kos"
	"../../ui"
)

const (
	pathButtonExit    kos.ButtonID = 1
	pathButtonRefresh kos.ButtonID = 2

	pathWindowTitle  = "KolibriOS Path Demo"
	pathWindowX      = 260
	pathWindowY      = 170
	pathWindowWidth  = 720
	pathWindowHeight = 276

	pathRawProbe = "/sys/./skins/../default.skn"
)

type App struct {
	summary    string
	cleanLine  string
	joinLine   string
	splitLine  string
	rootLine   string
	relLine    string
	infoLine   string
	ok         bool
	refreshBtn ui.Button
}

func NewApp() App {
	refresh := ui.NewButton(pathButtonRefresh, "Refresh", 28, 220)
	refresh.Width = 116

	app := App{
		refreshBtn: refresh,
	}
	app.refreshProbe()
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
	case pathButtonRefresh:
		app.refreshProbe()
		app.Redraw()
	case pathButtonExit:
		kos.Exit()
		return true
	}

	return false
}

func (app *App) Redraw() {
	exit := ui.NewButton(pathButtonExit, "Exit", 170, 220)
	exit.Width = 96

	kos.BeginRedraw()
	kos.OpenWindow(pathWindowX, pathWindowY, pathWindowWidth, pathWindowHeight, pathWindowTitle)
	kos.DrawText(28, 44, app.summaryColor(), app.summary)
	kos.DrawText(28, 66, ui.Silver, "This sample imports the ordinary path package: import \"path\"")
	kos.DrawText(28, 90, ui.Aqua, app.cleanLine)
	kos.DrawText(28, 112, ui.Lime, app.joinLine)
	kos.DrawText(28, 134, ui.Yellow, app.splitLine)
	kos.DrawText(28, 156, ui.White, app.rootLine)
	kos.DrawText(28, 178, ui.Silver, app.relLine)
	kos.DrawText(28, 200, ui.Black, app.infoLine)
	app.refreshBtn.Draw()
	exit.Draw()
	kos.EndRedraw()
}

func (app *App) refreshProbe() {
	const expectedPath = "/sys/default.skn"

	cleaned := path.Clean(pathRawProbe)
	joined := path.Join("/sys", ".", "skins", "..", "default.skn")
	dir, file := path.Split(cleaned)
	base := path.Base(cleaned)
	ext := path.Ext(cleaned)
	isAbs := path.IsAbs(cleaned)
	rootBase := path.Base("/")
	rootDir := path.Dir("/")
	relativeClean := path.Clean("skins/../../default.skn")

	app.cleanLine = "Clean: " + pathRawProbe + " -> " + cleaned
	app.joinLine = "Join: " + joined
	app.splitLine = "Split: dir " + dir + " / file " + file + " / base " + base + " / ext " + ext + " / abs " + formatBool(isAbs)
	app.rootLine = "Root: base " + rootBase + " / dir " + rootDir
	app.relLine = "Relative clean: skins/../../default.skn -> " + relativeClean

	if cleaned != expectedPath {
		app.fail("clean mismatch")
		return
	}
	if joined != expectedPath {
		app.fail("join mismatch")
		return
	}
	if dir != "/sys/" || file != "default.skn" || base != "default.skn" || ext != ".skn" || !isAbs {
		app.fail("split/base/ext mismatch")
		return
	}
	if rootBase != "/" || rootDir != "/" {
		app.fail("root semantics mismatch")
		return
	}
	if relativeClean != "../default.skn" {
		app.fail("relative clean mismatch")
		return
	}

	info, status := kos.GetPathInfo(cleaned)
	if status != kos.FileSystemOK {
		app.ok = false
		app.summary = "path probe failed / file info unavailable"
		app.infoLine = "Info: " + cleaned + " / " + formatFileSystemStatus(status)
		return
	}

	app.ok = true
	app.summary = "path probe ok / ordinary import path package resolved"
	app.infoLine = "Info: size " + formatHex64(info.Size) + " bytes / attrs " + formatHex32(uint32(info.Attributes))
}

func (app *App) fail(detail string) {
	app.ok = false
	app.summary = "path probe failed / " + detail
	app.infoLine = "Info: unavailable"
}

func (app *App) summaryColor() kos.Color {
	if app.ok {
		return ui.Lime
	}

	return ui.Red
}
