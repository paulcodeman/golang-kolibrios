package diagapp

import (
	"errors"
	"fmt"
	"io"
	"os"
	"syscall"

	"../../kos"
	"../../ui"
)

const (
	diagButtonExit    kos.ButtonID = 1
	diagButtonRefresh kos.ButtonID = 2

	diagWindowTitle  = "KolibriOS Go Diagnostics"
	diagWindowX      = 210
	diagWindowY      = 90
	diagWindowWidth  = 820
	diagWindowHeight = 420

	diagReportPath     = "/FD/1/GODIAG.TXT"
	diagProbePath      = "/FD/1/GODIAG.TMP"
	diagHeadlessPath   = "/FD/1/GODIAG.AUTO"
	diagFilesProbePath = "/sys/default.skn"
	diagOSProbeRoot    = "/FD/1"
	diagOSProbeDir     = "GOOSCHK"
	diagOSProbeFile    = "CHECK.TXT"
	diagOSRenamedFile  = "RENAMED.TXT"
	diagPreviewBytes   = 16
	diagLineHeight     = 20
	diagDebugPort      = 0x402

	diagHeadlessBegin = "[[GODIAG-BEGIN]]\r\n"
	diagHeadlessEnd   = "[[GODIAG-END]]\r\n"
)

var errFileProbe = &diagSentinel{text: "file probe failed"}

type diagSentinel struct {
	text string
}

func (err *diagSentinel) Error() string {
	return err.text
}

type wrappedError struct {
	cause error
}

func (err wrappedError) Error() string {
	return "wrapped"
}

func (err wrappedError) Unwrap() error {
	return err.cause
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

type checkResult struct {
	label  string
	ok     bool
	detail string
}

type snapshot struct {
	summary    string
	metaLine   string
	reportLine string
	reportBody string
	overallOK  bool
	results    []checkResult
}

type App struct {
	summary    string
	metaLine   string
	reportLine string
	reportBody string
	results    []checkResult
	overallOK  bool
	refreshBtn ui.Button
}

func NewApp() App {
	refresh := ui.NewButton(diagButtonRefresh, "Refresh", 28, 364)
	refresh.Width = 116

	app := App{
		refreshBtn: refresh,
	}
	app.refresh()
	return app
}

func (app *App) Run() {
	if app.tryHeadlessShutdown() {
		return
	}

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
	case diagButtonRefresh:
		app.refresh()
		app.Redraw()
	case diagButtonExit:
		kos.Exit()
		return true
	}

	return false
}

func (app *App) Redraw() {
	exit := ui.NewButton(diagButtonExit, "Exit", 170, 364)
	exit.Width = 96

	kos.BeginRedraw()
	kos.OpenWindow(diagWindowX, diagWindowY, diagWindowWidth, diagWindowHeight, diagWindowTitle)
	kos.DrawText(28, 44, app.summaryColor(), app.summary)
	kos.DrawText(28, 68, ui.Silver, app.metaLine)
	kos.DrawText(28, 92, app.summaryColor(), app.reportLine)

	y := 120
	for index := 0; index < len(app.results); index++ {
		kos.DrawText(28, y, app.resultColor(app.results[index].ok), formatCheckLine(app.results[index]))
		y += diagLineHeight
	}

	app.refreshBtn.Draw()
	exit.Draw()
	kos.EndRedraw()
}

func (app *App) refresh() {
	current := runDiagnostics()
	app.summary = current.summary
	app.metaLine = current.metaLine
	app.reportLine = current.reportLine
	app.reportBody = current.reportBody
	app.results = current.results
	app.overallOK = current.overallOK
}

func (app *App) summaryColor() kos.Color {
	if app.overallOK {
		return ui.Lime
	}

	return ui.Red
}

func (app *App) tryHeadlessShutdown() bool {
	if !headlessModeRequested() {
		return false
	}

	emitHeadlessReport(app.reportLine, app.reportBody)
	kos.SleepCentiseconds(2)
	return kos.PowerOff()
}

func (app *App) resultColor(ok bool) kos.Color {
	if ok {
		return ui.Lime
	}

	return ui.Red
}

