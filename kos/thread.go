package kos

const CurrentThreadSlot = -1

type ThreadInfo struct {
	CPUUsage            uint32
	WindowStackPosition uint16
	WindowStackSlot     uint16
	Name                string
	ProcessAddress      uint32
	UsedMemoryMinus1    uint32
	ID                  uint32
	WindowPosition      Point
	WindowSize          Point
	Status              ThreadStatus
	ClientPosition      Point
	ClientSize          Point
	WindowState         WindowState
	EventMask           EventMask
	KeyboardMode        KeyboardMode
}

func ReadThreadInfo(slot int) (info ThreadInfo, maxSlot int, ok bool) {
	var buffer [1024]byte

	maxSlot = GetThreadInfo(&buffer[0], slot)
	if maxSlot < 0 {
		return ThreadInfo{}, maxSlot, false
	}

	info = ThreadInfo{
		CPUUsage:            littleEndianUint32(buffer[:], 0),
		WindowStackPosition: littleEndianUint16(buffer[:], 4),
		WindowStackSlot:     littleEndianUint16(buffer[:], 6),
		Name:                trimASCIIField(buffer[10:21]),
		ProcessAddress:      littleEndianUint32(buffer[:], 22),
		UsedMemoryMinus1:    littleEndianUint32(buffer[:], 26),
		ID:                  littleEndianUint32(buffer[:], 30),
		WindowPosition: Point{
			X: int(littleEndianUint32(buffer[:], 34)),
			Y: int(littleEndianUint32(buffer[:], 38)),
		},
		WindowSize: Point{
			X: int(littleEndianUint32(buffer[:], 42)),
			Y: int(littleEndianUint32(buffer[:], 46)),
		},
		Status: ThreadStatus(littleEndianUint16(buffer[:], 50)),
		ClientPosition: Point{
			X: int(littleEndianUint32(buffer[:], 54)),
			Y: int(littleEndianUint32(buffer[:], 58)),
		},
		ClientSize: Point{
			X: int(littleEndianUint32(buffer[:], 62)),
			Y: int(littleEndianUint32(buffer[:], 66)),
		},
		WindowState:  WindowState(buffer[70]),
		EventMask:    EventMask(littleEndianUint32(buffer[:], 71)),
		KeyboardMode: KeyboardMode(buffer[75]),
	}

	return info, maxSlot, true
}

func ReadCurrentThreadInfo() (info ThreadInfo, maxSlot int, ok bool) {
	return ReadThreadInfo(CurrentThreadSlot)
}

func CurrentThreadID() (id uint32, ok bool) {
	var buffer [1024]byte

	maxSlot := GetThreadInfo(&buffer[0], CurrentThreadSlot)
	if maxSlot < 0 {
		return 0, false
	}

	return littleEndianUint32(buffer[:], 30), true
}
