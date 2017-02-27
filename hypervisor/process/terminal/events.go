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

//
//EVENT HANDLERS
//

// func (self *State) onFrameBufferSize(m msg.MessageFrameBufferSize) {
// 	println("process/terminal/events.onFrameBufferSize()")
// 	message := msg.Serialize(
// 		msg.TypeFrameBufferSize, msg.MessageFrameBufferSize{m.X, m.Y})
// 	hypervisor.DbusGlobal.PublishTo(
// 		dbus.ChannelId(self.proc.OutChannelId), message)
// }

func (self *State) onChar(m msg.MessageChar) {
	println("process/terminal/events.onChar()")
	commands[currCmd] += string(m.Char)
	cursPos++
	EchoWholeCommand(self.proc.OutChannelId)
}

func (self *State) onKey(m msg.MessageKey, serializedMsg []byte) {
	println("process/terminal/events.onKey()")

	switch msg.Action(m.Action) {

	case msg.Press:
		fallthrough
	case msg.Repeat:
		switch m.Key {

		case msg.KeyUp:
			traverseCommands(-1)
			EchoWholeCommand(self.proc.OutChannelId)

		case msg.KeyDown:
			traverseCommands(+1)
			EchoWholeCommand(self.proc.OutChannelId)

		case msg.KeyLeft:
			cursPos--

			if cursPos < 0 {
				cursPos = 0
			}

			EchoWholeCommand(self.proc.OutChannelId)

		case msg.KeyRight:
			cursPos++

			if cursPos >= maxCommandSize {
				cursPos = maxCommandSize - 1
			}

			EchoWholeCommand(self.proc.OutChannelId)

		case msg.KeyEnter:
			hypervisor.DbusGlobal.PublishTo(
				dbus.ChannelId(self.proc.OutChannelId), serializedMsg)
			log = append(log, commands[currCmd])
			commands = append(commands, prompt)
			traverseCommands(+1)

		case msg.KeyBackspace:
			hypervisor.DbusGlobal.PublishTo(
				dbus.ChannelId(self.proc.OutChannelId), serializedMsg)
		}
	case msg.Release:
		// most keys will do nothing upon release
	}
}