func runDiagnostics() snapshot {
	headless := headlessModeRequested()
	version := kos.KernelVersion()
	screenWidth, screenHeight := kos.ScreenSize()
	currentTime := kos.SystemTime()
	uptime := kos.UptimeCentiseconds()
	currentFolder := diagnosticsCurrentFolder()

	results := []checkResult{
		checkClock(),
		checkStrings(),
		checkSlices(),
		checkInterfaces(),
		checkAssertions(),
		checkErrors(),
		checkSyscall(),
		checkFmt(),
		checkFiles(),
		checkOS(),
	}
	results = append(results, checkDLL())
	results = append(results, checkSystem(version, screenWidth, screenHeight))
	results = append(results, checkReportProbe())
	if headless {
		results = append(results, checkConsole())
	}

	overallOK := allResultsOK(results)
	summary := "diag: FAIL / one or more checks failed"
	if overallOK {
		summary = "diag: PASS / runtime and system probes stable"
	}

	metaLine := "Kernel " + formatKernelVersion(version) +
		" / screen " + formatInt(screenWidth) + "x" + formatInt(screenHeight) +
		" / cwd " + currentFolder +
		" / clock " + formatClock(currentTime) +
		" / uptime " + formatUint32(uptime)

	reportBody := buildReport(summary, metaLine, diagReportPath, results)
	reportLine := exportReport(diagReportPath, reportBody)

	return snapshot{
		summary:    summary,
		metaLine:   metaLine,
		reportLine: reportLine,
		reportBody: reportBody,
		overallOK:  overallOK,
		results:    results,
	}
}

func allResultsOK(results []checkResult) bool {
	for index := 0; index < len(results); index++ {
		if !results[index].ok {
			return false
		}
	}

	return true
}

func buildReport(summary string, metaLine string, reportPath string, results []checkResult) string {
	report := "KolibriOS Go Diagnostics\r\n"
	if allResultsOK(results) {
		report += "RESULT: PASS\r\n"
	} else {
		report += "RESULT: FAIL\r\n"
	}
	report += "REPORT: " + reportPath + "\r\n"
	report += "SUMMARY: " + summary + "\r\n"
	report += "META: " + metaLine + "\r\n"

	for index := 0; index < len(results); index++ {
		report += "CHECK " + results[index].label + ": " + formatStatusWord(results[index].ok) +
			" / " + results[index].detail + "\r\n"
	}

	return report
}

func diagnosticsCurrentFolder() string {
	folder := kos.CurrentFolder()
	if folder == "" {
		return "?"
	}

	return folder
}

func exportReport(reportPath string, report string) string {
	data := []byte(report)
	written, status := kos.CreateOrRewriteFile(reportPath, data)
	if status != kos.FileSystemOK {
		return "Export: FAIL / " + reportPath + " / " + formatFileSystemStatus(status)
	}
	if int(written) != len(data) {
		return "Export: FAIL / short write " + formatUint32(written) + " of " + formatInt(len(data))
	}

	return "Export: PASS / wrote " + formatUint32(written) + " bytes to " + reportPath
}

func headlessModeRequested() bool {
	_, status := kos.GetPathInfo(diagHeadlessPath)
	return status == kos.FileSystemOK
}

func emitHeadlessReport(reportLine string, reportBody string) {
	if !kos.ReservePorts(diagDebugPort, diagDebugPort) {
		return
	}

	kos.WritePortString(diagDebugPort, diagHeadlessBegin)
	if reportLine != "" {
		kos.WritePortString(diagDebugPort, reportLine)
		kos.WritePortString(diagDebugPort, "\r\n")
	}
	kos.WritePortString(diagDebugPort, reportBody)
	kos.WritePortString(diagDebugPort, diagHeadlessEnd)
	kos.ReleasePorts(diagDebugPort, diagDebugPort)
}

func formatCheckLine(result checkResult) string {
	return result.label + ": " + formatStatusWord(result.ok) + " / " + result.detail
}

func checkClock() checkResult {
	start := kos.UptimeCentiseconds()
	kos.SleepCentiseconds(1)
	end := kos.UptimeCentiseconds()
	clock := formatClock(kos.SystemTime())
	if end >= start {
		return checkResult{
			label:  "clock",
			ok:     true,
			detail: clock + " / uptime " + formatUint32(start) + " -> " + formatUint32(end),
		}
	}

	return checkResult{
		label:  "clock",
		ok:     false,
		detail: "uptime regressed around " + clock,
	}
}

