package terminal

import (
	"github.com/skycoin/viscript/msg"
	"strconv"
)

func (ts *TerminalStack) OnUserCommandFinalStage(receiver msg.TerminalId, cmd msg.MessageTokenizedCommand) {
	switch cmd.Command {

	case "close_term":
		ts.commandCloseTerminalFinalStage(receiver, cmd)
	case "defocus":
		ts.Defocus()
	case "focus":
		ts.commandFocusFinalStage(receiver, cmd)
	case "list_terms":
		ts.commandListTerminals(receiver, cmd)
	case "new_term":
		//temporary
		//for now, we'll be testing the difference between fixed size and dynamic terminals.
		//the 1st/initial terminal will be dynamic.  new terms afterwards will all be fixed.
		ts.AddWithFixedSizeState(true)
	default:
		println("OnUserCommandFinalStage()   UNHANDLED COMMAND!!!:", cmd.Command)
		println("OnUserCommandFinalStage()   UNHANDLED COMMAND!!!:", cmd.Command)
		println("OnUserCommandFinalStage()   UNHANDLED COMMAND!!!:", cmd.Command)
		println("OnUserCommandFinalStage()   UNHANDLED COMMAND!!!:", cmd.Command)

	}
}

//
//
//private
//
//

func (ts *TerminalStack) commandListTerminals(receiver msg.TerminalId, cmd msg.MessageTokenizedCommand) {
	var m msg.MessageTerminalIds
	m.Focused = receiver

	for _, term := range ts.TermMap {
		m.TermIds = append(m.TermIds, term.TerminalId)
	}

	ts.GetFocusedTerminal().RelayToTask(msg.Serialize(msg.TypeTerminalIds, m))
}

func (ts *TerminalStack) commandCloseTerminalFinalStage(receiver msg.TerminalId, cmd msg.MessageTokenizedCommand) {
	//for simplicity, all edge cases are to be handled in the first stage, so that...
	//THIS SHOULD NEVER BE RUN WITHOUT A VALID ID
	termToClose, _ := strconv.Atoi(cmd.Args[0]) //(should always convert)
	ts.Remove(msg.TerminalId(termToClose))
}

func (ts *TerminalStack) commandFocusFinalStage(receiver msg.TerminalId, cmd msg.MessageTokenizedCommand) {
	//foundAnyTerm := false
	matchedTerm := msg.TerminalId(0)
	for _, t := range ts.TermMap {
		ruledOutMatch := false
		arg := cmd.Args[0]
		tId := strconv.Itoa(int(t.TerminalId))

		//chop runes off the end
		//(if user gave more digits than an id has)
		for len(arg) > len(tId) {
			arg = arg[:len(arg)-1]
		}

		println("arg:", arg)

		//compare each rune of user input
		//to leftmost runes of id
		for i, c := range arg {
			println("string(c):", string(c))
			println("string(tId[i]):", string(tId[i]))

			if c != rune(tId[i]) {
				ruledOutMatch = true
				break
			}
		}

		if !ruledOutMatch {
			matchedTerm = t.TerminalId
			break
		}
	}

	//set new focus (or show error)
	if matchedTerm != 0 {
		println("Found terminal starting with:", cmd.Args[0])
		ts.SetFocused(matchedTerm)
	} else {
		println("ERROR!!!   No terminal id matches:", cmd.Args[0])
	}
}
