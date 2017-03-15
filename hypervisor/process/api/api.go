package api

import (
	"fmt"

	"github.com/corpusc/viscript/hypervisor"
	"github.com/corpusc/viscript/msg"
)

func (st *State) publishToOut(message []byte) {
	hypervisor.DbusGlobal.PublishTo(st.Proc.OutChannelId, message)
}

func (st *State) NewLine() {
	keyEnter := msg.MessageKey{
		Key:    msg.KeyEnter,
		Scan:   0,
		Action: uint8(msg.Action(msg.Press)),
		Mod:    0}

	m := msg.Serialize(msg.TypeKey, keyEnter)
	st.publishToOut(m)
}

func (st *State) PrintLn(s string) {
	for _, c := range s {
		if uint32(c) == 10 { // newline
			st.NewLine()
			continue
		}
		st.SendChar(uint32(c))
	}

	st.NewLine()
}

func (st *State) Printf(format string, vars ...interface{}) {
	formattedString := fmt.Sprintf(format, vars)
	for _, c := range formattedString {
		if uint32(c) == 10 { // newline
			st.NewLine()
			continue
		}
		st.SendChar(uint32(c))
	}
}

func (st *State) SendChar(c uint32) {
	m := msg.Serialize(msg.TypePutChar, msg.MessagePutChar{0, c})
	st.publishToOut(m)
}
