package process

import (
	"fmt"
	"strings"

	"github.com/corpusc/viscript/hypervisor"
	"github.com/corpusc/viscript/msg"
)

func (st *State) onChar(m msg.MessageChar) {
	//println("process/terminal/events.onChar()")

	if len(commands[currCmd]) < maxCommandSize {
		// (we have free space to put character into)
		commands[currCmd] = commands[currCmd][:cursPos] + string(m.Char) + commands[currCmd][cursPos:]
		moveOneStepRight()
		EchoWholeCommand(st.proc.OutChannelId)
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

		EchoWholeCommand(st.proc.OutChannelId)
	case msg.Release:
		// most keys will do nothing upon release
	}
}

func (st *State) onMouseScroll(m msg.MessageMouseScroll, serializedMsg []byte) {
	//println("process/terminal/events.onMouseScroll()")
	hypervisor.DbusGlobal.PublishTo(st.proc.OutChannelId, serializedMsg)
}

func (st *State) actOnCommand() {
	words := strings.Split(commands[currCmd][len(prompt):], " ")

	switch strings.ToLower(words[0]) {

	case "?":
		fallthrough
	case "h":
		fallthrough
	case "help":
		st.printLn("Yes master, help is coming 'very soon'. (TM)")
	default:
		st.printLn("ERROR: \"" + words[0] + "\" is an unknown command.")
	}
}

func (st *State) newLine() {
	m := msg.Serialize(
		msg.TypeKey,
		msg.MessageKey{
			msg.KeyEnter,
			0, // Scan   uint32
			uint8(msg.Action(msg.Press)),
			0}) // Mod
	hypervisor.DbusGlobal.PublishTo(st.proc.OutChannelId, m)
}

func (st *State) printLn(s string) {
	for _, c := range s {
		st.sendChar(uint32(c))
	}

	st.newLine()
}

func (st *State) printf(format string, vars ...interface{}) {
	formattedString := fmt.Sprintf(format, vars)
	for _, c := range formattedString {
		st.sendChar(uint32(c))
	}
}

func (st *State) sendChar(c uint32) {
	m := msg.Serialize(msg.TypePutChar, msg.MessagePutChar{0, c})
	hypervisor.DbusGlobal.PublishTo(st.proc.OutChannelId, m) // EVERY publish action prefixes another chan id
}
