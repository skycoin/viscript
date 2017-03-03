package process

import (
	"github.com/corpusc/viscript/msg"
)

func (self *State) UnpackEvent(msgType uint16, message []byte) []byte {
	println("process/terminal/events.UnpackEvent()")

	switch msgType {

	case msg.TypeChar:
		var m msg.MessageChar
		msg.MustDeserialize(message, &m)
		self.onChar(m)

	case msg.TypeKey:
		var m msg.MessageKey
		msg.MustDeserialize(message, &m)
		self.onKey(m, message)

	// case msg.TypeFrameBufferSize:
	// 	// FIXME: BRAD SAYS THIS IS NOT INPUT
	// 	var m msg.MessageFrameBufferSize
	// 	msg.MustDeserialize(message, &m)
	// 	self.onFrameBufferSize(m)

	default:
		println("UNKNOWN MESSAGE TYPE!")
	}

	if self.DebugPrintInputEvents {
		println()
	}

	return message
}
