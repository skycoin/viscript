package process

import (
	"github.com/corpusc/viscript/msg"
	//"log"
)

type State struct {
	proc                  *Process
	DebugPrintInputEvents bool
}

func (self *State) Init(proc *Process) {
	println("(process/terminal/state.go).Init()")
	self.proc = proc
	self.DebugPrintInputEvents = true
}

func (self *State) HandleMessages() {
	//called per Tick()
	c := self.proc.InChannel

	for len(c) > 0 {
		m := <-c
		//TODO/FIXME:   cache channel id wherever it may be needed
		m = m[4:] //for now, DISCARD the chan id prefix
		msgType := msg.GetType(m)
		// println("msgType", msgType)
		msgTypeMask := msgType & 0xff00
		// println("msgTypeMask", msgTypeMask)

		switch msgTypeMask {
		case msg.CATEGORY_Input:
			println("(process/terminal/state.go)-----------CATEGORY_Input")
			self.UnpackEvent(msgType, m)
		case msg.CATEGORY_Terminal:
			println("(process/terminal/state.go)-----------CATEGORY_Terminal")
			self.UnpackEvent(msgType, m)
		default:
			println("(process/terminal/state.go)**************** UNHANDLED MESSAGE TYPE CATEGORY! ****************")
		}
	}
}
