package sysinfo

import (
	"../../kos"
	"../../ui"
)

const (
	sysinfoButtonExit kos.ButtonID = 1
	sysinfoButtonToggleTitle kos.ButtonID = 2
	sysinfoButtonRefresh kos.ButtonID = 3
	sysinfoButtonFocusSelf kos.ButtonID = 4
	sysinfoButtonReapplyLayout kos.ButtonID = 5

	sysinfoWindowX = 350
	sysinfoWindowY = 180
	sysinfoWindowWidth = 540
	sysinfoWindowHeight = 344
	sysinfoWindowTitle = "KolibriOS Sysinfo"
	sysinfoUTF8Title = "KolibriOS Проба UTF-8"
)

type App struct {
	version kos.KernelVersionInfo
	screenWidth int
	screenHeight int
	workArea kos.Rect
	skinHeight int
	skinMargins kos.SkinMargins
	keyboardLanguage kos.KeyboardLanguage
	normalLayout kos.KeyboardLayoutTable
	shiftLayout kos.KeyboardLayoutTable
	altLayout kos.KeyboardLayoutTable
	hasKeyboardLayouts bool
	currentSlot int
	activeSlot int
	hasCurrentSlot bool
	usingUTF8Title bool
	focusStatus string
	layoutStatus string
	toggleTitle ui.Button
	refresh ui.Button
	focusSelf ui.Button
	reapplyLayout ui.Button
}

func NewApp() App {
	toggleTitle := ui.NewButton(sysinfoButtonToggleTitle, "Use UTF-8", 28, 264)
	toggleTitle.Width = 128

	refresh := ui.NewButton(sysinfoButtonRefresh, "Refresh", 176, 264)
	refresh.Width = 112

	focusSelf := ui.NewButton(sysinfoButtonFocusSelf, "Focus self", 320, 264)
	focusSelf.Width = 120

	reapplyLayout := ui.NewButton(sysinfoButtonReapplyLayout, "Reapply layout", 28, 296)
	reapplyLayout.Width = 144

	app := App{
		toggleTitle: toggleTitle,
		refresh: refresh,
		focusSelf: focusSelf,
		reapplyLayout: reapplyLayout,
		focusStatus: "ready",
		layoutStatus: "ready",
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
		app.focusStatus = "refreshed"
		app.layoutStatus = "reloaded"
		app.Redraw()
	case sysinfoButtonFocusSelf:
		app.focusSelfWindow()
		app.Redraw()
	case sysinfoButtonReapplyLayout:
		app.reapplyKeyboardLayout()
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
	kos.DrawText(28, 190, ui.Lime, "Skin margins: "+formatSkinMargins(app.skinMargins))
	kos.DrawText(28, 208, ui.Yellow, "Keyboard lang: "+formatKeyboardLanguage(app.keyboardLanguage))
	kos.DrawText(28, 226, ui.White, "Layout sums: "+app.layoutChecksumsString())
	kos.DrawText(320, 46, ui.Yellow, "Title mode: "+app.titleMode())
	kos.DrawText(320, 64, ui.White, "Current slot: "+app.currentSlotString())
	kos.DrawText(320, 82, ui.Silver, "Active slot: "+formatInt(app.activeSlot))
	kos.DrawText(320, 100, ui.Aqua, "Focus state: "+app.focusStatus)
	kos.DrawText(320, 118, ui.Lime, "Layout state: "+app.layoutStatus)
	kos.DrawText(320, 136, ui.Silver, "21.2 replays current tables / 26.2 reads them")
	kos.DrawText(320, 154, ui.Silver, "18.3 focuses a slot / 18.7 reports the active slot")
	app.toggleTitle.Draw()
	app.refresh.Draw()
	app.focusSelf.Draw()
	app.reapplyLayout.Draw()
	kos.EndRedraw()
}

func (app *App) refreshInfo() {
	var normalOK bool
	var shiftOK bool
	var altOK bool

	app.version = kos.KernelVersion()
	app.screenWidth, app.screenHeight = kos.ScreenSize()
	app.workArea = kos.ScreenWorkingArea()
	app.skinHeight = kos.SkinHeight()
	app.skinMargins = kos.WindowSkinMargins()
	app.keyboardLanguage = kos.KeyboardLayoutLanguage()
	app.normalLayout, normalOK = kos.ReadKeyboardLayoutTable(kos.KeyboardLayoutNormal)
	app.shiftLayout, shiftOK = kos.ReadKeyboardLayoutTable(kos.KeyboardLayoutShift)
	app.altLayout, altOK = kos.ReadKeyboardLayoutTable(kos.KeyboardLayoutAlt)
	app.hasKeyboardLayouts = normalOK && shiftOK && altOK
	app.currentSlot, app.hasCurrentSlot = kos.CurrentThreadSlotIndex()
	app.activeSlot = kos.ActiveWindowSlot()
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

func (app *App) currentSlotString() string {
	if !app.hasCurrentSlot {
		return "-"
	}

	return formatInt(app.currentSlot)
}

func (app *App) focusSelfWindow() {
	if !app.hasCurrentSlot {
		app.focusStatus = "current slot unavailable"
		return
	}

	kos.FocusWindowSlot(app.currentSlot)
	app.refreshInfo()
	if app.activeSlot == app.currentSlot {
		app.focusStatus = "self active"
		return
	}

	app.focusStatus = "focus requested for slot " + formatInt(app.currentSlot)
}

func (app *App) layoutChecksumsString() string {
	if !app.hasKeyboardLayouts {
		return "unavailable"
	}

	return formatLayoutChecksums(app.normalLayout, app.shiftLayout, app.altLayout)
}

func (app *App) reapplyKeyboardLayout() {
	if !app.hasKeyboardLayouts {
		app.layoutStatus = "layout tables unavailable"
		return
	}

	if !kos.SetKeyboardLayoutTable(kos.KeyboardLayoutNormal, &app.normalLayout) {
		app.layoutStatus = "normal layout apply failed"
		return
	}

	if !kos.SetKeyboardLayoutTable(kos.KeyboardLayoutShift, &app.shiftLayout) {
		app.layoutStatus = "shift layout apply failed"
		return
	}

	if !kos.SetKeyboardLayoutTable(kos.KeyboardLayoutAlt, &app.altLayout) {
		app.layoutStatus = "alt layout apply failed"
		return
	}

	if !kos.SetKeyboardLayoutLanguage(app.keyboardLanguage) {
		app.layoutStatus = "language apply failed"
		return
	}

	app.refreshInfo()
	app.layoutStatus = "layout round-trip ok"
}
