package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	//"github.com/go-gl/glfw/v3.1/glfw"
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
		showUInt8("Button", message)
		showUInt8("Action", message)
		showUInt8("Mod", message)

	case MessageCharacter:
		s("MessageCharacter", message)
		showRune("Rune", message)

	case MessageKey:
		s("MessageKey", message)
		showUInt8("Key", message)
		showSInt32("Scan", message)
		showUInt8("Action", message)
		showUInt8("Mod", message)

	default:
		fmt.Println("UNKNOWN MESSAGE TYPE!")

	}

	fmt.Println()
	curRecByte = 0
}

func s(s string, message []byte) { // common to all messages
	fmt.Print(s)
	showUInt32("Len", message)
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

func showRune(s string, message []byte) { // almost generic, just top 2 vars customized (and string format)
	var value rune
	var size = 4

	rBuf := bytes.NewReader(message[curRecByte : curRecByte+size])
	err := binary.Read(rBuf, binary.LittleEndian, &value)
	curRecByte += size

	if err != nil {
		fmt.Println("binary.Read failed: ", err)
	} else {
		fmt.Printf("   [%s: %s]", s, string(value))
	}
}

func showUInt8(s string, message []byte) { // almost generic, just top 2 vars customized (and string format)
	var value uint8
	var size = 1

	rBuf := bytes.NewReader(message[curRecByte : curRecByte+size])
	err := binary.Read(rBuf, binary.LittleEndian, &value)
	curRecByte += size

	if err != nil {
		fmt.Println("binary.Read failed: ", err)
	} else {
		fmt.Printf("   [%s: %d]", s, value)
	}
}

func showSInt32(s string, message []byte) { // almost generic, just top 2 vars customized (and string format)
	var value int32
	var size = 4

	rBuf := bytes.NewReader(message[curRecByte : curRecByte+size])
	err := binary.Read(rBuf, binary.LittleEndian, &value)
	curRecByte += size

	if err != nil {
		fmt.Println("binary.Read failed: ", err)
	} else {
		fmt.Printf("   [%s: %d]", s, value)
	}
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
		fmt.Printf("   [%s: %d]", s, value)
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
		fmt.Printf("   [%s: %.1f]", s, value)
	}
}

//		b glfw.MouseButton
//		action glfw.Action
