package process

import (
	"fmt"

	"github.com/corpusc/viscript/hypervisor"
	"github.com/corpusc/viscript/msg"
)

func (st *State) publishToOut(message []byte) {
	hypervisor.DbusGlobal.PublishTo(st.proc.OutChannelId, message)
}

func (st *State) NewLine() {
	keyEnter := msg.MessageKey{
		Key:    msg.KeyEnter,
		Scan:   0,
		Action: uint8(msg.Action(msg.Press)),
		Mod:    0}

	st.publishToOut(msg.Serialize(msg.TypeKey, keyEnter))
}

func (st *State) PrintLn(s string) {
	for _, c := range s {
		st.sendChar(uint32(c))
	}

	st.NewLine()
}

func (st *State) PrintError(s string) {
	st.PrintLn("**** ERROR! ****    " + s)
}

func (st *State) Printf(format string, vars ...interface{}) {
	formattedString := fmt.Sprintf(format, vars)
	for _, c := range formattedString {
		st.sendChar(uint32(c))
	}
}

func (st *State) sendChar(c uint32) {
	var s string

	switch c {
	case msg.EscNewLine:
		st.NewLine()
		return
	case msg.EscTab:
		s = "Tab"
	case msg.EscCarriageReturn:
		s = "Carriage Return"
	case msg.EscBackSpace:
		s = "BackSpace"
		// case msg.EscBackSlash:
		// 	s = "BackSlash"
	}

	if s != "" {
		println("TASK ENCOUNTERED ESCAPE CHAR FOR [" + s + "], NOT SENDING TO TERMINAL")
		return
	}

	m := msg.Serialize(msg.TypePutChar, msg.MessagePutChar{0, c})
	st.publishToOut(m) // EVERY publish action prefixes another chan id
}
