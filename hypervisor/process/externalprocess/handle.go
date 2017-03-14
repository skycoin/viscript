package externalprocess

import "strings"

func (st *ExtState) actOnCommand(command []string) {
	if len(strings.ToLower(command[0])) > 0 {
		println("ActOnCommand() ->", strings.ToLower(command[0]))
		st.ExtProc.CmdOut <- []byte(strings.ToLower(command[0]))
	}
}
