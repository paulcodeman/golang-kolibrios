package consoledemo

import (
	"fmt"

	"../../kos"
)

const consoleDemoTitle = "KolibriOS Console Demo"
const consoleExitKey = 27

func Main() {
	console, ok := kos.OpenConsole(consoleDemoTitle)
	if !ok {
		kos.DebugString("console demo: failed to open /sys/lib/console.obj")
		kos.Exit()
		return
	}

	if _, err := fmt.Fprintf(console, "KolibriOS console demo\r\n"); err != nil {
		kos.DebugString("console demo: fmt write failed")
		console.Exit(true)
		kos.Exit()
		return
	}
	_, _ = fmt.Fprintf(console, "Loaded %s and resolved required exports.\r\n", kos.ConsoleDLLPath)
	_, _ = fmt.Fprintf(console, "fmt.Fprintf now writes through CONSOLE.OBJ as an io.Writer.\r\n")
	_, _ = fmt.Fprintf(console, "export table: 0x%x / version: 0x%x\r\n", uint32(console.ExportTable()), console.Version())
	if console.SupportsTitle() {
		console.SetTitle(consoleDemoTitle + " / ready")
	}

	if console.SupportsInput() {
		_, _ = fmt.Fprintf(console, "Press Esc to close this console.\r\n")
		waitForConsoleExitKey(console)
	} else {
		_, _ = fmt.Fprintf(console, "Input export missing, closing in three seconds.\r\n")
		kos.SleepSeconds(3)
	}

	console.Exit(true)
	kos.Exit()
}

func waitForConsoleExitKey(console kos.Console) {
	for {
		key := console.Getch()
		if key == consoleExitKey {
			return
		}
	}
}
