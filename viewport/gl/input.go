package gl

import (
	"fmt"

	"github.com/corpusc/viscript/msg"
	//"github.com/corpusc/viscript/script"
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
	// Send byte slice to the InputEvents chan
	InputEvents <- msg.Serialize(msgType, message)
}

func onClose(w *glfw.Window) {
	m := msg.MessageKey{Key: msg.KeyEscape}
	SerializeAndDispatch(msg.TypeKey, m)
}

// apparently every time this is fired, a mouse position event is ALSO fired
func onMouseButton(
	w *glfw.Window,
	bt glfw.MouseButton,
	action glfw.Action,
	mod glfw.ModifierKey) {

	//MessageMouseButton
	var m msg.MessageMouseButton
	m.Button = uint8(bt)
	m.Action = uint8(action)
	m.Mod = uint8(mod)

	SerializeAndDispatch(msg.TypeMouseButton, m)
}

// triggered both by moving **AND*** by pressing buttons
func onMouseCursorPos(w *glfw.Window, x float64, y float64) {
	var m msg.MessageMousePos
	m.X = x
	m.Y = y

	SerializeAndDispatch(msg.TypeMousePos, m)
}

func onMouseScroll(w *glfw.Window, xOff, yOff float64) {
	var m msg.MessageMouseScroll
	m.X = xOff
	m.Y = yOff

	SerializeAndDispatch(msg.TypeMouseScroll, m)
}

func onChar(w *glfw.Window, char rune) {
	var m msg.MessageChar
	m.Char = uint32(char)

	SerializeAndDispatch(msg.TypeChar, m)
}

func onKey(
	w *glfw.Window,
	key glfw.Key,
	scancode int,
	action glfw.Action,
	mod glfw.ModifierKey) {

	var m msg.MessageKey
	m.Key = uint32(key)
	m.Scan = uint32(scancode)
	m.Action = uint8(action)
	m.Mod = uint8(mod)

	if key != glfw.Key(m.Key) {
		fmt.Printf("ERROR KEY SERIALIZATION FUCKUP: key= %d, key= %d \n", key, m.Key)
	}

	SerializeAndDispatch(msg.TypeKey, m)
}
