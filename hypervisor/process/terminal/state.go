package process

import (
	"github.com/corpusc/viscript/msg"
	//"log"
)

type State struct {
	proc                  *Process
	DebugPrintInputEvents bool
	Cli                   *Cli
}

func (st *State) Init(proc *Process) {
	println("(process/terminal/state.go).Init()")
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
		m = m[4:] //for now, DISCARD the chan id prefix
		msgType := msg.GetType(m)
		msgCategory := msgType & 0xff00 // get back masked category

		switch msgCategory {
		case msg.CATEGORY_Input:
			println("(process/terminal/state.go)-----------CATEGORY_Input")
			st.UnpackEvent(msgType, m)
		case msg.CATEGORY_Terminal:
			println("(process/terminal/state.go)-----------CATEGORY_Terminal")
			st.UnpackEvent(msgType, m)
		default:
			println("(process/terminal/state.go)**************** UNHANDLED MESSAGE TYPE CATEGORY! ****************")
		}
	}
}
