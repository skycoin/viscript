package terminal

import (
	"strconv"

	"github.com/corpusc/viscript/msg"
)

func (ts *TerminalStack) ActOnCommand(tID msg.TerminalId, command msg.MessageCommand) {
	switch command.Command {
	case "add_new_term":
		ts.Add()
	case "list_terms":
		ts.ListTerminalsWithIds(tID)
	case "delete_term":
		if len(command.Args) != 1 {
			return
		}
		ts.DeleteTerminalIfExists(command.Args[0])
	case "clear":
		ts.ClearCurrentlyFocusedTerminal()
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

func (ts *TerminalStack) DeleteTerminalIfExists(termIdString string) {
	termIdInt, err := strconv.Atoi(termIdString)
	if err != nil {
		return
	}

	ts.RemoveTerminal(msg.TerminalId(termIdInt))
}

func (ts *TerminalStack) ClearCurrentlyFocusedTerminal() {
	// TODO: tried ts.focused.PrintPrompt() and Clear() but doesn't work
	// currently clearing by sending 100 \n characters from the hypervisor
}
