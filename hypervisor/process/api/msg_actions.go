package api

import (
	"strings"

	"github.com/corpusc/viscript/hypervisor"
	"github.com/corpusc/viscript/msg"
)

func (st *State) onChar(m msg.MessageChar) {
	println("process/terminal/events.onChar()")

	if len(commands[currCmd]) < maxCommandSize {
		// (we have free space to put character into)
		commands[currCmd] = commands[currCmd][:cursPos] + string(m.Char) + commands[currCmd][cursPos:]
		moveOneStepRight()
		EchoWholeCommand(st.Proc.OutChannelId)
	}
}

func (st *State) onKey(m msg.MessageKey, serializedMsg []byte) {
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
			st.actOnEnter(serializedMsg)
		}

		EchoWholeCommand(st.Proc.OutChannelId)
	case msg.Release:
		// most keys will do nothing upon release
	}
}

func (st *State) onMouseScroll(m msg.MessageMouseScroll, serializedMsg []byte) {
	//println("process/terminal/events.onMouseScroll()")
	hypervisor.DbusGlobal.PublishTo(st.Proc.OutChannelId, serializedMsg)
}

func (st *State) actOnCommand(command []string) {

	switch strings.ToLower(command[0]) {

	case "?":
		fallthrough
	case "h":
		fallthrough
	case "help":
		st.PrintLn("Yes master, help is coming 'very soon'. (TM)")
	default:
		st.PrintLn("ERROR: \"" + command[0] + "\" is an unknown command.")
	}
}
