package gl

import (
	"github.com/corpusc/viscript/msg"
	"github.com/go-gl/glfw/v3.2/glfw"
)

var InputEvents = make(chan []byte, 256) //event channel

// push events to the event queue
func PollEvents() {
	glfw.PollEvents() //move to gl
}

func InitInputEvents(w *glfw.Window) {
	//ui
	w.SetCloseCallback(onClose)
	//keyboard
	w.SetCharCallback(onChar)
	w.SetKeyCallback(onKey)
	//mouse
	w.SetMouseButtonCallback(onMouseButton)
	w.SetScrollCallback(onMouseScroll)
	w.SetCursorPosCallback(onMouseCursorPos)
}

func SerializeAndDispatch(msgType uint16, message interface{}) {
	// Send byte slice to channel
	InputEvents <- msg.Serialize(msgType, message)
}

func onClose(w *glfw.Window) {
	SerializeAndDispatch(
		msg.TypeKey,
		msg.MessageKey{Key: msg.KeyEscape})
}

// apparently every time this is fired, a mouse position event is ALSO fired
func onMouseButton(
	w *glfw.Window,
	bt glfw.MouseButton,
	action glfw.Action,
	mod glfw.ModifierKey) {

	SerializeAndDispatch(
		msg.TypeMouseButton,
		msg.MessageMouseButton{uint8(bt), uint8(action), uint8(mod)})
}

// triggered both by moving **AND*** by pressing buttons
func onMouseCursorPos(w *glfw.Window, x float64, y float64) {
	SerializeAndDispatch(
		msg.TypeMousePos,
		msg.MessageMousePos{x, y})
}

func onMouseScroll(w *glfw.Window, xOff, yOff float64) {
	SerializeAndDispatch(
		msg.TypeMouseScroll,
		msg.MessageMouseScroll{xOff, yOff, eitherControlKeyHeld(w)})
}

func onChar(w *glfw.Window, char rune) {
	SerializeAndDispatch(
		msg.TypeChar,
		msg.MessageChar{uint32(char)})
}

func onKey(
	w *glfw.Window,
	key glfw.Key,
	scancode int,
	action glfw.Action,
	mod glfw.ModifierKey) {

	SerializeAndDispatch(
		msg.TypeKey,
		msg.MessageKey{uint32(key), uint32(scancode), uint8(action), uint8(mod)})
}

func eitherControlKeyHeld(w *glfw.Window) bool {
	if w.GetKey(glfw.KeyLeftControl) == glfw.Press || w.GetKey(glfw.KeyRightControl) == glfw.Press {
		return true
	}

	return false
}
