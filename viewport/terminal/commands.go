package terminal

import (
	"github.com/skycoin/viscript/msg"
)

func (ts *TerminalStack) OnUserCommandFinalStage(tID msg.TerminalId, cmd msg.MessageTokenizedCommand) {
	switch cmd.Command {

	case "close_term":
		ts.commandCloseTerminalFinalStage(tID)
	case "list_terms":
		ts.ListTerminalsWithIds(tID)
	case "new_term":
		ts.AddWithFixedSizeState(true)
	default:
		//nope

	}
}

func (ts *TerminalStack) ListTerminalsWithIds(termId msg.TerminalId) {
	var m msg.MessageTerminalIds
	m.Focused = termId

	for _, term := range ts.Terms {
		m.TermIds = append(m.TermIds, term.TerminalId)
	}

	ts.Focused.RelayToTask(msg.Serialize(msg.TypeTerminalIds, m))
}

func (ts *TerminalStack) commandCloseTerminalFinalStage(id msg.TerminalId) {
	//for simplicity, all edge cases are to be handled in the first stage, so that...
	//THIS SHOULD NEVER BE RUN WITHOUT A VALID ID
	ts.Remove(id)
}
