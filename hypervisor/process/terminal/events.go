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

	//if we have one more free space to put character into
	if len(commands[currCmd]) < maxCommandSize {
		commands[currCmd] = commands[currCmd][:cursPos] + string(m.Char) + commands[currCmd][cursPos:]
		cursorForward()
		EchoWholeCommand(self.proc.OutChannelId)
	}
}

func (self *State) onKey(m msg.MessageKey, serializedMsg []byte) {
	println("process/terminal/events.onKey()")

	switch msg.Action(m.Action) {

	case msg.Press:
		fallthrough
	case msg.Repeat:
		switch m.Key {

		case msg.KeyHome:
			cursPos = len(prompt)
		case msg.KeyEnd:
			cursPos = len(commands[currCmd])

		case msg.KeyUp:
			if msg.ModifierKey(m.Mod) == msg.ModControl {
				currCmd = 0
			} else {
				traverseCommands(-1)
			}
		case msg.KeyDown:
			if msg.ModifierKey(m.Mod) == msg.ModControl {
				currCmd = len(commands) - 1 // this could crash if we don't make sure at least 1 command always exists
			} else {
				traverseCommands(+1)
			}
		case msg.KeyLeft:
			cursorBackward()
		case msg.KeyRight:
			cursorForward()

		case msg.KeyEnter:
			hypervisor.DbusGlobal.PublishTo(
				dbus.ChannelId(self.proc.OutChannelId), serializedMsg)

			log = append(log, commands[currCmd])
			commands = append(commands, prompt)
			traverseCommands(+1)

		case msg.KeyBackspace:
			if cursorBackward() { //...succeeded
				commands[currCmd] = commands[currCmd][:cursPos] + commands[currCmd][cursPos+1:]
			}
		case msg.KeyDelete:
			if cursPos < len(commands[currCmd]) {
				commands[currCmd] = commands[currCmd][:cursPos] + commands[currCmd][cursPos+1:]
			}
		}

		EchoWholeCommand(self.proc.OutChannelId)
	case msg.Release:
		// most keys will do nothing upon release
	}
}

func cursorForward() {
	cursPos++

	if cursPos > len(commands[currCmd]) {
		cursPos = len(commands[currCmd])
	} else if cursPos > maxCommandSize { //allows cursor to be one position beyond last char
		cursPos = maxCommandSize
	}
}

func cursorBackward() bool { //returns whether moved successfully
	cursPos--

	if cursPos < len(prompt) {
		cursPos = len(prompt)
		return false
	}

	return true
}
