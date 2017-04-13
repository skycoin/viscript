package process

import (
	"strconv"

	"strings"

	"github.com/corpusc/viscript/hypervisor"
	extTask "github.com/corpusc/viscript/hypervisor/task_ext"
	"github.com/corpusc/viscript/msg"
)

func (st *State) onChar(m msg.MessageChar) {
	//println("process/terminal/events.onChar()")
	st.Cli.InsertCharIfItFits(m.Char, st)
}

func (st *State) onKey(m msg.MessageKey, serializedMsg []byte) {
	switch msg.Action(m.Action) {

	case msg.Press: //one time, when key is first pressed
		//modifier key combos should probably never auto-repeat
		st.actOnOneTimeHotkeys(m)
		fallthrough
	case msg.Repeat: //constantly repeated for as long as key is pressed
		st.actOnRepeatableKeys(m, serializedMsg)
		st.Cli.EchoWholeCommand(st.proc.OutChannelId)

	case msg.Release:
		//most keys will do nothing upon release
	}
}

func (st *State) actOnOneTimeHotkeys(m msg.MessageKey) {
	switch m.Key {

	case msg.KeyC:
		if m.Mod == msg.GLFW_MOD_CONTROL {
			//FIXME: doesn't work yet, have to look into
			//implementation of extTask should change to Tick() like structure
			//then exiting the task will be different
			//current implementation i.e setting flag for the for statement in goroutine
			//doesn't work.

			//st.PrintError("Ctrl+C pressed")
			//err := st.proc.ExitExtProcess()
			//if err != nil {
			//	st.PrintError(err.Error())
			//}
		}

	case msg.KeyZ:
		if m.Mod == msg.GLFW_MOD_CONTROL {
			// st.PrintError("Ctrl+Z pressed")
			// err := st.proc.SendAttachedToBg()
			// if err != nil {
			// 	st.PrintError(err.Error())
			// }
			// st.PrintLn("Attached process sent to background.")
			// st.proc.DetachExternalProcess()
		}

	}
}

func (st *State) actOnRepeatableKeys(m msg.MessageKey, serializedMsg []byte) {
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
}

func (st *State) onMouseScroll(m msg.MessageMouseScroll, serializedMsg []byte) {
	if m.HoldingControl {
		hypervisor.DbusGlobal.PublishTo(st.proc.OutChannelId, serializedMsg)
	} else {
		st.Cli.AdjustBackscrollOffset(int(m.Y))
	}
}

func (st *State) actOnCommand() {
	cmd, args := st.Cli.CurrentCommandAndArgs()
	if len(cmd) < 1 {
		println("**** ERROR! ****   Command was empty!  Returning.")
		return
	}

	s := "actOnCommand() called with"
	println(s, "cmd:", cmd)

	for _, arg := range args {
		println(s, "arg:", arg)
	}

	if st.proc.HasExtProcessAttached() {
		extProcInChannel := st.proc.attachedExtProcess.GetProcessInChannel()
		extProcInChannel <- []byte(st.Cli.CurrentCommandLine())
		return
	}

	//internal task
	switch cmd {
	case "?":
		fallthrough
	case "h":
		fallthrough
	case "help":
		st.PrintLn("Current commands:")
		st.PrintLn("    help:               	This message ('?' or 'h' for short).")
		st.PrintLn("    start (-a) <command>: 	Start external task. (pass -a before command to aslo attach)")
		st.PrintLn("    attach <id>:       		Attach external process with given id to terminal.")
		st.PrintLn("    ls (-f):           	    List external processes (pass -f for full commands).")
		//st.PrintLn("    new_terminal:     Add new terminal.")
		//st.PrintLn("    rpc:              Issues command: \"go run rpc/cli/cli.go\"")
		//st.PrintLn("Current hotkeys:")
		//st.PrintLn("    CTRL+C:           ___description goes here___")
		//st.PrintLn("    CTRL+Z:           ___description goes here___")

	//list external process with ls
	case "ls":
		extTaskMap := hypervisor.ExtProcessListGlobal.ProcessMap
		if len(extTaskMap) == 0 {
			st.PrintLn("No external processes running.\n" +
				"Try starting one with \"start\" command (\"help\" or \"h\" for help).")
			break
		}

		fullPrint := false

		if len(args) > 0 && args[0] == "-f" {
			fullPrint = true
		}

		for procId, extProc := range extTaskMap {

			procCommand := ""

			if fullPrint {
				procCommand = extProc.GetFullCommandLine()
			} else {
				procCommand = strings.Split(
					extProc.GetFullCommandLine(), " ")[0]
			}

			st.Printf("[ %d ] -> [ %s ]\n", int(procId), procCommand)
		}

	//attach external task to terminal task
	case "attach":
		if len(args) < 1 {
			st.PrintError("No task id passed! eg: attach 1")
			break
		}

		passedId, err := strconv.Atoi(args[0])
		if err != nil {
			st.PrintError("Task id must be an integer.")
			break
		}

		extProcId := msg.ExtProcessId(passedId)

		extProc, err := hypervisor.GetExtProcess(extProcId)
		if err != nil {
			st.PrintError(err.Error())
			break
		}

		st.PrintLn(extProc.GetFullCommandLine())
		st.proc.AttachExternalProcess(extProc)

	case "r":
		fallthrough
	case "rpc":
		//Doesn't work yet with new implementation ! ! !
		// tokens := []string{"go", "run", "rpc/cli/cli.go"}
		// err := st.proc.AddAttachStart(tokens)
		// if err != nil {
		// 	st.PrintError(err.Error())
		// }

	//start new external task, detached running in bg by default
	case "s":
		fallthrough
	case "start":
		if len(args) < 1 {
			st.PrintError("Must pass a command into Start!")
		} else {
			detached := args[0] != "-a"

			if !detached {
				args = args[1:]
			}

			newExtProc, err := extTask.MakeNewTaskExternal(args, detached)
			if err != nil {
				st.PrintError(err.Error())
				break
			}

			err = newExtProc.Start()
			if err != nil {
				st.PrintError(err.Error())
				break
			}

			extProcInterface := newExtProc.GetExtProcessInterface()

			procId := hypervisor.AddExtProcess(extProcInterface)

			if !detached {
				st.proc.AttachExternalProcess(extProcInterface)
			}

			st.PrintLn("Added External Process (ID: " +
				strconv.Itoa(int(procId)) + ", Command: " +
				newExtProc.CommandLine + ")")
		}

	//add new terminal
	case "n":
		fallthrough
	case "new":
		fallthrough
	case "new_terminal":
		// viewport.Terms.Add()

	default:
		st.PrintError("\"" + cmd + "\" is an unknown command.")

	}
}
