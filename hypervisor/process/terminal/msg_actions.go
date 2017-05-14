package process

import (
	"bytes"
	"fmt"

	"github.com/corpusc/viscript/app"
	"github.com/corpusc/viscript/hypervisor"
	"github.com/corpusc/viscript/msg"
)

func (st *State) onMouseScroll(m msg.MessageMouseScroll, serializedMsg []byte) {
	if m.HoldingControl {
		hypervisor.DbusGlobal.PublishTo(st.proc.OutChannelId, serializedMsg)
	} else {
		st.Cli.AdjustBackscrollOffset(int(m.Y))
		st.onVisualInfo(st.VisualInfo) //refresh screen
		//FIXME: need to clear screen
	}
}

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

func (st *State) onVisualInfo(m msg.MessageVisualInfo) {
	//current position was reset to home (top left corner) inside viewport/term
	st.VisualInfo = m
	cl /* current log entry */ := len(st.Cli.Log) - 1 - st.Cli.BackscrollAmount
	page := []string{}

	//build a page (or less if term hasn't scrolled yet)
	for /* page isn't full & still more entries */ len(page) < int(m.NumRows) && cl >= 0 {
		ll /* last line */ := st.Cli.Log[cl]

		lineSections := []string{}

		x := int(m.NumColumns)
		for /* line needs breaking up */ len(ll) > int(m.NumColumns) {
			for /* decrementing towards start of word */ string(ll[x]) != " " &&
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

		cl--
	}

	for i := len(page) - 1; i >= 0; i-- {
		st.printLnAndMAYBELogIt(page[i], false)
	}

	st.Cli.EchoWholeCommand(st.proc.OutChannelId)
}

func (st *State) onTerminalIds(m msg.MessageTerminalIds) {
	var buffer bytes.Buffer

	header := fmt.Sprintf("Terminals (%d) (* -> focused):\n", len(m.TermIds))
	buffer.WriteString(header)

	for index, termID := range m.TermIds {
		var termDescription string
		if termID == m.Focused {
			termDescription = fmt.Sprintf("*[%d] %d\n", index, termID)
		} else {
			termDescription = fmt.Sprintf(" [%d] %d\n", index, termID)
		}
		buffer.WriteString(termDescription)
	}

	st.storedTerminalIds = m.TermIds
	st.PrintLn(buffer.String())
}

func (st *State) actOnOneTimeHotkeys(m msg.MessageKey) {
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

func (st *State) actOnCommand() {
	cmd, args := st.Cli.CurrentCommandAndArgs()
	if len(cmd) < 1 {
		println("**** ERROR! ****   Command was empty!  Returning.")
		return
	}

	//if having multiple lines is such a bad thing, then you should have went all
	//the way and put it all on one line.  that is much easier to overlook,
	//and therefore fights against the whole purpose of printing it, however.
	//also your feedback didn't read sensibly (in other ways) in my opinion
	s := "actOnCommand()      command \"" + cmd + "\"      args"

	for _, arg := range args {
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
		st.commandHelp()

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
		st.SendCommand("add_new_term", []string{})

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
