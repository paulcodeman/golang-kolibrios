package fmtdemo

import (
	"fmt"
	"os"

	"../../kos"
	"../../ui"
)

const (
	fmtButtonExit    kos.ButtonID = 1
	fmtButtonRefresh kos.ButtonID = 2

	fmtWindowTitle  = "KolibriOS Fmt Demo"
	fmtWindowX      = 228
	fmtWindowY      = 140
	fmtWindowWidth  = 852
	fmtWindowHeight = 356

	fmtProbePath    = "/sys/default.skn"
	fmtPreviewBytes = 12
)

type probeLabel struct {
	text string
}

func (label probeLabel) String() string {
	return label.text
}

type bufferWriter struct {
	data []byte
}

func (writer *bufferWriter) Write(data []byte) (int, error) {
	writer.data = append(writer.data, data...)
	return len(data), nil
}

type App struct {
	summary      string
	sprintfLine  string
	sprintlnLine string
	fprintfLine  string
	printLine    string
	errorLine    string
	scanLine     string
	infoLine     string
	ok           bool
	refreshBtn   ui.Button
}

func NewApp() App {
	refresh := ui.NewButton(fmtButtonRefresh, "Refresh", 28, 304)
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
	case fmtButtonRefresh:
		app.refreshProbe()
		app.Redraw()
	case fmtButtonExit:
		os.Exit(0)
		return true
	}

	return false
}

func (app *App) Redraw() {
	exit := ui.NewButton(fmtButtonExit, "Exit", 170, 304)
	exit.Width = 96

	kos.BeginRedraw()
	kos.OpenWindow(fmtWindowX, fmtWindowY, fmtWindowWidth, fmtWindowHeight, fmtWindowTitle)
	kos.DrawText(28, 44, app.summaryColor(), app.summary)
	kos.DrawText(28, 66, ui.Silver, "This sample imports the ordinary fmt package: import \"fmt\"")
	kos.DrawText(28, 92, ui.Aqua, app.sprintfLine)
	kos.DrawText(28, 114, ui.Lime, app.sprintlnLine)
	kos.DrawText(28, 136, ui.Navy, app.fprintfLine)
	kos.DrawText(28, 158, ui.Maroon, app.printLine)
	kos.DrawText(28, 180, ui.Yellow, app.errorLine)
	kos.DrawText(28, 202, ui.Aqua, app.scanLine)
	kos.DrawText(28, 224, ui.Black, app.infoLine)
	app.refreshBtn.Draw()
	exit.Draw()
	kos.EndRedraw()
}

