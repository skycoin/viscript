package msg

import (
	_"bytes"
	_"encoding/binary"
	_"fmt"
	//"math/rand"
)

const (
	TypeMousePos    = iota // 0
	TypeMouseScroll        // 1
	TypeMouseButton        // 2
	TypeCharacter
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



func (m * MessageMousePos)setMessageMousePosValue(x float64,y float64)  {
	m.X = x
	m.Y = y

}
func (m * MessageMouseScroll)setMessageMouseScrollValue(x float64,y float64)  {
	m.X = x
	m.Y = y

}
func (m * MessageMouseButton)setMessageMouseButtonValue(button uint8,action uint8,mod uint8)  {
	m.Button = button
	m.Action = action
	m.Mod = mod

}


//Terminal Driving Messages
