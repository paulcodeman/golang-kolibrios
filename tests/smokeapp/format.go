package smokeapp

import "../../kos"

var smokeDigits = [...]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

func formatInt(value int) string {
	if value < 0 {
		return "-" + formatUint32(uint32(-value))
	}

	return formatUint32(uint32(value))
}

func formatUint32(value uint32) string {
	if value < 10 {
		return smokeDigits[value]
	}

	return formatUint32(value/10) + smokeDigits[value%10]
}

func formatByte(value byte) string {
	if value < 10 {
		return smokeDigits[value]
	}

	return smokeDigits[value/10] + smokeDigits[value%10]
}

func formatClock(value kos.ClockTime) string {
	return formatByte(value.Hour) + ":" + formatByte(value.Minute) + ":" + formatByte(value.Second)
}

func formatKernelVersion(value kos.KernelVersionInfo) string {
	return formatUint32(uint32(value.Major)) + "." +
		formatUint32(uint32(value.Minor)) + "." +
		formatUint32(uint32(value.Patch)) + "." +
		formatUint32(uint32(value.Build))
}
