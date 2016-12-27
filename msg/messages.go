package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
	//"math/rand"
)

const (
	MsgMessageMousePos    = iota // HyperVisor -> Process, input event
	MessageMouseScroll           //hypervisor -> Process, input event
	MsgMessageMouseButton        //hypervisor -> Process, input event
	MsgMessageCharacter          //hypervisor -> Process, input event
	MsgMessageKey                // hyperisor -> Process, input event
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
	case MessageMouseScroll:
		s("MessageMouseScroll", message)
		showFloat64("X Offset", message)
		showFloat64("Y Offset", message)

	case MessageMouseButton:
		s("MessageMouseButton", message)
		gfx.Curs.ConvertMouseClickToTextCursorPosition(
			getAndShowUInt8("Button", message),
			getAndShowUInt8("Action", message))
		getAndShowUInt8("Mod", message)

	case MessageCharacter:
		s("MessageCharacter", message)
		insertRuneIntoDocument("Rune", message)
		script.Process(false)

	case MessageKey:
		s("MessageKey", message)
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
