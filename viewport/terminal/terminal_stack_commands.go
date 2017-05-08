package terminal

import (
	"github.com/corpusc/viscript/msg"
)

func (ts *TerminalStack) ActOnCommand(tID msg.TerminalId, command string) {
	switch command {
	case "add_new_term":
		ts.Add()
	case "list_terms":
		ts.ListTerminalsWithIds(tID)
	default:
	}
}

func (ts *TerminalStack) ListTerminalsWithIds(termId msg.TerminalId) {
	var m msg.MessageTerminalIds

	m.Focused = termId

	for _, term := range ts.Terms {
		m.TermIds = append(m.TermIds, term.TerminalId)
	}

	serializedMessage := msg.Serialize(msg.TypeTerminalIds, m)

	ts.Focused.RelayToTask(serializedMessage)
}
