package kos

type KeyEvent struct {
	Raw       int
	Empty     bool
	Hotkey    bool
	Code      byte
	ScanCode  byte
	Modifiers uint16
}

func ReadKey() KeyEvent {
	raw := GetKey()
	value := uint32(raw)
	event := KeyEvent{
		Raw: raw,
	}

	if value == 1 {
		event.Empty = true
		return event
	}

	if byte(value) == 2 {
		event.Hotkey = true
		event.ScanCode = byte(value >> 8)
		event.Modifiers = uint16(value >> 16)
		return event
	}

	event.Code = byte(value >> 8)
	event.ScanCode = byte(value >> 16)
	return event
}
