package osdemo

import (
	"errors"
	"io"
	"os"
	"path"

	"../../kos"
	"../../ui"
)

const (
	osButtonExit    kos.ButtonID = 1
	osButtonRefresh kos.ButtonID = 2

	osWindowTitle  = "KolibriOS OS Demo"
	osWindowX      = 230
	osWindowY      = 138
	osWindowWidth  = 820
	osWindowHeight = 344

	osDemoDirName      = "go-os-demo"
	osDemoFileName     = "sample.txt"
	osDemoRenamedName  = "renamed.txt"
	osDemoPayloadBase  = "KolibriOS os demo"
	osDemoPayloadExtra = " / append"
	osPreferredRoot    = "/FD/1"
)

type bufferWriter struct {
	data []byte
}

func (writer *bufferWriter) Write(buffer []byte) (int, error) {
	writer.data = append(writer.data, buffer...)
	return len(buffer), nil
}

type App struct {
	summary    string
	cwdLine    string
	writeLine  string
	readLine   string
	renameLine string
	infoLine   string
	ok         bool
	refreshBtn ui.Button
}

func NewApp() App {
	refresh := ui.NewButton(osButtonRefresh, "Refresh", 28, 286)
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
	case osButtonRefresh:
		app.refreshProbe()
		app.Redraw()
	case osButtonExit:
		kos.Exit()
		return true
	}

	return false
}

func (app *App) Redraw() {
	exit := ui.NewButton(osButtonExit, "Exit", 170, 286)
	exit.Width = 96

	kos.BeginRedraw()
	kos.OpenWindow(osWindowX, osWindowY, osWindowWidth, osWindowHeight, osWindowTitle)
	kos.DrawText(28, 44, app.summaryColor(), app.summary)
	kos.DrawText(28, 66, ui.Silver, "This sample imports the ordinary os package: import \"os\"")
	kos.DrawText(28, 92, ui.Aqua, app.cwdLine)
	kos.DrawText(28, 114, ui.Lime, app.writeLine)
	kos.DrawText(28, 136, ui.Yellow, app.readLine)
	kos.DrawText(28, 158, ui.Navy, app.renameLine)
	kos.DrawText(28, 180, ui.Black, app.infoLine)
	app.refreshBtn.Draw()
	exit.Draw()
	kos.EndRedraw()
}

func (app *App) refreshProbe() {
	cwd, err := os.Getwd()
	if err != nil {
		app.fail("getwd failed")
		return
	}

	baseDir := osPreferredRoot
	if _, status := kos.GetPathInfo(baseDir); status != kos.FileSystemOK {
		baseDir = cwd
	}

	demoDir := path.Join(baseDir, osDemoDirName)
	demoFile := path.Join(demoDir, osDemoFileName)
	renamedFile := path.Join(demoDir, osDemoRenamedName)
	payload := osDemoPayloadBase + osDemoPayloadExtra

	if err := removeIfExists(renamedFile); err != nil {
		app.fail("cleanup renamed file failed")
		return
	}
	if err := removeIfExists(demoFile); err != nil {
		app.fail("cleanup demo file failed")
		return
	}
	if err := removeIfExists(demoDir); err != nil {
		app.fail("cleanup demo dir failed")
		return
	}

	if err := os.Mkdir(demoDir, 0); err != nil {
		app.fail("mkdir failed")
		return
	}

	file, err := os.Create(demoFile)
	if err != nil {
		app.fail("create failed")
		return
	}

	wrote, err := io.WriteString(file, osDemoPayloadBase)
	if err == nil {
		var appendFile *os.File
		appendFile, err = os.OpenFile(demoFile, os.O_WRONLY|os.O_APPEND, 0)
		if err == nil {
			_, err = io.WriteString(appendFile, osDemoPayloadExtra)
			closeErr := appendFile.Close()
			if err == nil {
				err = closeErr
			}
		}
	}
	closeErr := file.Close()
	if err == nil {
		err = closeErr
	}
	if err != nil {
		app.fail("write failed")
		return
	}

	data, err := os.ReadFile(demoFile)
	if err != nil {
		app.fail("readfile failed")
		return
	}

	reader, err := os.Open(demoFile)
	if err != nil {
		app.fail("open failed")
		return
	}
	copyTarget := &bufferWriter{}
	copied, copyErr := io.Copy(copyTarget, reader)
	closeErr = reader.Close()
	if copyErr == nil {
		copyErr = closeErr
	}
	if copyErr != nil {
		app.fail("copy failed")
		return
	}

	info, status := kos.GetPathInfo(demoFile)
	if status != kos.FileSystemOK {
		app.ok = false
		app.summary = "os probe failed / file info unavailable"
		app.infoLine = "Info: " + demoFile + " / " + formatFileSystemStatus(status)
		return
	}

	if err := os.Rename(demoFile, renamedFile); err != nil {
		app.fail("rename failed")
		return
	}

	renamedData, err := os.ReadFile(renamedFile)
	if err != nil {
		app.fail("renamed read failed")
		return
	}

	if err := os.Remove(renamedFile); err != nil {
		app.fail("remove file failed")
		return
	}
	if err := os.Remove(demoDir); err != nil {
		app.fail("remove dir failed")
		return
	}

	app.cwdLine = "Getwd: " + cwd + " / probe root " + baseDir
	app.writeLine = "Mkdir/Create/OpenFile: wrote " + formatInt(wrote) + " + " + formatInt(len(osDemoPayloadExtra)) + " bytes into " + demoFile
	app.readLine = "ReadFile/Open+Copy: len " + formatInt(len(data)) + " / copy " + formatInt64(copied) + " / match " + formatBool(equalBytes(copyTarget.data, data))
	app.renameLine = "Rename/Remove: " + demoFile + " -> " + renamedFile + " / cleanup ok"

	if string(data) != payload || !equalBytes(copyTarget.data, []byte(payload)) {
		app.fail("payload mismatch")
		return
	}
	if copied != int64(len(payload)) {
		app.fail("copy length mismatch")
		return
	}
	if string(renamedData) != payload {
		app.fail("renamed payload mismatch")
		return
	}
	if info.Size != uint64(len(payload)) {
		app.fail("file size mismatch")
		return
	}

	app.ok = true
	app.summary = "os probe ok / ordinary import os package resolved"
	app.infoLine = "Info: size " + formatHex64(info.Size) + " bytes / attrs " + formatHex32(uint32(info.Attributes))
}

func (app *App) fail(detail string) {
	app.ok = false
	app.summary = "os probe failed / " + detail
	app.infoLine = "Info: unavailable"
}

func (app *App) summaryColor() kos.Color {
	if app.ok {
		return ui.Lime
	}

	return ui.Red
}

func removeIfExists(name string) error {
	err := os.Remove(name)
	if err == nil || errors.Is(err, os.ErrNotExist) {
		return nil
	}

	return err
}
