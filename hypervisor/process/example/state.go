package process

import (
	"github.com/corpusc/viscript/msg"
	"log"
)

//put all your process state here
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
	c := self.proc.MessageIn
	var msgType uint16
	var msgTypeMask uint16

	for len(c) > 0 {
		m := <-c //read from channel
		//route the message
		msgType = msg.GetType(m)
		msgTypeMask = msgType & 0xff00

		switch msgTypeMask {
		case msg.TypePrefix_Input:
			self.UnpackInputEvents(msgType, m)
		case msg.TypePrefix_Terminal: //process to hypervisor messages
			log.Panic("Error: Example process does not handle this type")
		}
	}
}
