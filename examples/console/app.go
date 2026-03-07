package consoledemo

import "../../kos"

const consoleDemoTitle = "KolibriOS Console Demo"
const consoleExitKey = 27

func Main() {
	console, ok := kos.OpenConsole(consoleDemoTitle)
	if !ok {
		kos.DebugString("console demo: failed to open /sys/lib/console.obj")
		kos.Exit()
		return
	}

	console.WriteString("KolibriOS console demo\r\n")
	console.WriteString("Loaded /sys/lib/console.obj and resolved required exports.\r\n")
	console.WriteString("Console output now goes through CONSOLE.OBJ instead of screenshots.\r\n")
	if console.SupportsTitle() {
		console.SetTitle(consoleDemoTitle + " / ready")
	}

	if console.SupportsInput() {
		console.WriteString("Press Esc to close this console.\r\n")
		waitForConsoleExitKey(console)
	} else {
		console.WriteString("Input export missing, closing in three seconds.\r\n")
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
