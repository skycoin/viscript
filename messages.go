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

func ProcessMessage(message *bytes.Buffer) {
	// do: extract length
	// do: get msgType

	/*switch msgType {
	case MessageMousePos:
	*/
	fmt.Println("MessageMousePos")

	var xAfter float64
	rBuf := bytes.NewReader(message.Bytes())
	err := binary.Read(rBuf, binary.LittleEndian, &xAfter)

	if err != nil {
		fmt.Println("binary.Read failed: ", err)
	} else {
		//fmt.Printf("onMouseCursorPos() x: %.1f   Y: %.1f\n", x, y)
		fmt.Printf("from byte buffer, x: %.1f\n", xAfter)
	}
	/*case MessageMouseScroll:
		fmt.Println("MessageMouseScroll")
	case MessageMouseButton:
		fmt.Println("MessageMouseButton")
	default:
		fmt.Println("UNKNOWN MESSAGE TYPE!")
	}
	*/
}