func checkStrings() checkResult {
	message := "go" + " diagnostics"
	if message == "go diagnostics" {
		return checkResult{
			label:  "text",
			ok:     true,
			detail: "string concat and equality",
		}
	}

	return checkResult{
		label:  "text",
		ok:     false,
		detail: "string runtime mismatch",
	}
}

func checkSlices() checkResult {
	buffer := make([]byte, 0, 4)
	buffer = append(buffer, []byte("go")...)
	buffer = append(buffer, '!')
	out := make([]byte, len(buffer))
	copied := copy(out, buffer)
	if copied == 3 && string(out) == "go!" {
		return checkResult{
			label:  "slice",
			ok:     true,
			detail: "append copy convert",
		}
	}

	return checkResult{
		label:  "slice",
		ok:     false,
		detail: "append or copy mismatch",
	}
}

type sourceText interface {
	Text() string
}

type targetText interface {
	Text() string
}

type diagText struct {
	text string
}

func (value diagText) Text() string {
	return value.text
}

type fmtLabel struct {
	text string
}

func (label fmtLabel) String() string {
	return label.text
}

type fmtError struct {
	text string
}

func (err fmtError) Error() string {
	return err.text
}

type fmtBufferWriter struct {
	data []byte
}

func (writer *fmtBufferWriter) Write(data []byte) (int, error) {
	writer.data = append(writer.data, data...)
	return len(data), nil
}

func checkInterfaces() checkResult {
	var left sourceText = diagText{text: "iface diag"}
	var right sourceText = diagText{text: "iface diag"}
	if left == right && left.Text() == "iface diag" {
		return checkResult{
			label:  "iface",
			ok:     true,
			detail: "dispatch and equality",
		}
	}

	return checkResult{
		label:  "iface",
		ok:     false,
		detail: "interface runtime mismatch",
	}
}

func checkAssertions() checkResult {
	var any interface{} = diagText{text: "assert diag"}
	var iface sourceText = diagText{text: "bridge diag"}

	direct, okDirect := any.(diagText)
	bridge, okAny := any.(targetText)
	converted, okIface := iface.(targetText)
	switch describeAssertionValue(any) {
	case "switch diag":
		if okDirect &&
			okAny &&
			okIface &&
			direct.Text() == "assert diag" &&
			bridge.Text() == "assert diag" &&
			converted.Text() == "bridge diag" {
			return checkResult{
				label:  "assert",
				ok:     true,
				detail: "e2t e2i i2i switch",
			}
		}
	}

	return checkResult{
		label:  "assert",
		ok:     false,
		detail: "assertion runtime mismatch",
	}
}

func describeAssertionValue(value interface{}) string {
	switch value.(type) {
	case diagText:
		return "switch diag"
	default:
		return "switch default"
	}
}

func checkErrors() checkResult {
	sentinel := &diagSentinel{text: "diag sentinel"}
	wrapped := wrappedError{cause: sentinel}
	if !errors.Is(sentinel, sentinel) {
		return checkResult{
			label:  "errors",
			ok:     false,
			detail: "sentinel self-match failed",
		}
	}
	if errors.Unwrap(wrapped) != sentinel {
		return checkResult{
			label:  "errors",
			ok:     false,
			detail: "Unwrap lost wrapped cause",
		}
	}
	if !errors.Is(wrapped, sentinel) {
		return checkResult{
			label:  "errors",
			ok:     false,
			detail: "Is failed through Unwrap",
		}
	}

	return checkResult{
		label:  "errors",
		ok:     true,
		detail: "ordinary import path stable",
	}
}

