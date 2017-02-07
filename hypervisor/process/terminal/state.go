package process

import (
	"github.com/corpusc/viscript/msg"
	//"log"
)

//put all your process state here
type State struct {
	proc                  *Process
	DebugPrintInputEvents bool
}

func (self *State) Init(proc *Process) {
	println("(process/terminal/state.go).Init()")
	self.proc = proc
	self.DebugPrintInputEvents = false
}

func (self *State) HandleMessages() {
	//println("(process/terminal/state.go).HandleMessages()")
	var c chan []byte = self.proc.MessageIn
	var msgType uint16
	var msgTypeMask uint16

	for len(c) > 0 {
		println("(process/terminal/state.go).HandleMessages() - ...ONE ITERATION OF CHANNEL")
		m := <-c // read from channel
		//route the message
		msgType = msg.GetType(m)
		msgTypeMask = msgType & 0xff00

		switch msgTypeMask {
		case msg.PrefixInput:
			self.UnpackInputEvents(msgType, m)
		case msg.PrefixTerminal: //process to hypervisor messages
			self.UnpackInputEvents(msgType, m)
		}
	}
}
