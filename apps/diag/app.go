package diagapp

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"kos"
	"ui"
)

const (
	diagButtonExit    kos.ButtonID = 1
	diagButtonRefresh kos.ButtonID = 2

	diagWindowTitle  = "KolibriOS Go Diagnostics"
	diagWindowX      = 210
	diagWindowY      = 90
	diagWindowWidth  = 820
	diagWindowHeight = 580

	diagReportPath     = "/FD/1/GODIAG.TXT"
	diagProbePath      = "/FD/1/GODIAG.TMP"
	diagHeadlessPath   = "/FD/1/GODIAG.AUTO"
	diagFilesProbePath = "/sys/default.skn"
	diagFilepathRaw    = "\\sys\\.\\skins\\..\\default.skn"
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
	refresh := ui.NewButton(diagButtonRefresh, "Refresh", 28, 524)
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
		os.Exit(0)
		return true
	}

	return false
}

func (app *App) Redraw() {
	exit := ui.NewButton(diagButtonExit, "Exit", 170, 524)
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
		checkMaps(),
		checkInterfaces(),
		checkAssertions(),
		checkErrors(),
		checkTime(),
		checkSyscall(),
		checkFmt(),
		checkBufio(),
		checkBuilders(),
		checkStrconv(),
		checkFilepath(),
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
	_, err := os.Stat(diagHeadlessPath)
	return err == nil
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

func checkMaps() checkResult {
	small := make(map[string]int)
	small["alpha"] = 1
	small["beta"] = small["alpha"] + 2
	delete(small, "alpha")
	if _, ok := small["alpha"]; ok {
		return checkResult{
			label:  "maps",
			ok:     false,
			detail: "string delete left comma-ok hit",
		}
	}
	if small["beta"] != 3 {
		return checkResult{
			label:  "maps",
			ok:     false,
			detail: "string map lost assigned value",
		}
	}

	hinted := make(map[int]diagMapPair, 100)
	hinted[7] = diagMapPair{label: "seven", count: 7}
	hinted[9] = diagMapPair{label: "nine", count: 9}
	pair, ok := hinted[7]
	if !ok || pair.label != "seven" || pair.count != 7 {
		return checkResult{
			label:  "maps",
			ok:     false,
			detail: "int lookup lost struct value",
		}
	}
	delete(hinted, 9)
	if _, ok := hinted[9]; ok {
		return checkResult{
			label:  "maps",
			ok:     false,
			detail: "int delete left comma-ok hit",
		}
	}

	sum := 0
	seenSeven := false
	for key, value := range hinted {
		sum += key + value.count
		if key == 7 && value.label == "seven" {
			seenSeven = true
		}
	}
	if !seenSeven || sum != 14 {
		return checkResult{
			label:  "maps",
			ok:     false,
			detail: "range lost surviving entry",
		}
	}

	return checkResult{
		label:  "maps",
		ok:     true,
		detail: "string int delete comma-ok range",
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

type diagMapPair struct {
	label string
	count int
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

func checkTime() checkResult {
	const (
		sleepRequest = 50 * time.Millisecond
		minimumDelta = 10 * time.Millisecond
	)

	start := time.Now()
	time.Sleep(sleepRequest)
	delta := time.Since(start)
	now := time.Now()

	if now.Year() < 2000 || now.Year() > 2099 {
		return checkResult{
			label:  "time",
			ok:     false,
			detail: "year outside bootstrap range",
		}
	}
	if now.Month() < time.January || now.Month() > time.December {
		return checkResult{
			label:  "time",
			ok:     false,
			detail: "invalid month",
		}
	}
	if now.Day() < 1 || now.Day() > 31 || now.Hour() > 23 || now.Minute() > 59 || now.Second() > 59 {
		return checkResult{
			label:  "time",
			ok:     false,
			detail: "invalid wall clock fields",
		}
	}
	if !time.Unix(now.Unix(), int64(now.Nanosecond())).Equal(now) {
		return checkResult{
			label:  "time",
			ok:     false,
			detail: "unix roundtrip mismatch",
		}
	}
	if delta < minimumDelta {
		return checkResult{
			label:  "time",
			ok:     false,
			detail: "sleep delta too short / " + formatDurationMilliseconds(delta),
		}
	}

	return checkResult{
		label:  "time",
		ok:     true,
		detail: formatTimeStamp(now) + " / sleep " + formatDurationMilliseconds(delta),
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
	currentFolder, cwdErr := os.Getwd()
	if cwdErr != nil {
		return checkResult{
			label:  "fmt",
			ok:     false,
			detail: "getwd " + cwdErr.Error(),
		}
	}

	writer := &fmtBufferWriter{}
	written, err := fmt.Fprintf(writer, "cwd=%s / head=%x", currentFolder, []byte{0x4B, 0x50})
	if err != nil {
		return checkResult{
			label:  "fmt",
			ok:     false,
			detail: "Fprintf returned error",
		}
	}
	if string(writer.data) != "cwd="+currentFolder+" / head=4b50" || written != len(writer.data) {
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

	previousStdout := os.DefaultStdout()
	os.Stdout = stdoutWriter
	printCount, printErr := fmt.Print(label, " print ", 7, "\n")
	printfCount, printfErr := fmt.Printf("cwd=%s", currentFolder)
	printlnCount, printlnErr := fmt.Println(" / tail", true)
	os.Stdout = previousStdout
	_ = stdoutWriter.Close()

	expectedPrint := "fmt print 7\ncwd=" + currentFolder + " / tail true\n"
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

	scanReader, scanWriter, scanPipeErr := os.Pipe()
	if scanPipeErr != nil {
		return checkResult{
			label:  "fmt",
			ok:     false,
			detail: "stdin pipe " + scanPipeErr.Error(),
		}
	}
	_, scanWriteErr := scanWriter.Write([]byte("scan 42 true\n"))
	_ = scanWriter.Close()
	if scanWriteErr != nil {
		_ = scanReader.Close()
		return checkResult{
			label:  "fmt",
			ok:     false,
			detail: "Fscanln write " + scanWriteErr.Error(),
		}
	}

	var scanWord string
	var scanValue int
	var scanOK bool

	scanned, scanErr := fmt.Fscanln(scanReader, &scanWord, &scanValue, &scanOK)
	_ = scanReader.Close()
	if scanErr != nil {
		return checkResult{
			label:  "fmt",
			ok:     false,
			detail: "Fscanln " + scanErr.Error(),
		}
	}
	if scanned != 3 || scanWord != "scan" || scanValue != 42 || !scanOK {
		return checkResult{
			label:  "fmt",
			ok:     false,
			detail: "Fscanln mismatch",
		}
	}

	defaultScanReader, defaultScanWriter, defaultScanPipeErr := os.Pipe()
	if defaultScanPipeErr != nil {
		return checkResult{
			label:  "fmt",
			ok:     false,
			detail: "default stdin pipe " + defaultScanPipeErr.Error(),
		}
	}
	_, defaultScanWriteErr := defaultScanWriter.Write([]byte("stdin 7 false\n"))
	_ = defaultScanWriter.Close()
	if defaultScanWriteErr != nil {
		_ = defaultScanReader.Close()
		return checkResult{
			label:  "fmt",
			ok:     false,
			detail: "Scanln write " + defaultScanWriteErr.Error(),
		}
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
		return checkResult{
			label:  "fmt",
			ok:     false,
			detail: "Scanln " + stdinErr.Error(),
		}
	}
	if stdinScanned != 3 || stdinWord != "stdin" || stdinValue != 7 || stdinOK {
		return checkResult{
			label:  "fmt",
			ok:     false,
			detail: "Scanln mismatch",
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
		detail: "sprintf fprintf print scan stdout stdin",
	}
}

func checkFilepath() checkResult {
	const expectedPath = "/sys/default.skn"

	cleaned := filepath.Clean(diagFilepathRaw)
	joined := filepath.Join("/sys", ".", "skins", "..", "default.skn")
	dir, file := filepath.Split(cleaned)
	base := filepath.Base(cleaned)
	ext := filepath.Ext(cleaned)
	slashed := filepath.ToSlash(diagFilepathRaw)
	restored := filepath.FromSlash(expectedPath)
	volume := filepath.VolumeName(cleaned)
	relativeAbs, err := filepath.Abs("default.skn")
	if err != nil {
		return checkResult{
			label:  "filepath",
			ok:     false,
			detail: "abs failed: " + err.Error(),
		}
	}

	if cleaned != expectedPath {
		return checkResult{
			label:  "filepath",
			ok:     false,
			detail: "clean mismatch: " + cleaned,
		}
	}
	if joined != expectedPath {
		return checkResult{
			label:  "filepath",
			ok:     false,
			detail: "join mismatch: " + joined,
		}
	}
	if dir != "/sys/" || file != "default.skn" || base != "default.skn" || ext != ".skn" || !filepath.IsAbs(cleaned) {
		return checkResult{
			label:  "filepath",
			ok:     false,
			detail: "split/base/ext mismatch",
		}
	}
	if filepath.Clean(slashed) != expectedPath || restored != expectedPath || volume != "" {
		return checkResult{
			label:  "filepath",
			ok:     false,
			detail: "slash or volume mismatch",
		}
	}
	if !filepath.IsAbs(relativeAbs) || filepath.Base(relativeAbs) != "default.skn" {
		return checkResult{
			label:  "filepath",
			ok:     false,
			detail: "abs semantics mismatch",
		}
	}

	info, err := os.Stat(cleaned)
	if err != nil {
		return checkResult{
			label:  "filepath",
			ok:     false,
			detail: "stat failed: " + err.Error(),
		}
	}

	return checkResult{
		label:  "filepath",
		ok:     true,
		detail: "clean join split abs volume / size " + formatHex64(uint64(info.Size())),
	}
}

func checkBufio() checkResult {
	readerPipe, writerPipe, err := os.Pipe()
	if err != nil {
		return checkResult{
			label:  "bufio",
			ok:     false,
			detail: "pipe unavailable: " + err.Error(),
		}
	}

	bufferedWriter := bufio.NewWriter(writerPipe)
	if _, err = bufferedWriter.WriteString("alpha beta\n"); err != nil {
		_ = readerPipe.Close()
		_ = writerPipe.Close()
		return checkResult{
			label:  "bufio",
			ok:     false,
			detail: "WriteString failed: " + err.Error(),
		}
	}
	if err = bufferedWriter.WriteByte('g'); err != nil {
		_ = readerPipe.Close()
		_ = writerPipe.Close()
		return checkResult{
			label:  "bufio",
			ok:     false,
			detail: "WriteByte failed: " + err.Error(),
		}
	}
	if _, err = bufferedWriter.WriteString("amma\n"); err != nil {
		_ = readerPipe.Close()
		_ = writerPipe.Close()
		return checkResult{
			label:  "bufio",
			ok:     false,
			detail: "WriteString tail failed: " + err.Error(),
		}
	}
	if err = bufferedWriter.Flush(); err != nil {
		_ = readerPipe.Close()
		_ = writerPipe.Close()
		return checkResult{
			label:  "bufio",
			ok:     false,
			detail: "Flush failed: " + err.Error(),
		}
	}
	_ = writerPipe.Close()

	bufferedReader := bufio.NewReader(readerPipe)
	firstByte, err := bufferedReader.ReadByte()
	if err != nil {
		_ = readerPipe.Close()
		return checkResult{
			label:  "bufio",
			ok:     false,
			detail: "ReadByte failed: " + err.Error(),
		}
	}
	if err = bufferedReader.UnreadByte(); err != nil {
		_ = readerPipe.Close()
		return checkResult{
			label:  "bufio",
			ok:     false,
			detail: "UnreadByte failed: " + err.Error(),
		}
	}
	firstLine, err := bufferedReader.ReadString('\n')
	if err != nil {
		_ = readerPipe.Close()
		return checkResult{
			label:  "bufio",
			ok:     false,
			detail: "ReadString failed: " + err.Error(),
		}
	}
	secondLine, err := bufferedReader.ReadBytes('\n')
	if err != nil {
		_ = readerPipe.Close()
		return checkResult{
			label:  "bufio",
			ok:     false,
			detail: "ReadBytes failed: " + err.Error(),
		}
	}
	_, eofErr := bufferedReader.ReadByte()
	_ = readerPipe.Close()

	brokenReader, brokenWriter, err := os.Pipe()
	if err != nil {
		return checkResult{
			label:  "bufio",
			ok:     false,
			detail: "broken pipe unavailable: " + err.Error(),
		}
	}
	_ = brokenReader.Close()
	_, brokenErr := brokenWriter.Write([]byte("x"))
	_ = brokenWriter.Close()

	linesReader, linesWriter, err := os.Pipe()
	if err != nil {
		return checkResult{
			label:  "bufio",
			ok:     false,
			detail: "line pipe unavailable: " + err.Error(),
		}
	}
	if _, err = linesWriter.Write([]byte("line one\nline two\n")); err != nil {
		_ = linesReader.Close()
		_ = linesWriter.Close()
		return checkResult{
			label:  "bufio",
			ok:     false,
			detail: "line pipe write failed: " + err.Error(),
		}
	}
	_ = linesWriter.Close()

	lineScanner := bufio.NewScanner(linesReader)
	lineA := ""
	lineB := ""
	if lineScanner.Scan() {
		lineA = lineScanner.Text()
	}
	if lineScanner.Scan() {
		lineB = lineScanner.Text()
	}
	lineScanErr := lineScanner.Err()
	_ = linesReader.Close()

	wordsReader, wordsWriter, err := os.Pipe()
	if err != nil {
		return checkResult{
			label:  "bufio",
			ok:     false,
			detail: "word pipe unavailable: " + err.Error(),
		}
	}
	if _, err = wordsWriter.Write([]byte("one two three\n")); err != nil {
		_ = wordsReader.Close()
		_ = wordsWriter.Close()
		return checkResult{
			label:  "bufio",
			ok:     false,
			detail: "word pipe write failed: " + err.Error(),
		}
	}
	_ = wordsWriter.Close()

	wordScanner := bufio.NewScanner(wordsReader)
	wordScanner.Split(bufio.ScanWords)
	wordA := ""
	wordB := ""
	wordC := ""
	if wordScanner.Scan() {
		wordA = wordScanner.Text()
	}
	if wordScanner.Scan() {
		wordB = wordScanner.Text()
	}
	if wordScanner.Scan() {
		wordC = wordScanner.Text()
	}
	wordScanErr := wordScanner.Err()
	_ = wordsReader.Close()

	bytesReader, bytesWriter, err := os.Pipe()
	if err != nil {
		return checkResult{
			label:  "bufio",
			ok:     false,
			detail: "byte pipe unavailable: " + err.Error(),
		}
	}
	if _, err = bytesWriter.Write([]byte("AZ")); err != nil {
		_ = bytesReader.Close()
		_ = bytesWriter.Close()
		return checkResult{
			label:  "bufio",
			ok:     false,
			detail: "byte pipe write failed: " + err.Error(),
		}
	}
	_ = bytesWriter.Close()

	byteScanner := bufio.NewScanner(bytesReader)
	byteScanner.Split(bufio.ScanBytes)
	byteA := ""
	byteB := ""
	if byteScanner.Scan() {
		byteA = byteScanner.Text()
	}
	if byteScanner.Scan() {
		byteB = byteScanner.Text()
	}
	byteScanErr := byteScanner.Err()
	_ = bytesReader.Close()

	if firstByte != 'a' || firstLine != "alpha beta\n" || string(secondLine) != "gamma\n" {
		return checkResult{
			label:  "bufio",
			ok:     false,
			detail: "reader path mismatch",
		}
	}
	if !errors.Is(eofErr, io.EOF) {
		return checkResult{
			label:  "bufio",
			ok:     false,
			detail: "EOF mismatch",
		}
	}
	if !errors.Is(brokenErr, syscall.EPIPE) {
		return checkResult{
			label:  "bufio",
			ok:     false,
			detail: "EPIPE mismatch",
		}
	}
	if lineScanErr != nil || lineA != "line one" || lineB != "line two" {
		return checkResult{
			label:  "bufio",
			ok:     false,
			detail: "ScanLines mismatch",
		}
	}
	if wordScanErr != nil || wordA != "one" || wordB != "two" || wordC != "three" {
		return checkResult{
			label:  "bufio",
			ok:     false,
			detail: "ScanWords mismatch",
		}
	}
	if byteScanErr != nil || byteA != "A" || byteB != "Z" {
		return checkResult{
			label:  "bufio",
			ok:     false,
			detail: "ScanBytes mismatch",
		}
	}

	return checkResult{
		label:  "bufio",
		ok:     true,
		detail: "reader eof epipe scanner / line one line two / one two three / A Z",
	}
}

func checkBuilders() checkResult {
	var builder strings.Builder
	builder.Grow(len(diagFilesProbePath))
	_, builderWriteErr := builder.WriteString("/sys/")
	if builderWriteErr != nil {
		return checkResult{
			label:  "builders",
			ok:     false,
			detail: "builder WriteString failed: " + builderWriteErr.Error(),
		}
	}
	_, builderBytesErr := builder.Write([]byte("default"))
	if builderBytesErr != nil {
		return checkResult{
			label:  "builders",
			ok:     false,
			detail: "builder Write failed: " + builderBytesErr.Error(),
		}
	}
	if err := builder.WriteByte('.'); err != nil {
		return checkResult{
			label:  "builders",
			ok:     false,
			detail: "builder WriteByte failed: " + err.Error(),
		}
	}
	_, builderTailErr := builder.WriteString("skn")
	if builderTailErr != nil {
		return checkResult{
			label:  "builders",
			ok:     false,
			detail: "builder tail failed: " + builderTailErr.Error(),
		}
	}
	built := builder.String()
	builderLen := builder.Len()
	builderCap := builder.Cap()
	builder.Reset()
	_, _ = builder.WriteString("builder ok")
	builderReset := builder.String()

	buffer := bytes.NewBuffer(nil)
	buffer.Grow(len(diagFilesProbePath))
	_, bufferWriteErr := buffer.WriteString("/sys/")
	if bufferWriteErr != nil {
		return checkResult{
			label:  "builders",
			ok:     false,
			detail: "buffer WriteString failed: " + bufferWriteErr.Error(),
		}
	}
	_, bufferBytesErr := buffer.Write([]byte("default"))
	if bufferBytesErr != nil {
		return checkResult{
			label:  "builders",
			ok:     false,
			detail: "buffer Write failed: " + bufferBytesErr.Error(),
		}
	}
	if err := buffer.WriteByte('.'); err != nil {
		return checkResult{
			label:  "builders",
			ok:     false,
			detail: "buffer WriteByte failed: " + err.Error(),
		}
	}
	_, bufferTailErr := buffer.WriteString("skn")
	if bufferTailErr != nil {
		return checkResult{
			label:  "builders",
			ok:     false,
			detail: "buffer tail failed: " + bufferTailErr.Error(),
		}
	}
	bufferSnapshot := append([]byte(nil), buffer.Bytes()...)
	bufferString := buffer.String()
	bufferLen := buffer.Len()
	bufferCap := buffer.Cap()
	buffer.Reset()
	_, _ = buffer.WriteString("buffer ok")
	bufferReset := buffer.String()
	bufferFromString := bytes.NewBufferString("demo")
	splitParts := strings.Split(diagFilesProbePath, "/")
	splitTwo := strings.SplitN(diagFilesProbePath, "/", 2)
	fields := strings.Fields("alpha  beta\tgamma")
	trimmed := strings.TrimSpace(" \tdefault \n")
	replaced := strings.ReplaceAll(diagFilesProbePath, ".skn", ".txt")
	stringReader := strings.NewReader(diagFilesProbePath)
	stringHead := make([]byte, 4)
	stringHeadRead, stringHeadErr := stringReader.Read(stringHead)
	stringHeadByte, stringHeadByteErr := stringReader.ReadByte()
	stringUnreadErr := stringReader.UnreadByte()
	stringHeadByteAgain, stringHeadByteAgainErr := stringReader.ReadByte()
	stringSeekPos, stringSeekErr := stringReader.Seek(-4, io.SeekEnd)
	stringTail := make([]byte, 4)
	stringTailRead, stringTailErr := stringReader.Read(stringTail)
	stringReadAt := make([]byte, 7)
	stringReadAtCount, stringReadAtErr := stringReader.ReadAt(stringReadAt, 5)
	stringReaderLen := stringReader.Len()
	stringReaderSize := stringReader.Size()
	stringCopyReader := strings.NewReader(diagFilesProbePath)
	var stringCopyBuilder strings.Builder
	stringCopied, stringCopyErr := io.Copy(&stringCopyBuilder, stringCopyReader)
	byteParts := bytes.Split([]byte(diagFilesProbePath), []byte("/"))
	byteSplitTwo := bytes.SplitN([]byte(diagFilesProbePath), []byte("/"), 2)
	byteFields := bytes.Fields([]byte("alpha  beta\tgamma"))
	byteTrimmed := bytes.TrimSpace([]byte(" \tdefault \n"))
	byteReplaced := bytes.ReplaceAll([]byte(diagFilesProbePath), []byte(".skn"), []byte(".txt"))
	byteReader := bytes.NewReader([]byte(diagFilesProbePath))
	byteHead := make([]byte, 4)
	byteHeadRead, byteHeadErr := byteReader.Read(byteHead)
	byteHeadValue, byteHeadValueErr := byteReader.ReadByte()
	byteUnreadErr := byteReader.UnreadByte()
	byteHeadValueAgain, byteHeadValueAgainErr := byteReader.ReadByte()
	byteSeekPos, byteSeekErr := byteReader.Seek(-4, io.SeekEnd)
	byteTail := make([]byte, 4)
	byteTailRead, byteTailErr := byteReader.Read(byteTail)
	byteReadAt := make([]byte, 7)
	byteReadAtCount, byteReadAtErr := byteReader.ReadAt(byteReadAt, 5)
	byteReaderLen := byteReader.Len()
	byteReaderSize := byteReader.Size()
	byteCopyReader := bytes.NewReader([]byte(diagFilesProbePath))
	byteCopyBuffer := bytes.NewBuffer(nil)
	byteCopied, byteCopyErr := io.Copy(byteCopyBuffer, byteCopyReader)

	if built != diagFilesProbePath || builderLen != len(diagFilesProbePath) || builderCap < builderLen || builderReset != "builder ok" {
		return checkResult{
			label:  "builders",
			ok:     false,
			detail: "strings.Builder mismatch",
		}
	}
	if !bytes.Equal(bufferSnapshot, []byte(diagFilesProbePath)) || bufferString != diagFilesProbePath || bufferLen != len(diagFilesProbePath) || bufferCap < bufferLen || bufferReset != "buffer ok" || bufferFromString.String() != "demo" {
		return checkResult{
			label:  "builders",
			ok:     false,
			detail: "bytes.Buffer mismatch",
		}
	}
	if len(splitParts) != 3 || splitParts[1] != "sys" || splitParts[2] != "default.skn" || len(splitTwo) != 2 || splitTwo[1] != "sys/default.skn" || len(fields) != 3 || fields[0] != "alpha" || fields[2] != "gamma" || trimmed != "default" || replaced != "/sys/default.txt" {
		return checkResult{
			label:  "builders",
			ok:     false,
			detail: "strings helper mismatch",
		}
	}
	if stringHeadErr != nil || stringHeadByteErr != nil || stringUnreadErr != nil || stringHeadByteAgainErr != nil || stringSeekErr != nil || stringTailErr != nil || stringReadAtErr != nil || stringCopyErr != nil || stringHeadRead != 4 || string(stringHead[:stringHeadRead]) != "/sys" || stringHeadByte != '/' || stringHeadByteAgain != '/' || stringSeekPos != 12 || stringTailRead != 4 || string(stringTail[:stringTailRead]) != ".skn" || stringReadAtCount != 7 || string(stringReadAt[:stringReadAtCount]) != "default" || stringReaderLen != 0 || stringReaderSize != int64(len(diagFilesProbePath)) || stringCopied != int64(len(diagFilesProbePath)) || stringCopyBuilder.String() != diagFilesProbePath {
		return checkResult{
			label:  "builders",
			ok:     false,
			detail: "strings reader mismatch",
		}
	}
	if len(byteParts) != 3 || !bytes.Equal(byteParts[1], []byte("sys")) || !bytes.Equal(byteParts[2], []byte("default.skn")) || len(byteSplitTwo) != 2 || !bytes.Equal(byteSplitTwo[1], []byte("sys/default.skn")) || len(byteFields) != 3 || !bytes.Equal(byteFields[0], []byte("alpha")) || !bytes.Equal(byteFields[2], []byte("gamma")) || !bytes.Equal(byteTrimmed, []byte("default")) || !bytes.Equal(byteReplaced, []byte("/sys/default.txt")) {
		return checkResult{
			label:  "builders",
			ok:     false,
			detail: "bytes helper mismatch",
		}
	}
	if byteHeadErr != nil || byteHeadValueErr != nil || byteUnreadErr != nil || byteHeadValueAgainErr != nil || byteSeekErr != nil || byteTailErr != nil || byteReadAtErr != nil || byteCopyErr != nil || byteHeadRead != 4 || !bytes.Equal(byteHead[:byteHeadRead], []byte("/sys")) || byteHeadValue != '/' || byteHeadValueAgain != '/' || byteSeekPos != 12 || byteTailRead != 4 || !bytes.Equal(byteTail[:byteTailRead], []byte(".skn")) || byteReadAtCount != 7 || !bytes.Equal(byteReadAt[:byteReadAtCount], []byte("default")) || byteReaderLen != 0 || byteReaderSize != int64(len(diagFilesProbePath)) || byteCopied != int64(len(diagFilesProbePath)) || !bytes.Equal(byteCopyBuffer.Bytes(), []byte(diagFilesProbePath)) {
		return checkResult{
			label:  "builders",
			ok:     false,
			detail: "bytes reader mismatch",
		}
	}

	return checkResult{
		label:  "builders",
		ok:     true,
		detail: "builder buffer reader split fields trim replace / len " + formatInt(builderLen) + " / cap " + formatInt(bufferCap),
	}
}

func checkStrconv() checkResult {
	info, err := os.Stat(diagFilesProbePath)
	if err != nil {
		return checkResult{
			label:  "strconv",
			ok:     false,
			detail: "stat failed: " + err.Error(),
		}
	}
	rawInfo, ok := info.Sys().(kos.FileInfo)
	if !ok {
		return checkResult{
			label:  "strconv",
			ok:     false,
			detail: "stat sys payload mismatch",
		}
	}
	currentFolder, err := os.Getwd()
	if err != nil {
		return checkResult{
			label:  "strconv",
			ok:     false,
			detail: "getwd failed: " + err.Error(),
		}
	}

	formatBool := strconv.FormatBool(true)
	formatInt := strconv.Itoa(-42)
	formatHex := strconv.FormatInt(-42, 16)
	formatUint := strconv.FormatUint(uint64(info.Size()), 16)

	parseBool, parseBoolErr := strconv.ParseBool("TRUE")
	parseInt, parseIntErr := strconv.Atoi("214")
	parseHex, parseHexErr := strconv.ParseInt("-0x2a", 0, 32)
	parseBin, parseBinErr := strconv.ParseUint("0b1010", 0, 32)

	appendInt := string(strconv.AppendInt([]byte("n="), -42, 10))
	appendUint := string(strconv.AppendUint([]byte("h="), uint64(info.Size()), 16))
	appendBool := string(strconv.AppendBool([]byte("ok="), true))

	_, rangeErr := strconv.ParseUint("999", 10, 8)
	_, syntaxErr := strconv.ParseBool("maybe")

	if formatBool != "true" || formatInt != "-42" || formatHex != "-2a" || formatUint == "" {
		return checkResult{
			label:  "strconv",
			ok:     false,
			detail: "format mismatch",
		}
	}
	if parseBoolErr != nil || !parseBool || parseIntErr != nil || parseInt != 214 {
		return checkResult{
			label:  "strconv",
			ok:     false,
			detail: "parse bool/int mismatch",
		}
	}
	if parseHexErr != nil || parseHex != -42 || parseBinErr != nil || parseBin != 10 {
		return checkResult{
			label:  "strconv",
			ok:     false,
			detail: "parse base mismatch",
		}
	}
	if appendInt != "n=-42" || appendUint != "h="+formatUint || appendBool != "ok=true" {
		return checkResult{
			label:  "strconv",
			ok:     false,
			detail: "append mismatch",
		}
	}
	if !errors.Is(rangeErr, strconv.ErrRange) || !errors.Is(syntaxErr, strconv.ErrSyntax) {
		return checkResult{
			label:  "strconv",
			ok:     false,
			detail: "error mismatch",
		}
	}

	return checkResult{
		label:  "strconv",
		ok:     true,
		detail: "format parse append / cwd " + currentFolder + " / attrs 0x" + strconv.FormatUint(uint64(rawInfo.Attributes), 16),
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
	info, err := os.Stat(diagFilesProbePath)
	if err != nil {
		return checkResult{
			label:  "files",
			ok:     false,
			detail: "stat " + diagFilesProbePath + " / " + err.Error(),
		}
	}
	rawInfo, ok := info.Sys().(kos.FileInfo)
	if !ok {
		return checkResult{
			label:  "files",
			ok:     false,
			detail: "stat sys payload mismatch",
		}
	}

	previewSize := diagPreviewBytes
	if info.Size() > 0 && info.Size() < int64(previewSize) {
		previewSize = int(info.Size())
	}
	if previewSize == 0 {
		previewSize = diagPreviewBytes
	}

	file, err := os.Open(diagFilesProbePath)
	if err != nil {
		return checkResult{
			label:  "files",
			ok:     false,
			detail: "open " + diagFilesProbePath + " / " + err.Error(),
		}
	}

	buffer := make([]byte, previewSize)
	read, err := file.Read(buffer)
	closeErr := file.Close()
	if err == nil {
		err = closeErr
	}
	if err != nil && !errors.Is(err, io.EOF) {
		return checkResult{
			label:  "files",
			ok:     false,
			detail: "read " + diagFilesProbePath + " / " + err.Error(),
		}
	}

	return checkResult{
		label: "files",
		ok:    true,
		detail: "size " + formatHex64(uint64(info.Size())) +
			" / attrs " + formatHex64(uint64(rawInfo.Attributes)) +
			" / head " + formatBytePreview(buffer[:int(read)]),
	}
}

func checkOS() checkResult {
	base := diagOSProbeRoot
	if _, err := os.Stat(base); err != nil {
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
	readerInfo, err := reader.Stat()
	if err != nil {
		_ = reader.Close()
		return checkResult{
			label:  "os",
			ok:     false,
			detail: "file stat " + err.Error(),
		}
	}
	headAt := make([]byte, len(payloadBase))
	headAtCount, headAtErr := reader.ReadAt(headAt, 0)
	if headAtErr != nil && headAtErr != io.EOF {
		_ = reader.Close()
		return checkResult{
			label:  "os",
			ok:     false,
			detail: "readat " + headAtErr.Error(),
		}
	}
	seekPos, seekErr := reader.Seek(-int64(len(payloadExtra)), io.SeekEnd)
	if seekErr != nil {
		_ = reader.Close()
		return checkResult{
			label:  "os",
			ok:     false,
			detail: "seek end " + seekErr.Error(),
		}
	}
	tail := make([]byte, len(payloadExtra))
	tailRead, tailErr := reader.Read(tail)
	if tailErr != nil && tailErr != io.EOF {
		_ = reader.Close()
		return checkResult{
			label:  "os",
			ok:     false,
			detail: "seek read " + tailErr.Error(),
		}
	}
	restartPos, restartErr := reader.Seek(0, io.SeekStart)
	if restartErr != nil {
		_ = reader.Close()
		return checkResult{
			label:  "os",
			ok:     false,
			detail: "seek start " + restartErr.Error(),
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
	if headAtCount != len(payloadBase) || string(headAt[:headAtCount]) != payloadBase {
		return checkResult{
			label:  "os",
			ok:     false,
			detail: "readat mismatch",
		}
	}
	if seekPos != int64(len(payloadBase)) || tailRead != len(payloadExtra) || string(tail[:tailRead]) != payloadExtra || restartPos != 0 {
		return checkResult{
			label:  "os",
			ok:     false,
			detail: "seek mismatch",
		}
	}
	if readerInfo.Size() != int64(len(payload)) {
		return checkResult{
			label:  "os",
			ok:     false,
			detail: "file stat size mismatch",
		}
	}

	info, err := os.Stat(demoFile)
	if err != nil {
		return checkResult{
			label:  "os",
			ok:     false,
			detail: "stat " + err.Error(),
		}
	}
	modTime := info.ModTime()
	if modTime.IsZero() || modTime.Year() < 2000 {
		return checkResult{
			label:  "os",
			ok:     false,
			detail: "modtime unavailable",
		}
	}
	if os.Getpid() <= 0 {
		return checkResult{
			label:  "os",
			ok:     false,
			detail: "getpid failed",
		}
	}
	os.Clearenv()
	if len(os.Environ()) != 0 {
		return checkResult{
			label:  "os",
			ok:     false,
			detail: "clearenv mismatch",
		}
	}
	if err := os.Setenv("GODIAG_ENV", "ok"); err != nil {
		return checkResult{
			label:  "os",
			ok:     false,
			detail: "setenv " + err.Error(),
		}
	}
	envValue, envOK := os.LookupEnv("GODIAG_ENV")
	envList := os.Environ()
	if !envOK || envValue != "ok" || len(envList) != 1 || envList[0] != "GODIAG_ENV=ok" {
		return checkResult{
			label:  "os",
			ok:     false,
			detail: "environment mismatch",
		}
	}
	if err := os.Unsetenv("GODIAG_ENV"); err != nil {
		return checkResult{
			label:  "os",
			ok:     false,
			detail: "unsetenv " + err.Error(),
		}
	}
	if value, ok := os.LookupEnv("GODIAG_ENV"); ok || value != "" || os.Getenv("GODIAG_ENV") != "" {
		return checkResult{
			label:  "os",
			ok:     false,
			detail: "unsetenv mismatch",
		}
	}
	if len(os.Args) < 1 {
		return checkResult{
			label:  "os",
			ok:     false,
			detail: "args bootstrap mismatch",
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
		detail: "cwd " + base + " / pid " + formatInt(os.Getpid()) + " / mod " + formatTimeStamp(modTime) + " / readat seek env append rename cleanup",
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
		label: "dll",
		ok:    true,
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
	if _, err := kos.WriteActiveConsole([]byte("active console preflight\n")); err != nil {
		_ = console.Close()
		return checkResult{
			label:  "console",
			ok:     false,
			detail: "active console bridge failed: " + err.Error(),
		}
	}

	if _, err := fmt.Println("golang-kolibrios console probe"); err != nil {
		_ = console.Close()
		return checkResult{
			label:  "console",
			ok:     false,
			detail: "stdout fmt header failed: " + err.Error(),
		}
	}
	if _, err := fmt.Printf("stdout fmt path active / table 0x%x / ver 0x%x\n", uint32(console.ExportTable()), console.Version()); err != nil {
		_ = console.Close()
		return checkResult{
			label:  "console",
			ok:     false,
			detail: "stdout fmt body failed: " + err.Error(),
		}
	}
	if _, err := fmt.Fprintf(console, "direct writer path active\n"); err != nil {
		_ = console.Close()
		return checkResult{
			label:  "console",
			ok:     false,
			detail: "writer fmt body failed: " + err.Error(),
		}
	}

	scanState := "line input missing"
	if console.SupportsLineInput() {
		_, _ = fmt.Println("stdin fmt scan path available for manual console demo")
		scanState = "line input ready"
	}

	_ = console.Close()
	return checkResult{
		label:  "console",
		ok:     true,
		detail: "init stdout fmt writer stdin exit / " + titleState + " / " + scanState,
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
	if err == nil || os.IsNotExist(err) {
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