func checkFmt() checkResult {
	label := fmtLabel{text: "fmt"}
	if fmt.Sprintf("ok") != "ok" {
		return checkResult{
			label:  "fmt",
			ok:     false,
			detail: "literal sprintf mismatch",
		}
	}
	if fmt.Sprintf("%s", "fmt") != "fmt" {
		return checkResult{
			label:  "fmt",
			ok:     false,
			detail: "%s mismatch",
		}
	}
	if fmt.Sprintf("%v", label) != "fmt" {
		return checkResult{
			label:  "fmt",
			ok:     false,
			detail: "%v stringer mismatch",
		}
	}
	if fmt.Sprintf("%d", 42) != "42" {
		return checkResult{
			label:  "fmt",
			ok:     false,
			detail: "%d mismatch",
		}
	}
	if fmt.Sprintf("%x", uint32(0x2A)) != "2a" {
		return checkResult{
			label:  "fmt",
			ok:     false,
			detail: "%x mismatch",
		}
	}
	if fmt.Sprintf("%t", true) != "true" {
		return checkResult{
			label:  "fmt",
			ok:     false,
			detail: "%t mismatch",
		}
	}
	if fmt.Sprintf("%%") != "%" {
		return checkResult{
			label:  "fmt",
			ok:     false,
			detail: "%% mismatch",
		}
	}
	if fmt.Sprintf("%v", fmtError{text: "diag fmt error"}) != "diag fmt error" {
		return checkResult{
			label:  "fmt",
			ok:     false,
			detail: "%v error mismatch",
		}
	}
	sprintfText := fmt.Sprintf("%v/%s/%d/%x/%t/%%", label, "ok", 42, uint32(0x2A), true)
	if sprintfText != "fmt/ok/42/2a/true/%" {
		return checkResult{
			label:  "fmt",
			ok:     false,
			detail: "Sprintf mismatch",
		}
	}

	sprintlnText := fmt.Sprintln(label, "line", true, 7)
	if sprintlnText != "fmt line true 7\n" {
		return checkResult{
			label:  "fmt",
			ok:     false,
			detail: "Sprintln mismatch",
		}
	}

	writer := &fmtBufferWriter{}
	written, err := fmt.Fprintf(writer, "cwd=%s / head=%x", kos.CurrentFolder(), []byte{0x4B, 0x50})
	if err != nil {
		return checkResult{
			label:  "fmt",
			ok:     false,
			detail: "Fprintf returned error",
		}
	}
	if string(writer.data) != "cwd="+kos.CurrentFolder()+" / head=4b50" || written != len(writer.data) {
		return checkResult{
			label:  "fmt",
			ok:     false,
			detail: "Fprintf mismatch",
		}
	}

	stdoutReader, stdoutWriter, pipeErr := os.Pipe()
	if pipeErr != nil {
		return checkResult{
			label:  "fmt",
			ok:     false,
			detail: "stdout pipe " + pipeErr.Error(),
		}
	}

	previousStdout := os.Stdout
	os.Stdout = stdoutWriter
	printCount, printErr := fmt.Print(label, " print ", 7, "\n")
	printfCount, printfErr := fmt.Printf("cwd=%s", kos.CurrentFolder())
	printlnCount, printlnErr := fmt.Println(" / tail", true)
	os.Stdout = previousStdout
	_ = stdoutWriter.Close()

	expectedPrint := "fmt print 7\ncwd=" + kos.CurrentFolder() + " / tail true\n"
	stdoutBuffer := make([]byte, len(expectedPrint))
	stdoutRead, stdoutReadErr := stdoutReader.Read(stdoutBuffer)
	_ = stdoutReader.Close()

	if printErr != nil || printfErr != nil || printlnErr != nil {
		return checkResult{
			label:  "fmt",
			ok:     false,
			detail: "Print returned error",
		}
	}
	if stdoutReadErr != nil {
		return checkResult{
			label:  "fmt",
			ok:     false,
			detail: "stdout read " + stdoutReadErr.Error(),
		}
	}
	if stdoutRead != len(expectedPrint) || string(stdoutBuffer[:stdoutRead]) != expectedPrint {
		return checkResult{
			label:  "fmt",
			ok:     false,
			detail: "Print mismatch",
		}
	}
	if printCount+printfCount+printlnCount != len(expectedPrint) {
		return checkResult{
			label:  "fmt",
			ok:     false,
			detail: "Print count mismatch",
		}
	}

	if fmt.Errorf("%v error %d", label, 7).Error() != "fmt error 7" {
		return checkResult{
			label:  "fmt",
			ok:     false,
			detail: "Errorf mismatch",
		}
	}
	return checkResult{
		label:  "fmt",
		ok:     true,
		detail: "sprintf fprintf print stdout",
	}
}

