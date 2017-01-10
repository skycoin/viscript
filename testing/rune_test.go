package testing

import (
	"fmt"
	"testing"

	"github.com/corpusc/viscript/gl"
	"github.com/corpusc/viscript/hypervisor"
	"github.com/corpusc/viscript/msg"
)

func TestRuneMsgFlow(t *testing.T) {

	var m msg.MessageOnCharacter
	m.Rune = uint32(107)
	b := msg.Serialize(msg.TypeChar, m)

	message := flowTest(b)

	var msgChar msg.MessageOnCharacter
	msg.MustDeserialize(message, &msgChar)

	runePrint := hypervisor.InsertRuneIntoDocument("Rune", msgChar.Rune)

	fmt.Printf("\n", runePrint)

	if msgChar.Rune != m.Rune {
		t.Error("Test rune msg flow not passed\n")
	} else {
		fmt.Println("rune test msg flow passed\n")
	}
}
func flowTest(b []byte) []byte {

	gl.InputEvents <- b
	message := hypervisor.DispatchInputEvents(gl.InputEvents)

	return message
}

func TestMsgMousePosFlow(t *testing.T) {
	var m msg.MessageMousePos
	m.X = 23
	m.Y = 24

	b := msg.Serialize(msg.TypeMousePos, m)
	message := flowTest(b)

	var msgPos msg.MessageMousePos
	msg.MustDeserialize(message, &msgPos)

	if msgPos.X != m.X && msgPos.Y != m.Y {
		t.Error("Test msg mouse pos flow not passed \n")
	} else {
		fmt.Println("Test msg mouse pos flow it's passed")
	}

}
func TestMsgMouseScrollFlow(t *testing.T) {
	var m msg.MessageMouseScroll
	m.X = 0
	m.Y = -0.100006103515625

	b := msg.Serialize(msg.TypeMouseScroll, m)
	message := flowTest(b)

	var msgPos msg.MessageMouseScroll
	msg.MustDeserialize(message, &msgPos)

	if msgPos.X != m.X && msgPos.Y != m.Y {
		t.Error("Test msg mouse scroll flow not passed \n")
	} else {
		fmt.Println("Test msg mouse scroll flow it's passed")
	}

}
func TestMsgMouseButtonFlow(t *testing.T) {

	var m msg.MessageMouseButton
	m.Button = uint8(0)
	m.Action = uint8(1)
	m.Mod = uint8(0)

	b := msg.Serialize(msg.TypeMouseButton, m)
	message := flowTest(b)

	var msgBtn msg.MessageMouseButton
	msg.MustDeserialize(message, &msgBtn)

	if msgBtn.Button != m.Button && msgBtn.Action != m.Action && msgBtn.Mod != m.Mod {
		t.Error("Test msg mouse button flow not passed \n")
	} else {
		fmt.Println("Test msg mouse button flow it's passed")
	}

}
func TestMsgKeyFlow(t *testing.T) {

	var m msg.MessageKey
	m.Key = uint8(69)
	m.Scan = uint32(14)
	m.Action = uint8(1)
	m.Mod = uint8(0)

	b := msg.Serialize(msg.TypeKey, m)
	message := flowTest(b)

	var msgKey msg.MessageKey
	msg.MustDeserialize(message, &msgKey)

	if msgKey.Key != m.Key && msgKey.Scan != m.Scan && msgKey.Action != m.Action && msgKey.Mod != m.Mod {
		t.Error("Test msg msg key flow not passed \n")
	} else {
		fmt.Println("Test msg key flow it's passed")
	}

}
