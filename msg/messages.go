package msg

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/corpusc/viscript/gfx"
	"github.com/corpusc/viscript/script"
)

var curRecByte = 0 // current receive message index

func MonitorEvents(ch chan []byte) {
	select {
	case v := <-ch:
		processMessage(v)
	default:
		//fmt.Println("monitorEvents() default")
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
	//		msg.Deserialize(msg, m1)
	//		//self.HandleInRouteMessage(m1)
	//		fmt.Printf("TypeMousePos: X= %f, Y= %f \n", m1.X, m2.X)
	//	case msg.MsgTypeMouseScroll:
	//		var m2 msg.TypeMouseScroll
	//		msg.Deserialize(msg, m1)
	//
	//	case msg.MsgTypeMouseButton:
	//		var m3 msg.TypeMouseButton
	//		mesg.Deserialize(msg, m1)
	//
	//	case msg.MsgTypeCharacter:
	//		var m4 msg.TypeCharacter
	//		msg.Deserialize(msg, m1)
	//
	//	case msg.MsgTypeKey:
	//		var m5 msg.TypeKey
	//		msg.Deserialize(msg, m1)
	//
	//	default:
	//		fmt.Println("UNKNOWN MESSAGE TYPE!")
	//
	//	}
	//}
}

func processMessage(message []byte) {

	switch getMessageTypeUInt8(".", message) {

	case TypeMousePos:
		s("TypeMousePos", message)
		msg := MessageMousePos{}
		msg.setMessageMousePosValue(getFloat64("X", message), getFloat64("Y", message))

	case TypeMouseScroll:
		s("TypeMouseScroll", message)
		msg := MessageMouseScroll{}
		msg.setMessageMouseScrollValue(getFloat64("X Offset", message), getFloat64("Y Offset", message))

	case TypeMouseButton:
		s("TypeMouseButton", message)
		msg := MessageMouseButton{}
		msg.setMessageMouseButtonValue(
			getAndShowUInt8("Button", message),
			getAndShowUInt8("Action", message),
			getAndShowUInt8("Mod", message))
		gfx.Curs.ConvertMouseClickToTextCursorPosition(msg.Button, msg.Action)

	case TypeCharacter:
		s("TypeCharacter", message)
		insertRuneIntoDocument("Rune", message)
		script.Process(false)

	case TypeKey:
		s("TypeKey", message)
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

func s(s string, message []byte) {
	fmt.Print(s)
	showUInt32("Len", message)
	curRecByte++ // skipping message type's space
}

/*
func getMessageTypeUInt8(s string, message []byte) (value uint8) {
	rBuf := bytes.NewReader(message[4:5])
	err := binary.Read(rBuf, binary.LittleEndian, &value)

	if err != nil {
		fmt.Println("binary.Read failed: ", err)
	} else {
		//fmt.Printf("from byte buffer, %s: %d\n", s, value)
	}

	return
}
*/

func GetMessageTypeUInt16(message []byte) uint16 {
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

		f := gfx.Rend.Focused
		b := f.TextBodies[0]
		b[f.CursY] = b[f.CursY][:f.CursX] + string(value) + b[f.CursY][f.CursX:len(b[f.CursY])]
		f.CursX++
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
		fmt.Printf("coa   [%s: %d]", s, value)
	}

	return value
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

func getFloat64(s string, message []byte) (val float64) {
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

	val = value
	return val
}
