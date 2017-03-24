package process

import (
	"strings"

	"github.com/corpusc/viscript/hypervisor"
	"github.com/corpusc/viscript/msg"
)

func (st *State) onChar(m msg.MessageChar) {
	//println("process/terminal/events.onChar()")
	if st.Cli.HasEnoughSpace() {
		st.Cli.InsertCharAtCursor(m.Char)
		st.Cli.EchoWholeCommand(st.proc.OutChannelId)
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
			st.Cli.CursPos = len(st.Cli.Prompt)
		case msg.KeyEnd:
			st.Cli.CursPos = len(st.Cli.Commands[st.Cli.CurrCmd])

		case msg.KeyUp:
			st.Cli.goUpCommandHistory(m.Mod)
		case msg.KeyDown:
			st.Cli.goDownCommandHistory(m.Mod)

		case msg.KeyLeft:
			st.Cli.moveOrJumpCursorLeft(m.Mod)
		case msg.KeyRight:
			st.Cli.moveOrJumpCursorRight(m.Mod)

		case msg.KeyBackspace:
			if st.Cli.moveCursorOneStepLeft() { //...succeeded
				st.Cli.DeleteCharAtCursor()
			}
		case msg.KeyDelete:
			if st.Cli.CursPos < len(st.Cli.Commands[st.Cli.CurrCmd]) {
				st.Cli.DeleteCharAtCursor()
			}

		case msg.KeyEnter:
			st.Cli.OnEnter(st, serializedMsg)
		}

		st.Cli.EchoWholeCommand(st.proc.OutChannelId)
	case msg.Release:
		// most keys will do nothing upon release
	}
}

func (st *State) onMouseScroll(m msg.MessageMouseScroll, serializedMsg []byte) {
	//println("process/terminal/events.onMouseScroll()")
	hypervisor.DbusGlobal.PublishTo(st.proc.OutChannelId, serializedMsg)
}

func (st *State) actOnCommand() {
	command, args := st.Cli.GetCommandWithArgs()

	var wholeCommand []string
	wholeCommand = append(wholeCommand, strings.ToLower(command))

	for _, v := range args {
		wholeCommand = append(wholeCommand, strings.ToLower(v))
	}

	if len(strings.ToLower(command)) > 0 {
		st.CmdOut <- []byte(strings.Join(wholeCommand, " "))
	}

	// switch strings.ToLower(command) {

	// case "?":
	// 	fallthrough
	// case "h":
	// 	fallthrough
	// case "help":
	// 	st.PrintLn("Yes master, help is coming 'very soon'. (TM)")
	// default:
	// 	st.PrintLn("ERROR: \"" + command + "\" is an unknown command.")
	// }
}
