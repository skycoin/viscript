package msg

import (
	_ "bytes"
	_ "encoding/binary"
	_ "fmt"
)

/*
Messages from:
	Opengl -> HyperVisor
	HyperVisor -> Process
*/

const PrefixInput uint16 = 0x0100 //message prefix

// HyperVisor -> Process, input event
const (
	TypeMousePos        = 1 + PrefixInput
	TypeMouseScroll     = 2 + PrefixInput
	TypeMouseButton     = 3 + PrefixInput
	TypeChar            = 4 + PrefixInput
	TypeKey             = 5 + PrefixInput
	TypeFrameBufferSize = 6 + PrefixInput
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

type MessageOnCharacter struct {
	Rune uint32
}

type MessageKey struct {
	Key    uint32
	Scan   uint32
	Action uint8
	Mod    uint8
}

// Terminal Driving Messages
type MessageFrameBufferSize struct {
	X uint32
	Y uint32
}
