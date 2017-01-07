package testing

import (
	"fmt"
	"github.com/corpusc/viscript/hypervisor"
	"github.com/corpusc/viscript/msg"
	"testing"
)

func TestRuneMsgFlow(t *testing.T) {

	var m msg.MessageOnCharacter
	m.Rune = uint32(107)
	b := msg.Serialize(msg.TypeOnCharacter, m)
	hypervisor.InputEvents <- b

	message := hypervisor.DispatchInputEvents(hypervisor.InputEvents)

	var msgChar msg.MessageOnCharacter
	msg.MustDeserialize(message, &msgChar)

	runePrint := hypervisor.InsertRuneIntoDocument("Rune",msgChar.Rune)

	fmt.Printf("\n",runePrint)

	if msgChar.Rune != m.Rune {
		t.Error("Test rune msg flow not passed\n")
	} else {
		fmt.Println("rune test msg flow passed\n")
	}
}
