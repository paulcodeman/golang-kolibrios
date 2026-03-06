package kos

type MouseButtonInfo struct {
	Raw              uint32
	LeftHeld         bool
	RightHeld        bool
	MiddleHeld       bool
	Button4Held      bool
	Button5Held      bool
	LeftPressed      bool
	RightPressed     bool
	MiddlePressed    bool
	LeftReleased     bool
	RightReleased    bool
	MiddleReleased   bool
	VerticalScroll   bool
	HorizontalScroll bool
	LeftDoubleClick  bool
}

func MouseScreenPosition() Point {
	return unpackUnsignedPoint(GetMouseScreenPosition())
}

func MouseWindowPosition() Point {
	return unpackSignedPackedPoint(GetMouseWindowPosition())
}

func MouseHeldButtons() MouseButtonInfo {
	return decodeMouseButtonInfo(GetMouseButtonState())
}

func MouseButtons() MouseButtonInfo {
	return decodeMouseButtonInfo(GetMouseButtonEventState())
}

func MouseScrollDelta() Point {
	return unpackSignedPackedPoint(GetMouseScrollData())
}

func decodeMouseButtonInfo(raw uint32) MouseButtonInfo {
	return MouseButtonInfo{
		Raw:              raw,
		LeftHeld:         raw&(1<<0) != 0,
		RightHeld:        raw&(1<<1) != 0,
		MiddleHeld:       raw&(1<<2) != 0,
		Button4Held:      raw&(1<<3) != 0,
		Button5Held:      raw&(1<<4) != 0,
		LeftPressed:      raw&(1<<8) != 0,
		RightPressed:     raw&(1<<9) != 0,
		MiddlePressed:    raw&(1<<10) != 0,
		VerticalScroll:   raw&(1<<15) != 0,
		LeftReleased:     raw&(1<<16) != 0,
		RightReleased:    raw&(1<<17) != 0,
		MiddleReleased:   raw&(1<<18) != 0,
		HorizontalScroll: raw&(1<<23) != 0,
		LeftDoubleClick:  raw&(1<<24) != 0,
	}
}
