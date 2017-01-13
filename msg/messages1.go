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

const MP1 uint16 = 0x0100 //input message prefix

// HyperVisor -> Process, input event
const (
	TypeMousePos        = 1 + MP1
	TypeMouseScroll     = 2 + MP1
	TypeMouseButton     = 3 + MP1
	TypeChar            = 4 + MP1
	TypeKey             = 5 + MP1
	TypeFrameBufferSize = 6 + MP1
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
	Key    uint8
	Scan   uint32
	Action uint8
	Mod    uint8
}

// Terminal Driving Messages
type MessageFrameBufferSize struct {
	X uint32
	Y uint32
}
