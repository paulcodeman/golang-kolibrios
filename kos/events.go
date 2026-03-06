package kos

func WaitEvent() EventType {
	return EventType(Event())
}

func CurrentButtonID() ButtonID {
	return ButtonID(GetButtonID())
}
