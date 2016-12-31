package msg

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/corpusc/viscript/gfx"
	"github.com/corpusc/viscript/script"
	//"log"
)

var curRecByte = 0 // current receive message index

func MonitorEvents(ch chan []byte) {
	select {
	case v := <-ch:
		processMessage(v)
	default:
		//fmt.Println("MonitorEvents() default")
	}
}

/* message processing example */

func ProcessIncomingMessages() {
	//have a channel for incoming messages
	//for msg := range self.IncomingChannel{
	//       print(msg)
	//}

	//for msg := range IncomingChannel {
	//	switch msg.GetMessageType(msg) {
	//	//InRouteMessage is the only message coming in to node from transports
	//	case msg.MsgTypeMousePos:
	//		var m1 msg.TypeMousePos
	//		msg.MustDeserialize(msg, m1)
	//		//self.HandleInRouteMessage(m1)
	//		fmt.Printf("TypeMousePos: X= %f, Y= %f \n", m1.X, m2.X)
	//	case msg.MsgTypeMouseScroll:
	//		var m2 msg.TypeMouseScroll
	//		msg.MustDeserialize(msg, m1)
	//
	//	case msg.MsgTypeMouseButton:
	//		var m3 msg.TypeMouseButton
	//		mesg.MustDeserialize(msg, m1)
	//
	//	case msg.MsgTypeCharacter:
	//		var m4 msg.TypeCharacter
	//		msg.MustDeserialize(msg, m1)
	//
	//	case msg.MsgTypeKey:
	//		var m5 msg.TypeKey
	//		msg.MustDeserialize(msg, m1)
	//
	//	default:
	//		fmt.Println("UNKNOWN MESSAGE TYPE!")
	//
	//	}
	//}
}

func processMessage(message []byte) {

	switch GetMessageTypeUInt16(message) {

	case TypeMousePos:
		var msgMousePos MessageMousePos
		MustDeserialize(message, &msgMousePos)

		s("TypeMousePos")
		getFloat64("X", msgMousePos.X)
		getFloat64("Y", msgMousePos.Y)

	case TypeMouseScroll:
		var msgScroll MessageMouseScroll
		MustDeserialize(message, &msgScroll)

		s("TypeMouseScroll")
		getFloat64("X Offset", msgScroll.X)
		getFloat64("Y Offset", msgScroll.Y)

	case TypeMouseButton:
		var msgBtn MessageMouseButton
		MustDeserialize(message, &msgBtn)

		s("TypeMouseButton")
		getAndShowUInt8("Button", msgBtn.Button)
		getAndShowUInt8("Action", msgBtn.Action)
		getAndShowUInt8("Mod", msgBtn.Mod)
		gfx.Curs.ConvertMouseClickToTextCursorPosition(msgBtn.Button, msgBtn.Action)

	case TypeCharacter:

		s("TypeCharacter")
		//var typeChar MessageCharacter
		//TODO aleksbgs i dont know what is this
		insertRuneIntoDocument("Rune", message)
		script.Process(false)

	case TypeKey:
		var keyMsg MessageKey
		s("TypeKey")
		MustDeserialize(message, &keyMsg)
		getAndShowUInt8("Key", keyMsg.Key)
		getUInt32Scan("Scan", keyMsg.Scan)
		getAndShowUInt8("Action", keyMsg.Action)
		getAndShowUInt8("Mod", keyMsg.Mod)

	default:
		fmt.Println("UNKNOWN MESSAGE TYPE!")
	}

	fmt.Println()
	curRecByte = 0
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

func insertRuneIntoDocument(s string, message []byte) {
	var value rune
	var size = 4

	rBuf := bytes.NewReader(message[curRecByte : curRecByte+size])
	err := binary.Read(rBuf, binary.LittleEndian, &value)
	curRecByte += size

	if err != nil {
		fmt.Println("binary.Read failed: ", err)
	} else {
		fmt.Printf("Rune   [%s: %s]", s, string(value))

		f := gfx.Rend.Focused
		b := f.TextBodies[0]
		b[f.CursY] = b[f.CursY][:f.CursX] + string(value) + b[f.CursY][f.CursX:len(b[f.CursY])]
		f.CursX++
	}
}

func getAndShowUInt8(s string, x uint8) uint8 {
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

func getUInt32Scan(s string, x uint32) uint32 {
	fmt.Printf("   [%s: %d]", s, x)
	return x
}

func getFloat64(s string, f float64) float64 {
	fmt.Printf("   [%s: %.1f]", s, f)
	return f
}
