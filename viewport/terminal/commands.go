package terminal

import (
	"github.com/skycoin/viscript/msg"
)

func (ts *TerminalStack) OnUserCommandFinalStage(tID msg.TerminalId, cmd msg.MessageTokenizedCommand) {
	switch cmd.Command {

	case "new_term":
		ts.AddWithFixedSizeState(true)
	case "list_terms":
		ts.ListTerminalsWithIds(tID)
	case "del_term":
		ts.commandDeleteTerminalsFinalStage(tID)
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

//for simplicity, we'll try to handle all the edge cases in the first stage, so that
//THIS SHOULD NEVER BE RUN WITHOUT A VALID ID
func (ts *TerminalStack) commandDeleteTerminalsFinalStage(id msg.TerminalId) {
	ts.Remove(id)
}
