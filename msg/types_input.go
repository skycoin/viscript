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

const TypePrefix_Input uint16 = 0x0100 //message prefix

// HyperVisor -> Process, input event
const (
	TypeMousePos        = 1 + TypePrefix_Input
	TypeMouseScroll     = 2 + TypePrefix_Input
	TypeMouseButton     = 3 + TypePrefix_Input
	TypeChar            = 4 + TypePrefix_Input
	TypeKey             = 5 + TypePrefix_Input
	TypeFrameBufferSize = 6 + TypePrefix_Input
)

//Input Messages

// HyperVisor -> Process

//message received by process, by hypervisor
type MessageMousePos struct {
	X float64
	Y float64
}

type MessageMouseScroll struct {
	X              float64
	Y              float64
	HoldingControl bool
}

type MessageMouseButton struct {
	Button uint8
	Action uint8
	Mod    uint8
}

type MessageChar struct {
	Char uint32
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
