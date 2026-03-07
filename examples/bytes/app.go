package bytesdemo

import (
	"bytes"
	"os"

	"../../kos"
	"../../ui"
)

const (
	bytesButtonExit    kos.ButtonID = 1
	bytesButtonRefresh kos.ButtonID = 2

	bytesWindowTitle  = "KolibriOS Bytes Demo"
	bytesWindowX      = 248
	bytesWindowY      = 156
	bytesWindowWidth  = 768
	bytesWindowHeight = 298

	bytesProbePath = "/sys/default.skn"
)

type App struct {
	summary    string
	joinLine   string
	matchLine  string
	indexLine  string
	trimLine   string
	cwdLine    string
	infoLine   string
	ok         bool
	refreshBtn ui.Button
}

func NewApp() App {
	refresh := ui.NewButton(bytesButtonRefresh, "Refresh", 28, 242)
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
	case bytesButtonRefresh:
		app.refreshProbe()
		app.Redraw()
	case bytesButtonExit:
		os.Exit(0)
		return true
	}

	return false
}

func (app *App) Redraw() {
	exit := ui.NewButton(bytesButtonExit, "Exit", 170, 242)
	exit.Width = 96

	kos.BeginRedraw()
	kos.OpenWindow(bytesWindowX, bytesWindowY, bytesWindowWidth, bytesWindowHeight, bytesWindowTitle)
	kos.DrawText(28, 44, app.summaryColor(), app.summary)
	kos.DrawText(28, 66, ui.Silver, "This sample imports the ordinary bytes package: import \"bytes\"")
	kos.DrawText(28, 92, ui.Aqua, app.joinLine)
	kos.DrawText(28, 114, ui.Lime, app.matchLine)
	kos.DrawText(28, 136, ui.Yellow, app.indexLine)
	kos.DrawText(28, 158, ui.White, app.trimLine)
	kos.DrawText(28, 180, ui.Silver, app.cwdLine)
	kos.DrawText(28, 202, ui.Black, app.infoLine)
	app.refreshBtn.Draw()
	exit.Draw()
	kos.EndRedraw()
}

func (app *App) refreshProbe() {
	joined := bytes.Join([][]byte{[]byte{}, []byte("sys"), []byte("default.skn")}, []byte("/"))
	hasPrefix := bytes.HasPrefix(joined, []byte("/sys/"))
	hasSuffix := bytes.HasSuffix(joined, []byte(".skn"))
	contains := bytes.Contains(joined, []byte("default"))
	index := bytes.Index(joined, []byte("default"))
	dot := bytes.IndexByte(joined, '.')
	before, after, found := bytes.Cut(joined, []byte("/default"))
	trimmed := bytes.TrimSuffix(bytes.TrimPrefix(joined, []byte("/sys/")), []byte(".skn"))
	trimmedOK := bytes.Equal(trimmed, []byte("default"))
	currentFolderPath, err := os.Getwd()
	if err != nil {
		app.fail("getwd failed")
		return
	}
	currentFolder := []byte(currentFolderPath)
	trimmedCWD := bytes.TrimPrefix(currentFolder, []byte("/"))

	app.joinLine = "Join: " + string(joined)
	app.matchLine = "Match: prefix " + formatBool(hasPrefix) + " / suffix " + formatBool(hasSuffix) + " / contains " + formatBool(contains)
	app.indexLine = "Index: default " + formatInt(index) + " / dot " + formatInt(dot) + " / cut " + string(before) + " | " + string(after) + " / found " + formatBool(found)
	app.trimLine = "Trim: /sys/ + .skn -> " + string(trimmed) + " / equal " + formatBool(trimmedOK)
	app.cwdLine = "Current folder: " + string(currentFolder) + " / trim leading slash -> " + string(trimmedCWD)

	if !bytes.Equal(joined, []byte(bytesProbePath)) {
		app.fail("join mismatch")
		return
	}
	if !hasPrefix || !hasSuffix || !contains {
		app.fail("prefix suffix contains mismatch")
		return
	}
	if index != 5 || dot != 12 || !found || !bytes.Equal(before, []byte("/sys")) || !bytes.Equal(after, []byte(".skn")) {
		app.fail("index or cut mismatch")
		return
	}
	if !trimmedOK {
		app.fail("trim mismatch")
		return
	}
	if len(currentFolder) == 0 || bytes.Equal(trimmedCWD, currentFolder) {
		app.fail("current folder trim mismatch")
		return
	}

	info, err := os.Stat(string(joined))
	if err != nil {
		app.ok = false
		app.summary = "bytes probe failed / file info unavailable"
		app.infoLine = "Info: " + string(joined) + " / " + err.Error()
		return
	}
	rawInfo, ok := info.Sys().(kos.FileInfo)
	if !ok {
		app.fail("stat sys payload mismatch")
		return
	}

	app.ok = true
	app.summary = "bytes probe ok / ordinary import bytes package resolved"
	app.infoLine = "Info: size " + formatHex64(uint64(info.Size())) + " bytes / attrs " + formatHex32(uint32(rawInfo.Attributes))
}

func (app *App) fail(detail string) {
	app.ok = false
	app.summary = "bytes probe failed / " + detail
	app.infoLine = "Info: unavailable"
}

func (app *App) summaryColor() kos.Color {
	if app.ok {
		return ui.Lime
	}

	return ui.Red
}
