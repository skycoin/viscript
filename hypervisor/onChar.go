package hypervisor

import (
	"github.com/corpusc/viscript/msg"
)

func onChar(m msg.MessageOnCharacter) {
	InsertRuneIntoDocument("Rune", m.Rune)
	//script.Process(false)
}
