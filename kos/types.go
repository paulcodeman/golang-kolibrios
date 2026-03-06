package kos

type Color uint32
type ButtonID int
type EventType int

const (
	EventNone EventType = 0
	EventRedraw EventType = 1
	EventKey EventType = 2
	EventButton EventType = 3
	EventDesktop EventType = 5
	EventMouse EventType = 6
	EventIPC EventType = 7
	EventNetwork EventType = 8
	EventDebug EventType = 9
	EventIRQBegin EventType = 16
)

const (
	EVENT_NONE = int(EventNone)
	EVENT_REDRAW = int(EventRedraw)
	EVENT_KEY = int(EventKey)
	EVENT_BUTTON = int(EventButton)
	EVENT_DESKTOP = int(EventDesktop)
	EVENT_MOUSE = int(EventMouse)
	EVENT_IPC = int(EventIPC)
	EVENT_NETWORK = int(EventNetwork)
	EVENT_DEBUG = int(EventDebug)
	EVENT_IRQBEGIN = int(EventIRQBegin)
)
