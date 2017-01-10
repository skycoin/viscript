package hypervisor

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/corpusc/viscript/msg"
	_ "strconv"
)

//DEPRECATE
//var curRecByte = 0 // current receive message index

func init() {
	//rune test
}

var DebugPrintInputEvents = false

//this is where input events are

func ProcessInputEvents(message []byte) []byte {

	switch GetMessageTypeUInt16(message) {

	case msg.TypeMousePos:
		var msgMousePos msg.MessageMousePos
		msg.MustDeserialize(message, &msgMousePos)

		if DebugPrintInputEvents {
			s("TypeMousePos")
			showFloat64("X", msgMousePos.X)
			showFloat64("Y", msgMousePos.Y)
		}

		onMouseCursorPos(msgMousePos)

	case msg.TypeMouseScroll:
		var msgScroll msg.MessageMouseScroll
		msg.MustDeserialize(message, &msgScroll)

		if DebugPrintInputEvents {
			s("TypeMouseScroll")
			showFloat64("X Offset", msgScroll.X)
			showFloat64("Y Offset", msgScroll.Y)
		}
		onMouseScroll(msgScroll)

	case msg.TypeMouseButton:
		var msgBtn msg.MessageMouseButton
		msg.MustDeserialize(message, &msgBtn)

		if DebugPrintInputEvents {
			s("TypeMouseButton")
			showUInt8("Button", msgBtn.Button)
			showUInt8("Action", msgBtn.Action)
			showUInt8("Mod", msgBtn.Mod)
		}
		onMouseButton(msgBtn)

	case msg.TypeOnCharacter:
		var msgChar msg.MessageOnCharacter
		msg.MustDeserialize(message, &msgChar)
		if DebugPrintInputEvents {
			s("TypeCharacter")
		}
		onChar(msgChar)

	case msg.TypeKey:
		var keyMsg msg.MessageKey
		msg.MustDeserialize(message, &keyMsg)
		onKey(keyMsg)

		if DebugPrintInputEvents {
			s("TypeKey")
			showUInt8("Key", keyMsg.Key)
			showUInt32Scan("Scan", keyMsg.Scan)
			showUInt8("Action", keyMsg.Action)
			showUInt8("Mod", keyMsg.Mod)
		}
	default:
		fmt.Println("UNKNOWN MESSAGE TYPE!")
	}

	fmt.Println()
	//curRecByte = 0
	return message
}

func s(s string) {
	fmt.Print(s)
}

func GetMessageTypeUInt16(message []byte) uint16 {
	var value uint16
	rBuf := bytes.NewReader(message[0:2])
	err := binary.Read(rBuf, binary.LittleEndian, &value)

	if err != nil {
		fmt.Println("binary.Read failed: ", err)
	} else {
		//fmt.Printf("from byte buffer, %s: %d\n", s, value)
	}

	return value
}

func showUInt8(s string, x uint8) uint8 {
	fmt.Printf("   [%s: %d]", s, x)
	return x
}

// the rest of these funcs are almost identical, just top 2 vars customized (and string format)
func showSInt32(s string, x int32) {
	fmt.Printf("   [%s: %d]", s, x)
}

func showUInt32(s string, x uint32) {
	fmt.Printf("   [%s: %d]", s, x)
}

func showUInt32Scan(s string, x uint32) uint32 {
	fmt.Printf("   [%s: %d]", s, x)
	return x
}

func showFloat64(s string, f float64) float64 {
	fmt.Printf("   [%s: %.1f]", s, f)
	return f
}
