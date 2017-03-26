package process

import (
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
	cmd, args := st.Cli.CurrentCommandAndArgs()
	println("actOnCommand() called with:", cmd)

	for _, s := range args {
		println("and arg:", s)
	}

	if len(cmd) < 1 {
		println("Command was empty, returning.")
		return
	}

	if st.proc.HasExtProcessAttached() {
		// Redirect input to the attached process
		extProc, err := st.proc.GetAttachedExtProcess()

		if err != nil {
			println(err.Error())
			return
		}

		extProc.CmdOut <- []byte(st.Cli.CurrentCommandLine())
	} else { //internal task
		switch cmd {

		case "?":
			fallthrough
		case "h":
			fallthrough
		case "help":
			st.PrintLn("Help will now work after you EXEC something, I presume...")

		case "exec":
			// extCommand := strings.Join(args, " ") //FIXME?  shouldn't this
			//take a list of args rather than joining them into a string?  they'll need to be
			//.Split() later right?  and from here on,
			//they are called cmd/command/etc. rather than args...which doesn't seem right?
			//shouldn't need to know it was "exec" that triggered this?

			err := st.proc.AddAndAttach(args)
			if err != nil {
				println(err.Error())
			}

		default:
			st.PrintLn("ERROR: \"" + cmd + "\" is an unknown command.")

		}
	}
}
