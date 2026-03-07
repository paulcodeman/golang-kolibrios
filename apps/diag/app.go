package diagapp

import (
	"errors"

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
	diagWindowHeight = 366

	diagReportPath     = "/FD/1/GODIAG.TXT"
	diagProbePath      = "/FD/1/GODIAG.TMP"
	diagHeadlessPath   = "/FD/1/GODIAG.AUTO"
	diagFilesProbePath = "/sys/default.skn"
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
	refresh := ui.NewButton(diagButtonRefresh, "Refresh", 28, 310)
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
	exit := ui.NewButton(diagButtonExit, "Exit", 170, 310)
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
		checkFiles(),
		checkSystem(version, screenWidth, screenHeight),
		checkReportProbe(),
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
