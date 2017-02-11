package process

import (
	"github.com/corpusc/viscript/msg"
	"log"
)

type State struct {
	proc                  *Process
	DebugPrintInputEvents bool
}

func (self *State) Init(proc *Process) {
	println("(process/example/state.go).Init()")
	self.proc = proc
	self.DebugPrintInputEvents = true
}

func (self *State) HandleMessages() {
	//println("(process/example/state.go).HandleMessages()")
	c := self.proc.InChannel

	for len(c) > 0 {
		m := <-c // FIXME IF this task ends up prefixing a channel id
		//(like terminal task does)
		//then we would need to do "m = m[4:]" before its multiple uses
		msgType := msg.GetType(m)
		msgTypeMask := msgType & 0xff00

		switch msgTypeMask {
		case msg.CATEGORY_Input:
			println("(process/example/state.go) -----------CATEGORY_Input")
			self.UnpackEvent(msgType, m)
		case msg.CATEGORY_Terminal:
			println("(process/example/state.go) -----------CATEGORY_Terminal")
			log.Panic("Error: Example process does NOT handle TERMINAL messages")
		default:
			println("(process/example/state.go) **************** UNHANDLED MESSAGE TYPE CATEGORY! ****************")
		}
	}
}
