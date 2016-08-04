package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/go-gl/glfw/v3.1/glfw"
)

const (
	MessageMousePos    = iota // 0
	MessageMouseScroll        // 1
	MessageMouseButton        // 2
	MessageCharacter
	MessageKey
)

var curRecByte = 0 // current receive message index

func monitorEvents(ch chan []byte) {
	select {
	case v := <-ch:
		processMessage(v)
	default:
		//fmt.Println("monitorEvents() default")
	}
}

func processMessage(message []byte) {
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
		convertMouseClickToTextCursorPosition(
			getAndShowUInt8("Button", message),
			getAndShowUInt8("Action", message))
		getAndShowUInt8("Mod", message)

	case MessageCharacter:
		s("MessageCharacter", message)
		insertRuneIntoDocument("Rune", message)

	case MessageKey:
		s("MessageKey", message)
		getAndShowUInt8("Key", message)
		showSInt32("Scan", message)
		getAndShowUInt8("Action", message)
		getAndShowUInt8("Mod", message)

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

func convertMouseClickToTextCursorPosition(button uint8, action uint8) {
	if glfw.MouseButton(button) == glfw.MouseButtonLeft &&
		glfw.Action(action) == glfw.Press {

		if mouseY < len(document) {
			cursY = mouseY

			if mouseX <= len(document[cursY]) {
				cursX = mouseX
			} else {
				cursX = len(document[cursY])
			}
		} else {
			cursY = len(document) - 1
		}
	}
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

func insertRuneIntoDocument(s string, message []byte) {
	var value rune
	var size = 4

	rBuf := bytes.NewReader(message[curRecByte : curRecByte+size])
	err := binary.Read(rBuf, binary.LittleEndian, &value)
	curRecByte += size

	if err != nil {
		fmt.Println("binary.Read failed: ", err)
	} else {
		fmt.Printf("   [%s: %s]", s, string(value))

		document[cursY] = document[cursY][:cursX] + string(value) + document[cursY][cursX:len(document[cursY])]
		cursX++
	}
}

func getAndShowUInt8(s string, message []byte) (value uint8) {
	var size = 1

	rBuf := bytes.NewReader(message[curRecByte : curRecByte+size])
	err := binary.Read(rBuf, binary.LittleEndian, &value)
	curRecByte += size

	if err != nil {
		fmt.Println("binary.Read failed: ", err)
	} else {
		fmt.Printf("   [%s: %d]", s, value)
	}

	return
}

// the rest of these funcs are almost identical, just top 2 vars customized (and string format)
func showSInt32(s string, message []byte) {
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

func showUInt32(s string, message []byte) {
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

func showFloat64(s string, message []byte) {
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