func (app *App) refreshProbe() {
	info, status := kos.GetPathInfo(fmtProbePath)
	if status != kos.FileSystemOK {
		app.fail("file info unavailable", "Info: "+fmt.Sprintf("%s / status %d", fmtProbePath, uint32(status)))
		return
	}

	data, status := kos.ReadAllFile(fmtProbePath)
	if status != kos.FileSystemOK {
		app.fail("file read unavailable", "Info: "+fmt.Sprintf("%s / status %d", fmtProbePath, uint32(status)))
		return
	}
	if len(data) > fmtPreviewBytes {
		data = data[:fmtPreviewBytes]
	}

	currentFolder := kos.CurrentFolder()
	label := probeLabel{text: "fmt"}

	sprintfText := fmt.Sprintf("%v/%s/%d/%x/%t/%%", label, "ok", 42, uint32(0x2A), true)
	printlnText := fmt.Sprintln(label, "line", true, 7)
	writer := &bufferWriter{}
	written, writeErr := fmt.Fprintf(writer, "cwd=%s / head=%x / size=%d", currentFolder, data, len(data))
	stdoutReader, stdoutWriter, pipeErr := os.Pipe()
	if pipeErr != nil {
		app.fail("stdout pipe unavailable", "Info: "+pipeErr.Error())
		return
	}
	previousStdout := os.DefaultStdout()
	os.Stdout = stdoutWriter
	printWritten, printErr := fmt.Print(label, " print ", 7, "\n")
	printfWritten, printfErr := fmt.Printf("cwd=%s / size=%d", currentFolder, len(data))
	printlnWritten, printlnErr := fmt.Println(" / tail", true)
	os.Stdout = previousStdout
	_ = stdoutWriter.Close()
	formatErr := fmt.Errorf("%v error %d", label, 7)
	scanReader, scanWriter, scanPipeErr := os.Pipe()
	if scanPipeErr != nil {
		app.fail("scan pipe unavailable", "Info: "+scanPipeErr.Error())
		return
	}
	_, scanWriteErr := scanWriter.Write([]byte("scan 42 true\n"))
	_ = scanWriter.Close()
	if scanWriteErr != nil {
		_ = scanReader.Close()
		app.fail("Fscanln write failed", "Info: "+scanWriteErr.Error())
		return
	}

	var scanWord string
	var scanValue int
	var scanOK bool

	scanned, scanErr := fmt.Fscanln(scanReader, &scanWord, &scanValue, &scanOK)
	_ = scanReader.Close()
	if scanErr != nil {
		app.fail("Fscanln failed", "Info: "+scanErr.Error())
		return
	}

	defaultScanReader, defaultScanWriter, defaultScanPipeErr := os.Pipe()
	if defaultScanPipeErr != nil {
		app.fail("stdin pipe unavailable", "Info: "+defaultScanPipeErr.Error())
		return
	}
	_, defaultScanWriteErr := defaultScanWriter.Write([]byte("stdin 7 false\n"))
	_ = defaultScanWriter.Close()
	if defaultScanWriteErr != nil {
		_ = defaultScanReader.Close()
		app.fail("Scanln write failed", "Info: "+defaultScanWriteErr.Error())
		return
	}

	previousStdin := os.DefaultStdin()
	os.Stdin = defaultScanReader
	var stdinWord string
	var stdinValue int
	var stdinOK bool
	stdinScanned, stdinErr := fmt.Scanln(&stdinWord, &stdinValue, &stdinOK)
	os.Stdin = previousStdin
	_ = defaultScanReader.Close()
	if stdinErr != nil {
		app.fail("Scanln failed", "Info: "+stdinErr.Error())
		return
	}

	expectedPrint := "fmt print 7\ncwd=" + currentFolder + " / size=" + formatInt(len(data)) + " / tail true\n"
	stdoutData := make([]byte, len(expectedPrint))
	stdoutRead, stdoutReadErr := stdoutReader.Read(stdoutData)
	_ = stdoutReader.Close()

	app.sprintfLine = "Sprintf: " + sprintfText
	app.sprintlnLine = "Sprintln: " + trimTrailingNewline(printlnText) + " / newline " + fmt.Sprintf("%t", hasTrailingNewline(printlnText))
	app.fprintfLine = "Fprintf: wrote " + formatInt(written) + " / " + string(writer.data)
	app.printLine = "Print*: wrote " + formatInt(printWritten+printfWritten+printlnWritten) + " / stdout match " + fmt.Sprintf("%t", stdoutRead == len(expectedPrint) && string(stdoutData[:stdoutRead]) == expectedPrint)
	app.errorLine = "Errorf: " + formatErr.Error()
	app.scanLine = "Scan*: Fscanln " + scanWord + "/" + formatInt(scanValue) + "/" + fmt.Sprintf("%t", scanOK) + " / Scanln " + stdinWord + "/" + formatInt(stdinValue) + "/" + fmt.Sprintf("%t", stdinOK)
	app.infoLine = fmt.Sprintf("File: %s / size %d / attrs 0x%x / head %x", fmtProbePath, info.Size, uint32(info.Attributes), data)

	expectedSprintf := "fmt/ok/42/2a/true/%"
	if sprintfText != expectedSprintf {
		app.fail("Sprintf mismatch", "Info: expected "+expectedSprintf)
		return
	}

	expectedSprintln := "fmt line true 7\n"
	if printlnText != expectedSprintln {
		app.fail("Sprintln mismatch", "Info: expected "+trimTrailingNewline(expectedSprintln))
		return
	}

	expectedFprintf := "cwd=" + currentFolder + " / head=" + formatHexBytes(data) + " / size=" + formatInt(len(data))
	if writeErr != nil {
		app.fail("Fprintf returned error", "Info: "+writeErr.Error())
		return
	}
	if written != len(expectedFprintf) || string(writer.data) != expectedFprintf {
		app.fail("Fprintf mismatch", "Info: expected "+expectedFprintf)
		return
	}

	if printErr != nil || printfErr != nil || printlnErr != nil {
		app.fail("Print returned error", "Info: stdout write failed")
		return
	}
	if stdoutReadErr != nil {
		app.fail("stdout read failed", "Info: "+stdoutReadErr.Error())
		return
	}
	if stdoutRead != len(expectedPrint) || string(stdoutData[:stdoutRead]) != expectedPrint {
		app.fail("Print mismatch", "Info: expected "+trimTrailingNewline(expectedPrint))
		return
	}

	if formatErr.Error() != "fmt error 7" {
		app.fail("Errorf mismatch", "Info: expected fmt error 7")
		return
	}
	if scanned != 3 || scanWord != "scan" || scanValue != 42 || !scanOK {
		app.fail("Fscanln mismatch", "Info: expected scan/42/true")
		return
	}
	if stdinScanned != 3 || stdinWord != "stdin" || stdinValue != 7 || stdinOK {
		app.fail("Scanln mismatch", "Info: expected stdin/7/false")
		return
	}

	app.ok = true
	app.summary = "fmt probe ok / ordinary import fmt package resolved"
}

func (app *App) fail(detail string, info string) {
	app.ok = false
	app.summary = "fmt probe failed / " + detail
	app.infoLine = info
}

func (app *App) summaryColor() kos.Color {
	if app.ok {
		return ui.Lime
	}

	return ui.Red
}
