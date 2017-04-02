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

		case msg.KeyC:
			// TODO: I guess, Doesn't work with new Implementation ! ! !
			if m.Mod == msg.GLFW_MOD_CONTROL {
				// st.PrintError("Ctrl+C pressed")
				// err := st.proc.ExitExtProcess()
				// if err != nil {
				// 	st.PrintError(err.Error())
				// }
			}

		case msg.KeyZ:
			// TODO: I guess, Doesn't work with new Implementation ! ! !
			// if m.Mod == msg.GLFW_MOD_CONTROL {
			// 	// st.PrintError("Ctrl+Z pressed")
			// 	err := st.proc.SendAttachedToBg()
			// 	if err != nil {
			// 		st.PrintError(err.Error())
			// 	}
			// 	st.PrintLn("Attached process sent to background.")
			// }

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
	if len(cmd) < 1 {
		println("**** ERROR! ****   Command was empty!  Returning.")
		return
	}

	s := "actOnCommand() called on"
	println(s, "cmd:", cmd)

	for _, arg := range args {
		println("with arg:", arg)
	}

	if st.proc.HasExtProcessAttached() {
		// Redirect input to the attached process
		extProc, err := st.proc.GetAttachedExtProcess()

		if err != nil {
			println(err.Error())
			return
		}

		extProc.ProcessIn <- []byte(st.Cli.CurrentCommandLine())
	} else { //internal task
		switch cmd {

		//helps
		case "?":
			fallthrough
		case "h":
			fallthrough
		case "help":
			st.PrintLn("Help will now work after you EXEC something, I presume...")

		case "j":
			// Doesn't work yet with new implementation ! ! !
			// println("Inside the jobs")
			// for id, extProc := range st.proc.extProcesses {
			// 	baseCommand := strings.Split(extProc.CommandLine, " ")[0]
			// 	st.PrintLn("[ " + strconv.Itoa(int(id)) + " ] -> [ " + baseCommand + " ]")
			// }
		case "jobs":
			// Doesn't work yet with new implementation ! ! !
			// println("Inside the jobs")
			// for id, extProc := range st.proc.extProcesses {
			// 	st.PrintLn("[ " + strconv.Itoa(int(id)) + " ] -> [ " + extProc.CommandLine + " ]")
			// }

		case "fg":
			// Doesn't work yet with new implementation ! ! !
			// println("Inside the FG")
			// // FIXME: Buggy tomorrow I'll continue to work on this
			// if len(args) < 1 {
			// 	st.PrintError("Must pass the job id! eg: fg 1")
			// } else {
			// 	extProcID, err := strconv.Atoi(args[0])
			// 	if err != nil {
			// 		st.PrintError(err.Error())
			// 	}

			// 	err = st.proc.SendExtToFg(msg.ExtProcessId(extProcID))
			// 	if err != nil {
			// 		st.PrintError(err.Error())
			// 	}
			// }

		case "rpc":
			// Doesn't work yet with new implementation ! ! !
			// tokens := []string{"go", "run", "rpc/cli/cli.go"}
			// err := st.proc.AddAttachStart(tokens)
			// if err != nil {
			// 	st.PrintLn(err.Error())
			// }

		case "e":
			fallthrough
		case "ex":
			fallthrough
		case "exec":
			// Doesn't work yet with new implementation ! ! !
			st.PrintLn("Doesn't work yet with new implementation ! ! ! Please wait...")
			// if len(args) < 1 {
			// 	st.PrintError("Must pass a command into EXEC!")
			// } else { //execute
			// 	err := st.proc.AddAttachStart(args) //(includes a command)
			// 	if err != nil {
			// 		for i := 0; i < 5; i++ {
			// 			println("********* " + err.Error() + " *********")
			// 		}
			// 		st.PrintLn(err.Error())
			// 	}
			// }

		default:
			st.PrintError("\"" + cmd + "\" is an unknown command.")

		}
	}
}
