package hypervisor

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/corpusc/viscript/gfx"
	"github.com/corpusc/viscript/msg"
	//"log"
	_ "strconv"
)

//DEPRECATE
//var curRecByte = 0 // current receive message index

func init() {
	//rune test
}

var DebugPrintInputEvents = false

//this is where input events are
func ProcessInputEvents(message []byte) {

	switch GetMessageTypeUInt16(message) {

	case msg.TypeMousePos:
		var msgMousePos msg.MessageMousePos
		msg.MustDeserialize(message, &msgMousePos)

		if DebugPrintInputEvents {
			s("TypeMousePos")
			showFloat64("X", msgMousePos.X)
			showFloat64("Y", msgMousePos.Y)
		}

	case msg.TypeMouseScroll:
		var msgScroll msg.MessageMouseScroll
		msg.MustDeserialize(message, &msgScroll)

		if DebugPrintInputEvents {
			s("TypeMouseScroll")
			showFloat64("X Offset", msgScroll.X)
			showFloat64("Y Offset", msgScroll.Y)
		}

	case msg.TypeMouseButton:
		var msgBtn msg.MessageMouseButton
		msg.MustDeserialize(message, &msgBtn)

		if DebugPrintInputEvents {
			s("TypeMouseButton")
			showUInt8("Button", msgBtn.Button)
			showUInt8("Action", msgBtn.Action)
			showUInt8("Mod", msgBtn.Mod)
		}
		gfx.Curs.ConvertMouseClickToTextCursorPosition(msgBtn.Button, msgBtn.Action)

	case msg.TypeCharacter:

		fmt.Printf("TypeCharacter/Rune, wtf fix\n")
		//WTF: FIX THIS
		//InsertRuneIntoDocumentTest("Rune", message)
		//script.Process(false)

		if DebugPrintInputEvents {
			s("TypeCharacter")
		}
	case msg.TypeKey:
		var keyMsg msg.MessageKey
		s("TypeKey")
		msg.MustDeserialize(message, &keyMsg)

		if DebugPrintInputEvents {
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

//WTF FIX THIS
/*
func InsertRuneIntoDocumentTest(s string, message []byte) {
	var value rune
	var size = 4

	rBuf := bytes.NewReader(message[curRecByte : curRecByte+size])
	err := binary.Read(rBuf, binary.LittleEndian, &value)
	curRecByte += size

	if err != nil {
		fmt.Println("binary.Read failed: ", err)
	} else {
		f := gfx.Rend.Focused
		b := f.TextBodies[0]
		resultsDif := f.CursX - len(b[f.CursY])
		fmt.Printf("Rune   [%s: %s]", s, string(value))

		if f.CursX > len(b[f.CursY]) {
			b[f.CursY] = b[f.CursY][:f.CursX-resultsDif] + b[f.CursY][:len(b[f.CursY])] + string(value)
			fmt.Printf("line is %s\n", b[f.CursY])
			f.CursX++
		} else {
			b[f.CursY] = b[f.CursY][:f.CursX] + string(value) + b[f.CursY][f.CursX:len(b[f.CursY])]
			f.CursX++
		}
	}
}
*/

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
