package process

import (
	"github.com/corpusc/viscript/hypervisor"
	"github.com/corpusc/viscript/hypervisor/dbus"
	"github.com/corpusc/viscript/msg"
)

// func (self *State) onFrameBufferSize(m msg.MessageFrameBufferSize) {
// 	println("process/terminal/events.onFrameBufferSize()")
// 	message := msg.Serialize(
// 		msg.TypeFrameBufferSize, msg.MessageFrameBufferSize{m.X, m.Y})
// 	hypervisor.DbusGlobal.PublishTo(
// 		dbus.ChannelId(self.proc.OutChannelId), message)
// }

func (self *State) onChar(m msg.MessageChar) {
	//println("process/terminal/events.onChar()")

	if len(commands[currCmd]) < maxCommandSize {
		// (we have free space to put character into)
		commands[currCmd] = commands[currCmd][:cursPos] + string(m.Char) + commands[currCmd][cursPos:]
		moveOneStepRight()
		EchoWholeCommand(self.proc.OutChannelId)
	}
}

func (self *State) onKey(m msg.MessageKey, serializedMsg []byte) {
	//println("process/terminal/events.onKey()")

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
			goUpCommandHistory(m.Mod)
		case msg.KeyDown:
			goDownCommandHistory(m.Mod)

		case msg.KeyLeft:
			moveOrJumpCursorLeft(m.Mod)
		case msg.KeyRight:
			moveOrJumpCursorRight(m.Mod)

		case msg.KeyBackspace:
			if moveOneStepLeft() { //...succeeded
				commands[currCmd] = commands[currCmd][:cursPos] + commands[currCmd][cursPos+1:]
			}
		case msg.KeyDelete:
			if cursPos < len(commands[currCmd]) {
				commands[currCmd] = commands[currCmd][:cursPos] + commands[currCmd][cursPos+1:]
			}

		case msg.KeyEnter:
			self.actOnEnter(serializedMsg)
		}

		EchoWholeCommand(self.proc.OutChannelId)
	case msg.Release:
		// most keys will do nothing upon release
	}
}

func moveOneStepLeft() bool { //returns whether moved successfully
	cursPos--

	if cursPos < len(prompt) {
		cursPos = len(prompt)
		return false
	}

	return true
}

func moveOneStepRight() bool { //returns whether moved successfully
	cursPos++

	if cursPos > len(commands[currCmd]) {
		cursPos = len(commands[currCmd])
		return false
	} else if cursPos > maxCommandSize { //allows cursor to be one position beyond last char
		cursPos = maxCommandSize
		return false
	}

	return true
}

func moveOrJumpCursorLeft(mod uint8) {
	if msg.ModifierKey(mod) == msg.ModControl {
		numSpaces := 0
		numVisible := 0 //NON-space

		for moveOneStepLeft() == true {
			if commands[currCmd][cursPos] == ' ' {
				numSpaces++

				if numVisible > 0 || numSpaces > 1 {
					moveOneStepRight()
					break
				}
			} else {
				numVisible++
			}
		}
	} else {
		moveOneStepLeft()
	}
}

func moveOrJumpCursorRight(mod uint8) {
	if msg.ModifierKey(mod) == msg.ModControl {
		for moveOneStepRight() == true {
			if cursPos < len(commands[currCmd]) &&
				commands[currCmd][cursPos] == ' ' {
				moveOneStepRight()
				break
			}
		}
	} else {
		moveOneStepRight()
	}
}

func goUpCommandHistory(mod uint8) {
	if msg.ModifierKey(mod) == msg.ModControl {
		currCmd = 0
	} else {
		traverseCommands(-1)
	}
}

func goDownCommandHistory(mod uint8) {
	if msg.ModifierKey(mod) == msg.ModControl {
		currCmd = len(commands) - 1 // this could cause crash if we don't make sure at least 1 command always exists
	} else {
		traverseCommands(+1)
	}
}

func (self *State) actOnEnter(serializedMsg []byte) {
	numLineFeeds := 1

	if cursPos >= 64 { //FIXME using Terminal's self.GridSize.X
		numLineFeeds++
	}

	for numLineFeeds > 0 {
		numLineFeeds--
		hypervisor.DbusGlobal.PublishTo(
			dbus.ChannelId(self.proc.OutChannelId), serializedMsg)
	}

	log = append(log, commands[currCmd])
	commands = append(commands, prompt)
	currCmd = len(commands) - 1
	cursPos = len(commands[currCmd])
}
