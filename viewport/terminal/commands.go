package terminal

import (
	"github.com/skycoin/viscript/msg"
	"strconv"
)

func (ts *TerminalStack) OnUserCommandFinalStage(reciever msg.TerminalId, cmd msg.MessageTokenizedCommand) {
	switch cmd.Command {

	case "close_term":
		ts.commandCloseTerminalFinalStage(reciever, cmd)
	case "defocus":
		ts.Defocus()
	case "focus":
		//
	case "list_terms":
		ts.commandListTerminals(reciever)
	case "new_term":
		//temporary
		//for now, we'll be testing the difference between fixed size and dynamic terminals.
		//the 1st/initial terminal will be dynamic.  new terms afterwards will all be fixed.
		ts.AddWithFixedSizeState(true)
	default:
		//nope

	}
}

//
//
//private
//
//

func (ts *TerminalStack) commandListTerminals(reciever msg.TerminalId) {
	var m msg.MessageTerminalIds
	m.Focused = reciever

	for _, term := range ts.TermMap {
		m.TermIds = append(m.TermIds, term.TerminalId)
	}

	ts.Focused.RelayToTask(msg.Serialize(msg.TypeTerminalIds, m))
}

func (ts *TerminalStack) commandCloseTerminalFinalStage(reciever msg.TerminalId, cmd msg.MessageTokenizedCommand) {
	//for simplicity, all edge cases are to be handled in the first stage, so that...
	//THIS SHOULD NEVER BE RUN WITHOUT A VALID ID
	termToClose, _ := strconv.Atoi(cmd.Args[0]) //(should always convert)
	ts.Remove(msg.TerminalId(termToClose))
}
