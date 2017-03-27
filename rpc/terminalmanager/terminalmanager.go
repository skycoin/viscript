package terminalmanager

import (
	"github.com/corpusc/viscript/hypervisor"
	"github.com/corpusc/viscript/hypervisor/dbus"
	"github.com/corpusc/viscript/viewport"
	"github.com/corpusc/viscript/viewport/terminal"
)

type TerminalManager struct {
	dbus          *dbus.DbusInstance
	terminalStack *terminal.TerminalStack
	processList   *hypervisor.ProcessList
}

func newTerminalManager() *TerminalManager {
	ntm := new(TerminalManager)
	ntm.dbus = &hypervisor.DbusGlobal
	ntm.terminalStack = &god.Terms
	ntm.processList = &hypervisor.ProcessListGlobal
	return ntm
}
