package stringsdemo

import (
	"os"
	"strings"

	"../../kos"
	"../../ui"
)

const (
	stringsButtonExit    kos.ButtonID = 1
	stringsButtonRefresh kos.ButtonID = 2

	stringsWindowTitle  = "KolibriOS Strings Demo"
	stringsWindowX      = 250
	stringsWindowY      = 160
	stringsWindowWidth  = 748
	stringsWindowHeight = 298

	stringsProbePath = "/sys/default.skn"
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
	refresh := ui.NewButton(stringsButtonRefresh, "Refresh", 28, 242)
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
	case stringsButtonRefresh:
		app.refreshProbe()
		app.Redraw()
	case stringsButtonExit:
		os.Exit(0)
		return true
	}

	return false
}

func (app *App) Redraw() {
	exit := ui.NewButton(stringsButtonExit, "Exit", 170, 242)
	exit.Width = 96

	kos.BeginRedraw()
	kos.OpenWindow(stringsWindowX, stringsWindowY, stringsWindowWidth, stringsWindowHeight, stringsWindowTitle)
	kos.DrawText(28, 44, app.summaryColor(), app.summary)
	kos.DrawText(28, 66, ui.Silver, "This sample imports the ordinary strings package: import \"strings\"")
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
	joined := strings.Join([]string{"", "sys", "default.skn"}, "/")
	var builder strings.Builder
	builder.Grow(len(joined))
	_, _ = builder.WriteString("/sys/")
	_, _ = builder.WriteString("default")
	_ = builder.WriteByte('.')
	_, _ = builder.Write([]byte("skn"))
	built := builder.String()
	builderLen := builder.Len()
	builderCap := builder.Cap()
	builder.Reset()
	_, _ = builder.WriteString("builder ok")
	builderReset := builder.String()
	hasPrefix := strings.HasPrefix(joined, "/sys/")
	hasSuffix := strings.HasSuffix(joined, ".skn")
	contains := strings.Contains(joined, "default")
	index := strings.Index(joined, "default")
	lastSlash := strings.LastIndex(joined, "/")
	before, after, found := strings.Cut(joined, "/default")
	trimmed := strings.TrimSuffix(strings.TrimPrefix(joined, "/sys/"), ".skn")
	currentFolder, err := os.Getwd()
	if err != nil {
		app.fail("getwd failed")
		return
	}
	trimmedCWD := strings.TrimPrefix(currentFolder, "/")

	app.joinLine = "Join: " + joined + " / builder " + built + " / reset " + builderReset
	app.matchLine = "Match: prefix " + formatBool(hasPrefix) + " / suffix " + formatBool(hasSuffix) + " / contains " + formatBool(contains)
	app.indexLine = "Index: default " + formatInt(index) + " / last slash " + formatInt(lastSlash) + " / cut " + before + " | " + after + " / found " + formatBool(found)
	app.trimLine = "Trim: /sys/ + .skn -> " + trimmed
	app.cwdLine = "Current folder: " + currentFolder + " / trim leading slash -> " + trimmedCWD

	if joined != stringsProbePath {
		app.fail("join mismatch")
		return
	}
	if built != stringsProbePath || builderLen != len(stringsProbePath) || builderCap < builderLen || builderReset != "builder ok" {
		app.fail("builder mismatch")
		return
	}
	if !hasPrefix || !hasSuffix || !contains {
		app.fail("prefix suffix contains mismatch")
		return
	}
	if index != 5 || lastSlash != 4 || !found || before != "/sys" || after != ".skn" {
		app.fail("index or cut mismatch")
		return
	}
	if trimmed != "default" {
		app.fail("trim mismatch")
		return
	}
	if currentFolder == "" || trimmedCWD == currentFolder {
		app.fail("current folder trim mismatch")
		return
	}

	info, err := os.Stat(joined)
	if err != nil {
		app.ok = false
		app.summary = "strings probe failed / file info unavailable"
		app.infoLine = "Info: " + joined + " / " + err.Error()
		return
	}
	rawInfo, ok := info.Sys().(kos.FileInfo)
	if !ok {
		app.fail("stat sys payload mismatch")
		return
	}

	app.ok = true
	app.summary = "strings probe ok / ordinary import strings package resolved"
	app.infoLine = "Info: size " + formatHex64(uint64(info.Size())) + " bytes / attrs " + formatHex32(uint32(rawInfo.Attributes))
}

func (app *App) fail(detail string) {
	app.ok = false
	app.summary = "strings probe failed / " + detail
	app.infoLine = "Info: unavailable"
}

func (app *App) summaryColor() kos.Color {
	if app.ok {
		return ui.Lime
	}

	return ui.Red
}
