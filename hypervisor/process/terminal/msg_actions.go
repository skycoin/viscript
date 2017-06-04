package process

import (
	"fmt"

	//"github.com/corpusc/viscript/app"
	"github.com/corpusc/viscript/hypervisor"
	"github.com/corpusc/viscript/msg"
)

func (st *State) onMouseScroll(m msg.MessageMouseScroll, serializedMsg []byte) {
	if m.HoldingControl {
		hypervisor.DbusGlobal.PublishTo(st.proc.OutChannelId, serializedMsg)
	} else {
		st.Cli.AdjustBackscrollOffset(int(m.Y))
		st.drawScreenfulOfLog(st.VisualInfo)
	}
}

func (st *State) onChar(m msg.MessageChar) {
	//println("process/terminal/events.onChar()")
	st.Cli.InsertCharIfItFits(m.Char, st)
}

func (st *State) onKey(m msg.MessageKey, serializedMsg []byte) {
	switch msg.Action(m.Action) {

	case msg.Press: //one time, when key is first pressed
		//modifier key combos should never auto-repeat
		st.onOneTimeHotkeys(m)
		fallthrough
	case msg.Repeat: //constantly repeated for as long as key is pressed
		st.onRepeatableKeys(m, serializedMsg)
		st.Cli.EchoWholeCommand(st.proc.OutChannelId)

	case msg.Release:
		//most keys will probably never do anything upon release
	}
}

func (st *State) drawScreenfulOfLog(m msg.MessageVisualInfo) {
	//app.At("process/terminal/msg_action", "drawScreenfulOfLog")

	//called by
	//* viewport/term.setupNewGrid()
	//* backscrolling actions

	st.VisualInfo = m

	if //there's not a full screenful in log yet
	st.VisualInfo.CurrRow <
		st.VisualInfo.NumRows-
			st.VisualInfo.PromptRows {

		//don't allow a setting that can't be used yet.
		//it would give no visual feedback anyways.
		//later it might (in a buggy way),
		//once some random state changes.
		//by that time, the user forgets they
		//might have backscrolled (they saw no scrolling)
		st.Cli.BackscrollAmount = 0
	}

	ei /* current log entry index */ := len(st.Cli.Log) - 1 - st.Cli.BackscrollAmount
	page := []string{} //(screenful of visible text)

	//build a page (or less if term hasn't scrolled yet)
	usableRows := int(m.NumRows - m.PromptRows)
	for /* page isn't full & more entries */ len(page) < usableRows && ei >= 0 {
		ll /* last line */ := st.Cli.Log[ei]

		lineSections := []string{} //pieces of broken/divided-up lines

		x := int(m.NumColumns)
		for /* line needs breaking up */ len(ll) > int(m.NumColumns) {
			/* decrement towards start of word */
			for string(ll[x]) != " " &&
				/* still fits on 1 line */ (len(ll)-x) < int(m.NumColumns) {
				x--
			}

			lineSections = append(lineSections, ll[:x])
			ll = ll[x+1:]
		}

		//the last section, if anything remains
		if len(ll) > 0 {
			lineSections = append(lineSections, ll)
		}

		//add line or line sections to page
		for i := len(lineSections) - 1; i >= 0; i-- {
			page = append(page, lineSections[i])
		}

		ei--
	}

	for i := len(page) - 1; i >= 0; i-- {
		st.printLnAndMAYBELogIt(page[i], false)
	}

	st.Cli.EchoWholeCommand(st.proc.OutChannelId)
}

func (st *State) onTerminalIds(m msg.MessageTerminalIds) {
	st.storedTerminalIds = m.TermIds
	st.PrintLn(fmt.Sprintf("Terminals (%d) (* -> focused):", len(m.TermIds)))

	for index, termID := range m.TermIds {
		s := fmt.Sprintf("[%d] %d", index, termID)

		if termID == m.Focused {
			st.PrintLn("*" + s)
		} else {
			st.PrintLn(" " + s)
		}
	}

	st.Cli.EchoWholeCommand(st.proc.OutChannelId)
}

func (st *State) onOneTimeHotkeys(m msg.MessageKey) {
	switch m.Key {

	case msg.KeyC:
		if m.Mod == msg.GLFW_MOD_CONTROL {
			if !st.proc.HasExtProcessAttached() {
				return
			}

			//TODO: send the exit request i.e sigint or
			//just maybe kill the process?
		}

	case msg.KeyZ:
		if m.Mod == msg.GLFW_MOD_CONTROL {
			if !st.proc.HasExtProcessAttached() {
				return
			}

			st.PrintLn("Detaching External Process")
			st.proc.DetachExternalProcess()
		}

	}
}

func (st *State) onRepeatableKeys(m msg.MessageKey, serializedMsg []byte) {
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

func (st *State) actOnCommand() {
	cmd, args := st.Cli.CurrentCommandAndArgs()
	if len(cmd) < 1 {
		println("**** ERROR! ****   Command was empty!  Returning.")
		return
	}

	//if having multiple lines is such a bad thing, then you should have went all
	//the way and put it all on one line.  it is much easier to overlook,
	//and therefore fights against the whole purpose of printing it, however.
	s := "actOnCommand()      command: \"" + cmd + "\""

	for i, arg := range args {
		if i == 0 {
			s += "      args:"
		}

		s += " [" + arg + "]"
	}

	println(s)

	if st.proc.HasExtProcessAttached() {
		extProcInChannel := st.proc.attachedExtProcess.GetProcessInChannel()
		extProcInChannel <- []byte(st.Cli.CurrentCommandLine())
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

	case "c":
		fallthrough
	case "cls":
		fallthrough
	case "clear":
		st.commandClearTerminal()

	case "ls":
		st.commandListExternalTasks(args)

	//attach external task to terminal task
	case "attach":
		st.commandAttach(args)

	case "r":
		fallthrough
	case "rpc":
		st.commandStart([]string{
			"-a", "go", "run", "rpc/cli/cli.go"})

	//start new external task, detached running in bg by default
	case "s":
		fallthrough
	case "start":
		st.commandStart(args)

	case "shutdown":
		st.commandShutDown(args)

	//add new terminal
	case "n":
		fallthrough
	case "new_term":
		st.SendCommand("new_term", []string{})

	//list all terminals marking focused with •
	case "list_terms":
		st.SendCommand("list_terms", []string{})

	//delete terminal with given index to the stored termIds
	case "delete_term":
		st.deleteTerminal(args)

	//display all apps with descriptions
	case "apps":
		st.commandDisplayApps()

	default:
		st.PrintError("\"" + cmd + "\" is an unknown command.")
	}
}
