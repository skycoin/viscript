package externalprocess

import (
	"github.com/corpusc/viscript/hypervisor/process/api"
)

type ExtState struct {
	api.State
	ExtProc *ExternalProcess
}

func (st *ExtState) Init(p *api.Process, proc *ExternalProcess) {
	println("(process/externalprocess/state.go).Init()")
	st.State.Init(p)
	st.ExtProc = proc
	st.DebugPrintInputEvents = true
}
