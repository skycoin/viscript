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

func ProcessMessage(message []byte) {
	// do: extract length
	// do: get msgType

	/*switch msgType {
	case MessageMousePos:
	*/
	fmt.Println("MessageMousePos")

	showFloat64("X", message)
	showFloat64("Y", message)
	/*case MessageMouseScroll:
		fmt.Println("MessageMouseScroll")
	case MessageMouseButton:
		fmt.Println("MessageMouseButton")
	default:
		fmt.Println("UNKNOWN MESSAGE TYPE!")
	}
	*/

	curRecByte = 0
}

func showFloat64(s string, message []byte) {
	var value float64

	rBuf := bytes.NewReader(message[curRecByte : curRecByte+8]) // later _:_+8
	err := binary.Read(rBuf, binary.LittleEndian, &value)
	curRecByte += 8

	if err != nil {
		fmt.Println("binary.Read failed: ", err)
	} else {
		fmt.Printf("from byte buffer, %s: %.1f\n", s, value)
	}
}
