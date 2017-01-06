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
	Rune uint32
}

type MessageKey struct {
	Key    uint8
	Scan   uint32
	Action uint8
	Mod    uint8
}

//Terminal Driving Messages
