package sysinfo

import "../../kos"

var decimalDigits = [...]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
var hexDigits = [...]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F"}

func formatInt(value int) string {
	if value < 0 {
		return "-" + formatUint32(uint32(-value))
	}

	return formatUint32(uint32(value))
}

func formatUint32(value uint32) string {
	if value < 10 {
		return decimalDigits[value]
	}

	return formatUint32(value/10) + decimalDigits[value%10]
}

func formatHex8(value byte) string {
	return "0x" +
		hexDigits[(value>>4)&0x0F] +
		hexDigits[value&0x0F]
}

func formatHex32(value uint32) string {
	return "0x" +
		hexDigits[(value>>28)&0x0F] +
		hexDigits[(value>>24)&0x0F] +
		hexDigits[(value>>20)&0x0F] +
		hexDigits[(value>>16)&0x0F] +
		hexDigits[(value>>12)&0x0F] +
		hexDigits[(value>>8)&0x0F] +
		hexDigits[(value>>4)&0x0F] +
		hexDigits[value&0x0F]
}

func formatKernelVersion(info kos.KernelVersionInfo) string {
	return formatUint32(uint32(info.Major)) + "." +
		formatUint32(uint32(info.Minor)) + "." +
		formatUint32(uint32(info.Patch)) + "." +
		formatUint32(uint32(info.Build))
}

func formatKernelABI(info kos.KernelVersionInfo) string {
	return formatUint32(uint32(info.ABIMajor)) + "." +
		formatUint32(uint32(info.ABIMinor))
}

func formatRect(rect kos.Rect) string {
	return "(" + formatInt(rect.Left) + "," + formatInt(rect.Top) + ")-(" +
		formatInt(rect.Right) + "," + formatInt(rect.Bottom) + ")"
}
