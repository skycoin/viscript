package terminalmanager

import (
	"github.com/corpusc/viscript/hypervisor"
	"github.com/corpusc/viscript/viewport"
	"github.com/corpusc/viscript/viewport/terminal"
)

type TerminalManager struct {
	terminalStack *terminal.TerminalStack
	processList   *hypervisor.ProcessList
}

func newTerminalManager() *TerminalManager {
	ntm := new(TerminalManager)
	ntm.terminalStack = &viewport.Terms
	ntm.processList = &hypervisor.ProcessListGlobal
	return ntm
}
