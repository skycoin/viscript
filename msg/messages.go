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
	TypeOnCharacter        // etc.
	TypeKey
)

//Input Messages

// HyperVisor -> Process

//message received by process, by hypervisor

var oldMousePositionX float64
var oldMousePostionY float64

type MessageMousePos struct {
	X float64
	Y float64
}

func (m *MessageMousePos) GetDelta(currentX, currentY float64) (float64, float64) {

	xMovement := oldMousePositionX - currentX
	yMovement := oldMousePostionY - currentY
	oldMousePositionX = currentX
	oldMousePostionY = currentY

	return xMovement, yMovement
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

type MessageOnCharacter struct {
	Rune uint32
}

type MessageKey struct {
	Key    uint8
	Scan   uint32
	Action uint8
	Mod    uint8
}

//Terminal Driving Messages
