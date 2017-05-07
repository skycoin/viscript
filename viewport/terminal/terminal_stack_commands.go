package terminal

func (ts *TerminalStack) ActOnCommand(command string) {
	switch command {
	case "new":
		ts.Add()
	case "list":
		ts.ListTerminalsWithIds()
	default:
	}
}

func (ts *TerminalStack) ListTerminalsWithIds() {

}
