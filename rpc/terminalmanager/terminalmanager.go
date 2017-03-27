package terminalmanager

import (
	"github.com/corpusc/viscript/god"
	"github.com/corpusc/viscript/god/terminal"
	"github.com/corpusc/viscript/hypervisor"
	"github.com/corpusc/viscript/hypervisor/dbus"
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
