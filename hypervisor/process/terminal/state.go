package process

import (
	"github.com/corpusc/viscript/app"
	"github.com/corpusc/viscript/msg"
)

var stPath = "hypervisor/process/terminal/state"

type State struct {
	DebugPrintInputEvents bool
	Cli                   *Cli
	proc                  *Process
}

func (st *State) Init(proc *Process) {
	st.proc = proc
	st.DebugPrintInputEvents = true
	st.Cli = NewCli()
}

func (st *State) HandleMessages() {
	//called per Tick()
	c := st.proc.InChannel

	for len(c) > 0 {
		m := <-c
		//TODO/FIXME:   cache channel id wherever it may be needed
		m = m[4:] //.....for now, DISCARD the chan id prefix
		msgType := msg.GetType(m)
		msgCategory := msgType & 0xff00 // get back masked category

		switch msgCategory {

		case msg.CATEGORY_Input:
			st.UnpackMessage(msgType, m)
		case msg.CATEGORY_Terminal:
			app.At(stPath, "----CATEGORY_Terminal----FIXME? SHOULD NEVER GET THIS TYPE HERE?!?!")
			st.UnpackMessage(msgType, m)
		default:
			app.At(stPath, "**************** UNHANDLED MESSAGE TYPE CATEGORY! ****************")

		}
	}
}