func checkSyscall() checkResult {
	var pipefd [2]int

	if err := syscall.Pipe(pipefd[:]); err != nil {
		return checkResult{
			label:  "syscall",
			ok:     false,
			detail: "pipe " + err.Error(),
		}
	}

	payload := []byte("diag syscall pipe")
	written, err := syscall.Write(pipefd[1], payload)
	if err != nil {
		return checkResult{
			label:  "syscall",
			ok:     false,
			detail: "write " + err.Error(),
		}
	}
	if written != len(payload) {
		return checkResult{
			label:  "syscall",
			ok:     false,
			detail: "short write",
		}
	}

	buffer := make([]byte, len(payload))
	read, err := syscall.Read(pipefd[0], buffer)
	if err != nil {
		return checkResult{
			label:  "syscall",
			ok:     false,
			detail: "read " + err.Error(),
		}
	}
	if read != len(payload) || string(buffer[:read]) != string(payload) {
		return checkResult{
			label:  "syscall",
			ok:     false,
			detail: "pipe payload mismatch",
		}
	}

	return checkResult{
		label:  "syscall",
		ok:     true,
		detail: "pipe read write " + formatInt(read) + " bytes",
	}
}

func checkFiles() checkResult {
	info, status := kos.GetPathInfo(diagFilesProbePath)
	if status != kos.FileSystemOK {
		err := probeError{
			op:     "get info",
			path:   diagFilesProbePath,
			status: status,
			cause:  errFileProbe,
		}
		return checkResult{
			label:  "files",
			ok:     false,
			detail: err.Error(),
		}
	}

	previewSize := diagPreviewBytes
	if info.Size > 0 && info.Size < uint64(previewSize) {
		previewSize = int(info.Size)
	}
	if previewSize == 0 {
		previewSize = diagPreviewBytes
	}

	buffer := make([]byte, previewSize)
	read, status := kos.ReadFile(diagFilesProbePath, buffer, 0)
	if status != kos.FileSystemOK && status != kos.FileSystemEOF {
		err := probeError{
			op:     "read",
			path:   diagFilesProbePath,
			status: status,
			cause:  errFileProbe,
		}
		return checkResult{
			label:  "files",
			ok:     false,
			detail: err.Error(),
		}
	}

	return checkResult{
		label: "files",
		ok:    true,
		detail: "size " + formatHex64(info.Size) +
			" / head " + formatBytePreview(buffer[:int(read)]),
	}
}

