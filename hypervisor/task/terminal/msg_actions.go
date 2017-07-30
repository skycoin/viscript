package task

import (
	"fmt"

	//"github.com/skycoin/viscript/app"
	"github.com/skycoin/viscript/hypervisor"
	"github.com/skycoin/viscript/msg"
)

func (st *State) onMouseScroll(m msg.MessageMouseScroll, serializedMsg []byte) {
	if m.HoldingControl {
		//send message to 'Terminal'
		//(i believe this was only used for sidescrolling in the text editor,
		//but we might want to use this as an alternative to PGUP/PGDN keys
		//when focused on a CLI terminal)
		hypervisor.DbusGlobal.PublishTo(st.task.OutChannelId, serializedMsg)
	} else {
		st.Cli.AdjustBackscrollOffset(int(m.Y))
		//just pass unchanged VisualInfo
		//(which is acted upon like a boolean flag)
		st.makePageOfLog(st.VisualInfo)
	}
}

func (st *State) onChar(m msg.MessageChar) {
	//println("task/terminal/msg_actions.onChar()")
	st.Cli.InsertCharIfItFits(m.Char, st)
}

func (st *State) onKey(m msg.MessageKey, serializedMsg []byte) {
	switch msg.Action(m.Action) {

	case msg.Press: //one time, when key is first pressed
		//modifier key combos should never auto-repeat?
		st.onNONRepeatableKey(m)
		fallthrough
		//^^^^^^^^^^
	case msg.Repeat: //constantly repeated for as long as key is pressed
		//THIS IS ALSO DONE ON 'PRESS' EVENTS, via FALLTHROUGH!
		st.onRepeatableKey(m, serializedMsg)
		st.Cli.EchoWholeCommand(st.task.OutChannelId)

	case msg.Release:
		//most keys will probably never do anything upon release
	}
}

func (st *State) onTerminalIds(m msg.MessageTerminalIds) {
	st.storedTerminalIds = m.TermIds
	num := len(m.TermIds)
	str := fmt.Sprintf("%d Terminal", num)

	if num == 1 {
		str += ":"
	} else {
		str += "s:"
	}

	st.PrintLn(str)

	for i, termID := range m.TermIds {
		s := fmt.Sprintf("    [%d] %d", i, termID)

		if termID == m.Focused {
			st.PrintLn(s + "    (FOCUSED)")
		} else {
			st.PrintLn(s)
		}
	}

	st.Cli.EchoWholeCommand(st.task.OutChannelId)
}

func (st *State) onNONRepeatableKey(m msg.MessageKey) {
	switch m.Key {

	case msg.KeyC:
		if m.Mod == msg.GLFW_MOD_CONTROL {
			if !st.task.HasExtTaskAttached() {
				return
			}

			//TODO: send the exit request i.e sigint or
			//just maybe kill the task?
		}

	case msg.KeyZ:
		if m.Mod == msg.GLFW_MOD_CONTROL {
			if !st.task.HasExtTaskAttached() {
				return
			}

			st.PrintLn("Detaching External Task")
			st.task.DetachExternalTask()
		}

	}
}

func (st *State) onRepeatableKey(m msg.MessageKey, serializedMsg []byte) {
	//ALSO EXECUTED ON PRESS EVENTS via case fallthrough!

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

func (st *State) onUserCommand() {
	cmd, args := st.Cli.CurrentCommandAndArgsInLowerCase()

	if len(cmd) < 1 {
		println("**** ERROR! ****   Command was empty!  Returning.")
		return
	}

	s := "onUserCommand()      command: \"" + cmd + "\""

	for i, arg := range args {
		if i == 0 {
			s += "      args:"
		}

		s += "   [" +
			fmt.Sprintf("%d", i) +
			"] \"" + arg + "\""
	}

	println(s)

	if st.task.HasExtTaskAttached() {
		extTaskInChannel := st.task.attachedExtTask.GetTaskInChannel()
		extTaskInChannel <- []byte(st.Cli.CurrentCommandLineInLowerCase())
		return
	}

	//internal task handling
	switch cmd {

	case "?":
		fallthrough
	case "h":
		fallthrough
	case "help":
		if len(args) != 0 {
			st.commandAppHelp(args)
			break
		} else {
			st.commandHelp()
		}

	//display all apps with descriptions
	case "apps":
		st.commandApps()

	//attach external task to terminal task
	case "attach":
		st.commandAttach(args)

	case "c":
		fallthrough
	case "cls":
		fallthrough
	case "clear":
		st.commandClearTerminal()

	//delete terminal with given index
	case "xt":
		fallthrough
	case "ct":
		fallthrough
	case "close_term":
		st.commandCloseTerminalFirstStage(args)

	case "list_tasks":
		st.commandListExternalTasks(args)

	//list all terminals
	case "lt":
		fallthrough
	case "list_term":
		fallthrough
	case "list_terms":
		st.SendCommand("list_terms", []string{})

	//add new terminal
	case "nt":
		fallthrough
	case "new_term":
		st.SendCommand("new_term", []string{})

	//ping app
	case "ping":
		st.commandAppPing(args)

	//resource usage
	case "ru":
		fallthrough
	case "res_usage":
		st.commandResourceUsage(args)

	case "r":
		fallthrough
	case "rpc":
		st.commandStart([]string{"-a", "go", "run", "rpc/cli/cli.go"})

	//shutdown running app
	case "sd":
		fallthrough
	case "shutdown":
		st.commandShutDown(args)

	//start new external task, detached running in bg by default
	case "s":
		fallthrough
	case "start":
		st.commandStart(args)

	default:
		st.PrintError("\"" + cmd + "\" is an unknown command.")

	}
}
