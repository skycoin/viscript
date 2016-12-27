package msg

import (
	"bytes"
	"encoding/binary"
	"fmt"
	//"math/rand"
)

const (
	TypeMousePos    = iota // 0
	TypeMouseScroll        // 1
	TypeMouseButton        // 2
	TypeCharacter
	TypeKey
)
const (
	MsgTypeMousePos    = iota // HyperVisor -> Process, input event
	MsgTypeMouseScroll        //hypervisor -> Process, input event
	MsgTypeMouseButton        //hypervisor -> Process, input event
	MsgTypeCharacter          //hypervisor -> Process, input event
	MsgTypeKey                // hyperisor -> Process, input event
)

func GetMessageType(message []byte) uint16 {
	var value uint16
	rBuf := bytes.NewReader(message[4:6])
	err := binary.Read(rBuf, binary.LittleEndian, &value)
	if err != nil {
		fmt.Println("binary.Read failed: ", err)
	} else {
		//fmt.Printf("from byte buffer, %s: %d\n", s, value)
	}
	return value
}

/*
	case TypeMouseScroll:
		s("TypeMouseScroll", message)
		showFloat64("X Offset", message)
		showFloat64("Y Offset", message)

	case TypeMouseButton:
		s("TypeMouseButton", message)
		gfx.Curs.ConvertMouseClickToTextCursorPosition(
			getAndShowUInt8("Button", message),
			getAndShowUInt8("Action", message))
		getAndShowUInt8("Mod", message)

	case TypeCharacter:
		s("TypeCharacter", message)
		insertRuneIntoDocument("Rune", message)
		script.Process(false)

	case TypeKey:
		s("TypeKey", message)
		getAndShowUInt8("Key", message)
		showSInt32("Scan", message)
		getAndShowUInt8("Action", message)
		getAndShowUInt8("Mod", message)
*/

//Input Messages

// HyperVisor -> Process

//message received by process, by hypervisor
type MessageMousePos struct {
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

//Terminal Driving Messages