func checkOS() checkResult {
	base := diagOSProbeRoot
	if _, status := kos.GetPathInfo(base); status != kos.FileSystemOK {
		cwd, err := os.Getwd()
		if err != nil {
			return checkResult{
				label:  "os",
				ok:     false,
				detail: "getwd " + err.Error(),
			}
		}
		base = cwd
	}

	demoDir := diagJoinPath(base, diagOSProbeDir)
	demoFile := diagJoinPath(demoDir, diagOSProbeFile)
	renamedFile := diagJoinPath(demoDir, diagOSRenamedFile)
	payloadBase := "diag os"
	payloadExtra := " append"
	payload := payloadBase + payloadExtra

	if err := diagRemoveIfExists(renamedFile); err != nil {
		return checkResult{
			label:  "os",
			ok:     false,
			detail: "cleanup renamed " + err.Error(),
		}
	}
	if err := diagRemoveIfExists(demoFile); err != nil {
		return checkResult{
			label:  "os",
			ok:     false,
			detail: "cleanup file " + err.Error(),
		}
	}
	if err := diagRemoveIfExists(demoDir); err != nil {
		return checkResult{
			label:  "os",
			ok:     false,
			detail: "cleanup dir " + err.Error(),
		}
	}

	if err := os.Mkdir(demoDir, 0); err != nil {
		return checkResult{
			label:  "os",
			ok:     false,
			detail: "mkdir " + err.Error(),
		}
	}

	file, err := os.Create(demoFile)
	if err != nil {
		return checkResult{
			label:  "os",
			ok:     false,
			detail: "create " + err.Error(),
		}
	}

	if _, err = io.WriteString(file, payloadBase); err != nil {
		closeErr := file.Close()
		if closeErr != nil {
			err = closeErr
		}
		return checkResult{
			label:  "os",
			ok:     false,
			detail: "write " + err.Error(),
		}
	}
	if err = file.Close(); err != nil {
		return checkResult{
			label:  "os",
			ok:     false,
			detail: "close " + err.Error(),
		}
	}

	appendFile, err := os.OpenFile(demoFile, os.O_WRONLY|os.O_APPEND, 0)
	if err != nil {
		return checkResult{
			label:  "os",
			ok:     false,
			detail: "append open " + err.Error(),
		}
	}
	if _, err = io.WriteString(appendFile, payloadExtra); err != nil {
		closeErr := appendFile.Close()
		if closeErr != nil {
			err = closeErr
		}
		return checkResult{
			label:  "os",
			ok:     false,
			detail: "append write " + err.Error(),
		}
	}
	if err = appendFile.Close(); err != nil {
		return checkResult{
			label:  "os",
			ok:     false,
			detail: "append close " + err.Error(),
		}
	}

	reader, err := os.Open(demoFile)
	if err != nil {
		return checkResult{
			label:  "os",
			ok:     false,
			detail: "open " + err.Error(),
		}
	}

	data, readErr := io.ReadAll(reader)
	closeErr := reader.Close()
	if readErr == nil {
		readErr = closeErr
	}
	if readErr != nil {
		return checkResult{
			label:  "os",
			ok:     false,
			detail: "read " + readErr.Error(),
		}
	}

	if string(data) != payload {
		return checkResult{
			label:  "os",
			ok:     false,
			detail: "payload mismatch",
		}
	}

	if err := os.Rename(demoFile, renamedFile); err != nil {
		return checkResult{
			label:  "os",
			ok:     false,
			detail: "rename " + err.Error(),
		}
	}

	renamedData, err := os.ReadFile(renamedFile)
	if err != nil {
		return checkResult{
			label:  "os",
			ok:     false,
			detail: "renamed read " + err.Error(),
		}
	}
	if string(renamedData) != payload {
		return checkResult{
			label:  "os",
			ok:     false,
			detail: "renamed payload mismatch",
		}
	}

	if err := os.Remove(renamedFile); err != nil {
		return checkResult{
			label:  "os",
			ok:     false,
			detail: "remove file " + err.Error(),
		}
	}
	if err := os.Remove(demoDir); err != nil {
		return checkResult{
			label:  "os",
			ok:     false,
			detail: "remove dir " + err.Error(),
		}
	}

	return checkResult{
		label:  "os",
		ok:     true,
		detail: "cwd " + base + " / append rename cleanup",
	}
}

func checkDLL() checkResult {
	table := kos.LoadConsoleDLL()
	if table == 0 {
		return checkResult{
			label:  "dll",
			ok:     false,
			detail: "load " + kos.ConsoleDLLPath + " failed",
		}
	}

	initProc := table.Lookup("con_init")
	writeProc := table.Lookup("con_write_string")
	exitProc := table.Lookup("con_exit")
	startProc := table.Lookup("START")
	version := uint32(table.Lookup("version"))
	if !startProc.Valid() || !initProc.Valid() || !writeProc.Valid() || !exitProc.Valid() {
		return checkResult{
			label:  "dll",
			ok:     false,
			detail: "required console export missing",
		}
	}

	return checkResult{
		label:  "dll",
		ok:     true,
		detail: kos.ConsoleDLLPath +
			" / table " + formatHex64(uint64(table)) +
			" / ver " + formatHex64(uint64(version)) +
			" / start " + formatHex64(uint64(startProc)) +
			" / init " + formatHex64(uint64(initProc)) +
			" / write " + formatHex64(uint64(writeProc)),
	}
}

