package msg

import (
	_ "bytes"
	_ "encoding/binary"
	_ "fmt"
	//"math/rand"
)

const PREFIX_SIZE = 6 // guaranteed minimum size of every message (4 for length + 2 for type)

// HyperVisor -> Process, input event
const (
	TypeMousePos    = iota // 0
	TypeMouseScroll        // 1
	TypeMouseButton        // 2
	TypeCharacter          // etc.
	TypeKey
)

//Input Messages

// HyperVisor -> Process

//message received by process, by hypervisor
type MessageMousePos struct {
	X float64
	Y float64
}

type MessageMouseScroll struct {
	X float64
	Y float64
}

type MessageMouseButton struct {
	Button uint8
	Action uint8
	Mod    uint8
}

type MessageCharacter struct {
	Rune uint16 //what type is this?
}

type MessageKey struct {
	Key    uint8
	Scan   uint32
	Action uint8
	Mod    uint8
}

func (m *MessageMousePos) setMessageMousePosValue(x, y float64) {
	m.X = x
	m.Y = y
}
func (m *MessageMouseScroll) setMessageMouseScrollValue(x, y float64) {
	m.X = x
	m.Y = y
}
func (m *MessageMouseButton) setMessageMouseButtonValue(button, action, mod uint8) {
	m.Button = button
	m.Action = action
	m.Mod = mod
}

//Terminal Driving Messages
