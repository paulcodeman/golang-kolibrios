package filesdemo

import (
	"errors"

	"../../kos"
	"../../ui"
)

const (
	filesButtonExit kos.ButtonID = 1
	filesButtonRefresh kos.ButtonID = 2

	filesProbePath = "/sys/default.skn"

	filesWindowX = 320
	filesWindowY = 180
	filesWindowWidth = 620
	filesWindowHeight = 252
	filesWindowTitle = "KolibriOS Files Demo"
	filesPreviewBytes = 16
)

var (
	errPathInfo = &pathSentinel{text: "path info failed"}
	errPathRead = &pathSentinel{text: "path read failed"}
)

type pathSentinel struct {
	text string
}

func (err *pathSentinel) Error() string {
	return err.text
}

type probeError struct {
	op     string
	path   string
	status kos.FileSystemStatus
	cause  error
}

func (err probeError) Error() string {
	return err.op + " " + err.path + " / " + formatFileSystemStatus(err.status)
}

func (err probeError) Unwrap() error {
	return err.cause
}

type App struct {
	path       string
	summary    string
	infoLine   string
	readLine   string
	classLine  string
	errorLine  string
	errorsOK   bool
	lastError  error
	refreshBtn ui.Button
}

func NewApp() App {
	refresh := ui.NewButton(filesButtonRefresh, "Refresh", 28, 198)
	refresh.Width = 116

	app := App{
		path:       filesProbePath,
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
	case filesButtonRefresh:
		app.refreshProbe()
		app.Redraw()
	case filesButtonExit:
		kos.Exit()
		return true
	}

	return false
}

func (app *App) Redraw() {
	exit := ui.NewButton(filesButtonExit, "Exit", 170, 198)
	exit.Width = 96

	kos.BeginRedraw()
	kos.OpenWindow(filesWindowX, filesWindowY, filesWindowWidth, filesWindowHeight, filesWindowTitle)
	kos.DrawText(28, 44, app.summaryColor(), app.summary)
	kos.DrawText(28, 66, ui.Silver, "This sample imports errors with the ordinary import path: import \"errors\"")
	kos.DrawText(28, 88, ui.Black, "Path: "+app.path)
	kos.DrawText(28, 112, ui.Aqua, app.infoLine)
	kos.DrawText(28, 134, ui.Lime, app.readLine)
	kos.DrawText(28, 156, ui.Yellow, app.classLine)
	kos.DrawText(28, 178, ui.Black, app.errorLine)
	app.refreshBtn.Draw()
	exit.Draw()
	kos.EndRedraw()
}

func (app *App) refreshProbe() {
	app.lastError = nil
	app.errorsOK = checkErrorsCompatibility()

	info, status := kos.GetPathInfo(app.path)
	if status != kos.FileSystemOK {
		app.fail(&probeError{
			op:     "get info",
			path:   app.path,
			status: status,
			cause:  errPathInfo,
		})
		return
	}

	app.infoLine = "Info: size " + formatHex64(info.Size) + " bytes / attrs " + formatHex32(uint32(info.Attributes))

	previewSize := filesPreviewBytes
	if info.Size > 0 && info.Size < uint64(previewSize) {
		previewSize = int(info.Size)
	}
	if previewSize == 0 {
		previewSize = filesPreviewBytes
	}

	buffer := make([]byte, previewSize)
	read, status := kos.ReadFile(app.path, buffer, 0)
	if status != kos.FileSystemOK && status != kos.FileSystemEOF {
		app.fail(&probeError{
			op:     "read",
			path:   app.path,
			status: status,
			cause:  errPathRead,
		})
		return
	}

	app.readLine = "Read: " + formatUint32(read) + " bytes / head " + formatBytePreview(buffer[:int(read)])
	if app.errorsOK {
		app.classLine = "errors: sentinel and unwrap chain ok"
		app.errorLine = "Error: none"
		app.summary = "file probe ok / bootstrap errors package resolved"
		return
	}

	app.errorLine = "Error: none"
	app.classLine = "errors: bootstrap self-check failed"
	app.summary = "file probe ok / bootstrap errors self-check failed"
}

func (app *App) fail(err error) {
	app.lastError = err
	app.readLine = "Read: unavailable"
	app.errorLine = "Error: " + err.Error()

	if errors.Is(err, errPathInfo) {
		app.infoLine = "Info: unavailable"
		app.summary = "file probe failed / errors.Is(errPathInfo)"
		app.classLine = "errors.Is(errPathInfo) = true"
		return
	}

	if errors.Is(err, errPathRead) {
		app.summary = "file probe failed / errors.Is(errPathRead)"
		app.classLine = "errors.Is(errPathRead) = true"
		return
	}

	app.summary = "file probe failed / unclassified error"
	app.classLine = "errors.Is: no sentinel match"
}

func (app *App) summaryColor() kos.Color {
	if app.lastError == nil && app.errorsOK {
		return ui.Lime
	}

	return ui.Red
}

func checkErrorsCompatibility() bool {
	wrapped := &probeError{
		op:     "diag",
		path:   filesProbePath,
		status: kos.FileSystemNotFound,
		cause:  errPathInfo,
	}
	return errors.Is(errPathInfo, errPathInfo) &&
		errors.Unwrap(wrapped) == errPathInfo &&
		errors.Is(wrapped, errPathInfo)
}