func checkConsole() checkResult {
	console, ok := kos.OpenConsole("KolibriOS Go Diagnostics Console")
	if !ok {
		return checkResult{
			label:  "console",
			ok:     false,
			detail: "open console failed",
		}
	}

	titleState := "title skipped"
	if console.SupportsTitle() {
		if !console.SetTitle("KolibriOS Go Diagnostics Console / live") {
			_ = console.Close()
			return checkResult{
				label:  "console",
				ok:     false,
				detail: "set title failed",
			}
		}
		titleState = "title ok"
	}

	if _, err := fmt.Fprintf(console, "golang-kolibrios console probe\r\n"); err != nil {
		_ = console.Close()
		return checkResult{
			label:  "console",
			ok:     false,
			detail: "fmt header failed",
		}
	}
	if _, err := fmt.Fprintf(console, "fmt writer path active / table 0x%x / ver 0x%x\r\n", uint32(console.ExportTable()), console.Version()); err != nil {
		_ = console.Close()
		return checkResult{
			label:  "console",
			ok:     false,
			detail: "fmt body failed",
		}
	}

	_ = console.Close()
	return checkResult{
		label:  "console",
		ok:     true,
		detail: "init fmt exit / " + titleState,
	}
}

func diagJoinPath(base string, name string) string {
	if base == "" || base == "/" {
		return "/" + name
	}
	if base[len(base)-1] == '/' {
		return base + name
	}

	return base + "/" + name
}

func diagRemoveIfExists(name string) error {
	err := os.Remove(name)
	if err == nil || errors.Is(err, os.ErrNotExist) {
		return nil
	}

	return err
}

func checkSystem(version kos.KernelVersionInfo, screenWidth int, screenHeight int) checkResult {
	workArea := kos.ScreenWorkingArea()
	skinHeight := kos.SkinHeight()
	margins := kos.WindowSkinMargins()
	if screenWidth <= 0 || screenHeight <= 0 {
		return checkResult{
			label:  "system",
			ok:     false,
			detail: "invalid screen size",
		}
	}
	if workArea.Width() <= 0 || workArea.Height() <= 0 {
		return checkResult{
			label:  "system",
			ok:     false,
			detail: "invalid work area",
		}
	}
	if workArea.Left < 0 || workArea.Top < 0 || workArea.Right >= screenWidth || workArea.Bottom >= screenHeight {
		return checkResult{
			label:  "system",
			ok:     false,
			detail: "work area outside screen",
		}
	}
	if skinHeight < 0 || margins.Left < 0 || margins.Top < 0 || margins.Right < 0 || margins.Bottom < 0 {
		return checkResult{
			label:  "system",
			ok:     false,
			detail: "invalid skin geometry",
		}
	}

	return checkResult{
		label: "system",
		ok:    true,
		detail: "kernel " + formatKernelVersion(version) +
			" / work " + formatInt(workArea.Width()) + "x" + formatInt(workArea.Height()) +
			" / skin " + formatInt(skinHeight),
	}
}

func checkReportProbe() checkResult {
	deleteStatus := kos.DeletePath(diagProbePath)
	if deleteStatus != kos.FileSystemOK && deleteStatus != kos.FileSystemNotFound {
		return checkResult{
			label:  "report",
			ok:     false,
			detail: "cleanup " + formatFileSystemStatus(deleteStatus),
		}
	}

	payload := []byte("go diag probe\r\n")
	written, status := kos.CreateOrRewriteFile(diagProbePath, payload)
	if status != kos.FileSystemOK {
		return checkResult{
			label:  "report",
			ok:     false,
			detail: "create " + formatFileSystemStatus(status),
		}
	}
	if int(written) != len(payload) {
		return checkResult{
			label:  "report",
			ok:     false,
			detail: "short write " + formatUint32(written),
		}
	}

	data, status := kos.ReadAllFile(diagProbePath)
	if status != kos.FileSystemOK && status != kos.FileSystemEOF {
		return checkResult{
			label:  "report",
			ok:     false,
			detail: "readback " + formatFileSystemStatus(status),
		}
	}
	if string(data) != string(payload) {
		return checkResult{
			label:  "report",
			ok:     false,
			detail: "readback mismatch",
		}
	}

	deleteStatus = kos.DeletePath(diagProbePath)
	if deleteStatus != kos.FileSystemOK {
		return checkResult{
			label:  "report",
			ok:     false,
			detail: "delete " + formatFileSystemStatus(deleteStatus),
		}
	}

	return checkResult{
		label:  "report",
		ok:     true,
		detail: "create read delete on " + diagProbePath,
	}
}

func Run() {
	app := NewApp()
	app.Run()
}
