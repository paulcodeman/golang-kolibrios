package kos

const ConsoleDefaultDimension = ^uint32(0)
const ConsoleDLLStart = 1

type Console struct {
	table           DLLExportTable
	startProc       DLLProc
	initProc        DLLProc
	writeStringProc DLLProc
	exitProc        DLLProc
	setTitleProc    DLLProc
	getchProc       DLLProc
	keyHitProc      DLLProc
	version         uint32
}

var activeConsoleTable DLLExportTable
var activeConsoleExitProc DLLProc

func LoadConsole() (Console, bool) {
	return LoadConsoleFromDLL(LoadConsoleDLL())
}

func LoadConsoleFromDLL(table DLLExportTable) (Console, bool) {
	console := Console{
		table:           table,
		startProc:       table.Lookup("START"),
		initProc:        table.Lookup("con_init"),
		writeStringProc: table.Lookup("con_write_string"),
		exitProc:        table.Lookup("con_exit"),
		setTitleProc:    table.Lookup("con_set_title"),
		getchProc:       table.Lookup("con_getch"),
		keyHitProc:      table.Lookup("con_kbhit"),
		version:         uint32(table.Lookup("version")),
	}
	if !console.Valid() {
		return Console{}, false
	}
	console.start()

	return console, true
}

func OpenConsole(title string) (Console, bool) {
	console, ok := LoadConsole()
	if !ok {
		return Console{}, false
	}
	if !console.InitDefault(title) {
		return Console{}, false
	}

	return console, true
}

func (console Console) ExportTable() DLLExportTable {
	return console.table
}

func (console Console) Valid() bool {
	return console.table != 0 &&
		console.initProc.Valid() &&
		console.writeStringProc.Valid() &&
		console.exitProc.Valid()
}

func (console Console) SupportsTitle() bool {
	return console.setTitleProc.Valid()
}

func (console Console) Version() uint32 {
	return console.version
}

func (console Console) SupportsInput() bool {
	return console.getchProc.Valid()
}

func (console Console) start() {
	if console.startProc.Valid() {
		CallStdcall1VoidRaw(uint32(console.startProc), ConsoleDLLStart)
	}
}

func (console Console) Init(windowWidth uint32, windowHeight uint32, scrollWidth uint32, scrollHeight uint32, title string) bool {
	titlePtr, titleAddr := stringAddress(title)
	if !console.Valid() || titlePtr == nil {
		return false
	}

	CallStdcall5VoidRaw(uint32(console.initProc), windowWidth, windowHeight, scrollWidth, scrollHeight, titleAddr)
	freeCString(titlePtr)
	registerActiveConsole(console)
	return true
}

func (console Console) InitDefault(title string) bool {
	return console.Init(
		ConsoleDefaultDimension,
		ConsoleDefaultDimension,
		ConsoleDefaultDimension,
		ConsoleDefaultDimension,
		title,
	)
}

func (console Console) SetTitle(title string) bool {
	titlePtr, titleAddr := stringAddress(title)
	if !console.SupportsTitle() || titlePtr == nil {
		return false
	}

	CallStdcall1VoidRaw(uint32(console.setTitleProc), titleAddr)
	freeCString(titlePtr)
	return true
}

func (console Console) WriteString(text string) bool {
	textPtr, textAddr := stringAddress(text)
	if !console.Valid() || textPtr == nil {
		return false
	}

	CallStdcall2VoidRaw(uint32(console.writeStringProc), textAddr, uint32(len(text)))
	freeCString(textPtr)
	return true
}

func (console Console) Write(data []byte) (int, error) {
	if len(data) == 0 {
		return 0, nil
	}
	if console.WriteString(string(data)) {
		return len(data), nil
	}

	return 0, &consoleError{text: "console write failed"}
}

func (console Console) KeyHit() bool {
	return console.keyHitProc.Valid() && CallStdcall0Raw(uint32(console.keyHitProc)) != 0
}

func (console Console) Getch() int {
	if !console.SupportsInput() {
		return 0
	}

	return int(int32(CallStdcall0Raw(uint32(console.getchProc))))
}

func (console Console) Close() error {
	if !console.Valid() {
		return &consoleError{text: "console close failed"}
	}

	console.Exit(true)
	return nil
}

func (console Console) Exit(closeWindow bool) {
	if !console.exitProc.Valid() {
		return
	}

	CallStdcall1VoidRaw(uint32(console.exitProc), boolToUint32(closeWindow))
	unregisterActiveConsole(console)
}

func boolToUint32(value bool) uint32 {
	if value {
		return 1
	}

	return 0
}

type consoleError struct {
	text string
}

func (err *consoleError) Error() string {
	return err.text
}

func registerActiveConsole(console Console) {
	activeConsoleTable = console.table
	activeConsoleExitProc = console.exitProc
}

func unregisterActiveConsole(console Console) {
	if activeConsoleTable == console.table {
		activeConsoleTable = 0
		activeConsoleExitProc = 0
	}
}

func closeActiveConsole(closeWindow bool) {
	if activeConsoleTable == 0 || !activeConsoleExitProc.Valid() {
		return
	}

	CallStdcall1VoidRaw(uint32(activeConsoleExitProc), boolToUint32(closeWindow))
	activeConsoleTable = 0
	activeConsoleExitProc = 0
}
