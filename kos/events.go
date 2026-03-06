package kos

func WaitEvent() EventType {
	return EventType(Event())
}

func PollEvent() EventType {
	return EventType(CheckEvent())
}

func WaitEventFor(timeout uint32) EventType {
	return EventType(WaitEventTimeout(timeout))
}

func SwapEventMask(mask EventMask) EventMask {
	return EventMask(SetEventMask(uint32(mask)))
}

func CurrentButtonID() ButtonID {
	return ButtonID(GetButtonID())
}
