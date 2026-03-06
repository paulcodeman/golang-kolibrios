package kos

type Color uint32
type ButtonID int
type EventType int
type EventMask uint32
type WindowState byte
type ThreadStatus uint16
type KeyboardMode byte

type Point struct {
	X int
	Y int
}

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
	EventMaskRedraw EventMask = 1 << 0
	EventMaskKey EventMask = 1 << 1
	EventMaskButton EventMask = 1 << 2
	EventMaskDesktop EventMask = 1 << 4
	EventMaskMouse EventMask = 1 << 5
	EventMaskIPC EventMask = 1 << 6
	EventMaskNetwork EventMask = 1 << 7
	EventMaskDebug EventMask = 1 << 8
	EventMaskMouseInsideWindowOnly EventMask = 1 << 30
	EventMaskMouseActiveWindowOnly EventMask = 1 << 31
	DefaultEventMask EventMask = EventMaskRedraw | EventMaskKey | EventMaskButton
)

const (
	WindowStateMaximized WindowState = 1 << 0
	WindowStateMinimized WindowState = 1 << 1
	WindowStateRolledUp  WindowState = 1 << 2
)

const (
	ThreadRunning              ThreadStatus = 0
	ThreadSuspended            ThreadStatus = 1
	ThreadSuspendedWaiting     ThreadStatus = 2
	ThreadTerminating          ThreadStatus = 3
	ThreadExceptionTerminating ThreadStatus = 4
	ThreadWaitingForEvent      ThreadStatus = 5
	ThreadSlotFree             ThreadStatus = 9
)

const (
	KeyboardASCII KeyboardMode = 0
	KeyboardScan  KeyboardMode = 1
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
