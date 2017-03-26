package process

import (
	"strings"

	"github.com/corpusc/viscript/hypervisor"
	"github.com/corpusc/viscript/msg"
)

func (st *State) onChar(m msg.MessageChar) {
	//println("process/terminal/events.onChar()")
	st.Cli.InsertCharIfItFits(m.Char, st)
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
	println("actOnCommand() called with:", command, "and args:", args)
	command = strings.ToLower(command)

	if len(command) == 0 {
		println("Command was empty, returning.")
		return
	}

	var wholeCommand []string
	wholeCommand = append(wholeCommand, command)

	for _, v := range args {
		wholeCommand = append(wholeCommand, strings.ToLower(v))
	}

	if st.proc.HasExtProcessAttached() {
		// Redirect input to the attached process
		extProc, err := st.proc.GetAttachedExtProcess()

		if err != nil {
			println(err.Error())
			return
		}

		extProc.CmdOut <- []byte(strings.Join(wholeCommand, " "))
	} else {
		// Handle terminal command here
		switch command {

		case "?":
			fallthrough
		case "h":
			fallthrough
		case "help":
			st.PrintLn("Yes master, help is coming 'very soon'. (TM)")
		case "exec":
			extCommand := strings.Join(wholeCommand[1:], " ")
			err := st.proc.AddAndAttach(extCommand)
			if err != nil {
				println(err.Error())
			}
		default:
			st.PrintLn("ERROR: \"" + command + "\" is an unknown command.")

		}
	}
}
