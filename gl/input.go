package gl

import (
	//"fmt"
	//"log"

	//"bytes"
	//"math"
	//"strconv"

	//"encoding/binary"

	//"github.com/corpusc/viscript/gfx"
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
	w.SetCharCallback(onChar)
	w.SetKeyCallback(onKey)
	w.SetMouseButtonCallback(onMouseButton)
	w.SetScrollCallback(onMouseScroll)
	w.SetCursorPosCallback(onMouseCursorPos)
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
	//DispatchEvent(msg.TypeMouseButton, m)
	b := msg.Serialize(msg.TypeMouseButton, m)
	InputEvents <- b

}

// triggered both by moving **AND*** by pressing buttons
func onMouseCursorPos(w *glfw.Window, x float64, y float64) {
	var m msg.MessageMousePos
	m.X = x
	m.Y = y
	//DispatchEvent(msg.TypeMousePos, m)
	b := msg.Serialize(msg.TypeMousePos, m)
	InputEvents <- b
}

func onMouseScroll(w *glfw.Window, xOff, yOff float64) {
	var m msg.MessageMouseScroll
	m.X = xOff
	m.Y = yOff
	//DispatchEvent(msg.TypeMouseScroll, m)
	b := msg.Serialize(msg.TypeMouseScroll, m)
	InputEvents <- b
}

func onChar(w *glfw.Window, char rune) {
	var m msg.MessageOnCharacter
	m.Rune = uint32(char)
	b := msg.Serialize(msg.TypeChar, m)
	InputEvents <- b
}

func onKey(
	w *glfw.Window,
	key glfw.Key,
	scancode int,
	action glfw.Action,
	mod glfw.ModifierKey) {

	var m msg.MessageKey
	m.Key = uint8(key)
	m.Scan = uint32(scancode)
	m.Action = uint8(action)
	m.Mod = uint8(mod)

	//DispatchEvent(msg.TypeKey, m)
	b := msg.Serialize(msg.TypeKey, m)
	InputEvents <- b
}
