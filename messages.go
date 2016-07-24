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
)

var curRecByte = 0 // current receive message index

func processMessage(message []byte) {
	// do: extract length
	// do: get msgType

	switch getMessageType("Message Type", message) {
	case MessageMousePos:
		fmt.Println("MessageMousePos")
		showUInt32("Length", message)
		curRecByte++ // skipping message type's space
		showFloat64("X", message)
		showFloat64("Y", message)
	case MessageMouseScroll:
		fmt.Println("MessageMouseScroll")
	case MessageMouseButton:
		fmt.Println("MessageMouseButton")
	default:
		fmt.Println("UNKNOWN MESSAGE TYPE!")
	}

	curRecByte = 0
}

func getMessageType(s string, message []byte) (value uint8) {
	rBuf := bytes.NewReader(message[4:5])
	err := binary.Read(rBuf, binary.LittleEndian, &value)

	if err != nil {
		fmt.Println("binary.Read failed: ", err)
	} else {
		fmt.Printf("from byte buffer, %s: %d\n", s, value)
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
		fmt.Printf("from byte buffer, %s: %d\n", s, value)
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
		fmt.Printf("from byte buffer, %s: %.1f\n", s, value)
	}
}
