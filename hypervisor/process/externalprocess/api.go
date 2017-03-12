package externalprocess

import (
	"fmt"

	"github.com/corpusc/viscript/hypervisor"
	"github.com/corpusc/viscript/msg"
)

func (st *State) publishToOut(message []byte) {
	hypervisor.DbusGlobal.PublishTo(st.proc.OutChannelId, message)
}

func (st *State) newLine() {
	m := msg.Serialize(
		msg.TypeKey,
		msg.MessageKey{
			msg.KeyEnter,
			0, // Scan   uint32
			uint8(msg.Action(msg.Press)),
			0}) // Mod
	// hypervisor.DbusGlobal.PublishTo(st.proc.OutChannelId, m)
	st.publishToOut(m)
}

func (st *State) printLn(s string) {
	for _, c := range s {
		if uint32(c) == 10 { // newline
			st.newLine()
			continue
		}
		st.sendChar(uint32(c))
	}

	st.newLine()
}

func (st *State) printf(format string, vars ...interface{}) {
	formattedString := fmt.Sprintf(format, vars)
	for _, c := range formattedString {
		if uint32(c) == 10 { // newline
			st.newLine()
			continue
		}
		st.sendChar(uint32(c))
	}
}

func (st *State) sendChar(c uint32) {
	m := msg.Serialize(msg.TypePutChar, msg.MessagePutChar{0, c})
	// hypervisor.DbusGlobal.PublishTo(st.proc.OutChannelId, m) // EVERY publish action prefixes another chan id
	st.publishToOut(m)
}
