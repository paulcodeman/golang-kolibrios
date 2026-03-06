package kos

type ClockTime struct {
	Hour   byte
	Minute byte
	Second byte
}

func SystemTime() ClockTime {
	value := GetTime()

	return ClockTime{
		Hour:   decodeBCDByte(byte(value)),
		Minute: decodeBCDByte(byte(value >> 8)),
		Second: decodeBCDByte(byte(value >> 16)),
	}
}

func UptimeCentiseconds() uint32 {
	return GetTimeCounter()
}

func UptimeNanoseconds() uint64 {
	return GetTimeCounterPro()
}

func SleepCentiseconds(centiseconds uint32) {
	Sleep(centiseconds)
}

func SleepMilliseconds(milliseconds uint32) {
	if milliseconds == 0 {
		return
	}

	centiseconds := milliseconds / 10
	if milliseconds%10 != 0 {
		centiseconds++
	}

	Sleep(centiseconds)
}

func SleepSeconds(seconds uint32) {
	const maxUint32 = ^uint32(0)

	centiseconds := uint64(seconds) * 100
	if centiseconds > uint64(maxUint32) {
		Sleep(maxUint32)
		return
	}

	Sleep(uint32(centiseconds))
}

func decodeBCDByte(value byte) byte {
	return ((value >> 4) * 10) + (value & 0x0F)
}
