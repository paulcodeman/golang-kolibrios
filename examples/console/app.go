package consoledemo

import (
	"fmt"

	"kos"
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

	if _, err := fmt.Println("KolibriOS console demo"); err != nil {
		kos.DebugString("console demo: stdout write failed")
		kos.Exit()
		return
	}
	_, _ = fmt.Printf("Loaded %s and resolved required exports.\n", kos.ConsoleDLLPath)
	_, _ = fmt.Println("fmt.Print* now routes through os.Stdout into CONSOLE.OBJ.")
	_, _ = fmt.Printf("export table: 0x%x / version: 0x%x\n", uint32(console.ExportTable()), console.Version())
	if console.SupportsTitle() {
		console.SetTitle(consoleDemoTitle + " / ready")
	}

	if console.SupportsInput() {
		_, _ = fmt.Println("Press Esc to close this console.")
		waitForConsoleExitKey(console)
	} else {
		_, _ = fmt.Println("Input export missing, closing in three seconds.")
		kos.SleepSeconds(3)
	}

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
