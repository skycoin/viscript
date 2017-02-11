package process

import (
	"github.com/corpusc/viscript/hypervisor"
	"github.com/corpusc/viscript/hypervisor/dbus"
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
		self.onKey(m)

	case msg.TypeFrameBufferSize:
		// FIXME: BRAD SAYS THIS IS NOT INPUT
		var m msg.MessageFrameBufferSize
		msg.MustDeserialize(message, &m)
		self.onFrameBufferSize(m)

	default:
		println("UNKNOWN MESSAGE TYPE!")
	}

	if self.DebugPrintInputEvents {
		println()
	}

	return message
}

//
//EVENT HANDLERS
//

func (self *State) onFrameBufferSize(m msg.MessageFrameBufferSize) {
	println("process/terminal/events.onFrameBufferSize()")
}

func (self *State) onChar(m msg.MessageChar) {
	println("process/terminal/events.onChar()")

	//FIXME? actual terminal id really needed?
	message := msg.Serialize(
		msg.TypePutChar, msg.MessagePutChar{0, m.Char}) //....just gave it 0 for now
	hypervisor.DbusGlobal.PublishTo(
		dbus.ChannelId(self.proc.OutChannelId), message) //EVERY publish action prefixes another chan id
}

func (self *State) onKey(m msg.MessageKey) {
	println("process/terminal/events.onKey()")
}
