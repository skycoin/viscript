package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	MessageMousePos    = iota // 0
	MessageMouseScroll        // 1
	MessageMouseButton        // 2
	MessageCharacter
	MessageKey
)

var curRecByte = 0 // current receive message index

func processMessage(message []byte) {
	// do: extract length
	// do: get msgType

	switch getMessageType(".", message) {

	case MessageMousePos:
		s("MessageMousePos", message)
		showFloat64("X", message)
		showFloat64("Y", message)

	case MessageMouseScroll:
		s("MessageMouseScroll", message)
		showFloat64("X Offset", message)
		showFloat64("Y Offset", message)

	case MessageMouseButton:
		s("MessageMouseButton", message)
		//showMouseButton()
		//showAction()

	case MessageCharacter:
		s("MessageCharacter", message)

	case MessageKey:
		s("MessageKey", message)

	default:
		fmt.Println("UNKNOWN MESSAGE TYPE!")

	}

	fmt.Println()
	curRecByte = 0
}

func s(s string, message []byte) { // common to all messages
	fmt.Print(s)
	showUInt32("Length", message)
	curRecByte++ // skipping message type's space
}

func getMessageType(s string, message []byte) (value uint8) {
	rBuf := bytes.NewReader(message[4:5])
	err := binary.Read(rBuf, binary.LittleEndian, &value)

	if err != nil {
		fmt.Println("binary.Read failed: ", err)
	} else {
		//fmt.Printf("from byte buffer, %s: %d\n", s, value)
	}

	return
}

func showUInt32(s string, message []byte) { // almost generic, just top 2 vars customized (and string format)
	var value uint32
	var size = 4

	rBuf := bytes.NewReader(message[curRecByte : curRecByte+size])
	err := binary.Read(rBuf, binary.LittleEndian, &value)
	curRecByte += size

	if err != nil {
		fmt.Println("binary.Read failed: ", err)
	} else {
		fmt.Printf("   < %s: %d >", s, value)
	}
}

func showFloat64(s string, message []byte) { // almost generic, just top 2 vars customized (and string format)
	var value float64
	var size = 8

	rBuf := bytes.NewReader(message[curRecByte : curRecByte+size])
	err := binary.Read(rBuf, binary.LittleEndian, &value)
	curRecByte += size

	if err != nil {
		fmt.Println("binary.Read failed: ", err)
	} else {
		fmt.Printf("   < %s: %.1f >", s, value)
	}
}

//		b glfw.MouseButton
//		action glfw.Action
