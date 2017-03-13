package externalprocess

import (
	"github.com/corpusc/viscript/msg"
)

func (st *State) UnpackEvent(msgType uint16, message []byte) []byte {
	println("process/terminal/events.UnpackEvent()")

	switch msgType {

	case msg.TypeChar:
		var m msg.MessageChar
		msg.MustDeserialize(message, &m)
		st.onChar(m)

	case msg.TypeKey:
		var m msg.MessageKey
		msg.MustDeserialize(message, &m)
		st.onKey(m, message)

	// case msg.TypeFrameBufferSize:
	// 	// FIXME: BRAD SAYS THIS IS NOT INPUT
	// 	var m msg.MessageFrameBufferSize
	// 	msg.MustDeserialize(message, &m)
	// 	st.onFrameBufferSize(m)

	default:
		println("UNKNOWN MESSAGE TYPE!")
	}

	if st.DebugPrintInputEvents {
		println()
	}

	return message
}
