package sysinfo

import (
	"../../kos"
	"../../ui"
)

const (
	sysinfoButtonExit kos.ButtonID = 1
	sysinfoButtonToggleTitle kos.ButtonID = 2
	sysinfoButtonRefresh kos.ButtonID = 3

	sysinfoWindowX = 350
	sysinfoWindowY = 180
	sysinfoWindowWidth = 540
	sysinfoWindowHeight = 250
	sysinfoWindowTitle = "KolibriOS Sysinfo"
	sysinfoUTF8Title = "KolibriOS Проба UTF-8"
)

type App struct {
	version kos.KernelVersionInfo
	screenWidth int
	screenHeight int
	workArea kos.Rect
	skinHeight int
	usingUTF8Title bool
	toggleTitle ui.Button
	refresh ui.Button
}

func NewApp() App {
	toggleTitle := ui.NewButton(sysinfoButtonToggleTitle, "Use UTF-8", 28, 196)
	toggleTitle.Width = 128

	refresh := ui.NewButton(sysinfoButtonRefresh, "Refresh", 176, 196)
	refresh.Width = 112

	app := App{
		toggleTitle: toggleTitle,
		refresh: refresh,
	}
	app.refreshInfo()

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
	case sysinfoButtonToggleTitle:
		app.usingUTF8Title = !app.usingUTF8Title
		if app.usingUTF8Title {
			app.toggleTitle.Label = "Use ASCII"
			kos.SetWindowTitleWithEncodingPrefix(kos.EncodingUTF8, sysinfoUTF8Title)
		} else {
			app.toggleTitle.Label = "Use UTF-8"
			kos.SetWindowTitle(sysinfoWindowTitle)
		}
		app.Redraw()
	case sysinfoButtonRefresh:
		app.refreshInfo()
		app.Redraw()
	case sysinfoButtonExit:
		kos.Exit()
		return true
	}

	return false
}

func (app *App) Redraw() {
	kos.BeginRedraw()
	kos.OpenWindow(sysinfoWindowX, sysinfoWindowY, sysinfoWindowWidth, sysinfoWindowHeight, sysinfoWindowTitle)
	kos.DrawText(28, 46, ui.White, "Kernel version: "+formatKernelVersion(app.version))
	kos.DrawText(28, 64, ui.Silver, "Kernel ABI: "+formatKernelABI(app.version))
	kos.DrawText(28, 82, ui.Aqua, "Commit id: "+formatHex32(app.version.CommitID))
	kos.DrawText(28, 100, ui.Lime, "Debug tag: "+app.debugTagString())
	kos.DrawText(28, 118, ui.Yellow, "Screen size: "+formatInt(app.screenWidth)+"x"+formatInt(app.screenHeight))
	kos.DrawText(28, 136, ui.White, "Work area: "+formatRect(app.workArea))
	kos.DrawText(28, 154, ui.Silver, "Work size: "+formatInt(app.workArea.Width())+"x"+formatInt(app.workArea.Height()))
	kos.DrawText(28, 172, ui.Aqua, "Skin height: "+formatInt(app.skinHeight))
	kos.DrawText(320, 196, ui.Yellow, "Title mode: "+app.titleMode())
	kos.DrawText(320, 214, ui.Silver, "Refresh after skin or taskbar changes")
	app.toggleTitle.Draw()
	app.refresh.Draw()
	kos.EndRedraw()
}

func (app *App) refreshInfo() {
	app.version = kos.KernelVersion()
	app.screenWidth, app.screenHeight = kos.ScreenSize()
	app.workArea = kos.ScreenWorkingArea()
	app.skinHeight = kos.SkinHeight()
}

func (app *App) debugTagString() string {
	if !app.version.IsDebug() {
		return "release"
	}

	return formatHex8(app.version.DebugTag)
}

func (app *App) titleMode() string {
	if app.usingUTF8Title {
		return "71.1 UTF-8 prefix"
	}

	return "71.2 direct encoding"
}
