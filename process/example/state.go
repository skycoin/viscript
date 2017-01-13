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

func (self *State) InitState(proc *Process) {
	self.proc = proc
	self.DebugPrintInputEvents = false
}

func (self *State) HandleMessages() {

	var c chan []byte = self.proc.MessageIn

	var msgType uint16
	var msgTypeMask uint16

	for len(c) > 0 {
		m := <-c // read from channel
		//route the message
		msgType = msg.GetMessageTypeUInt16(m)
		msgTypeMask = msgType & 0xff00
		switch msgTypeMask {
		case msg.MP1: //input events
			self.ProcessInputEvents(msgType, m) //input event
		case msg.MP2: //process to hypervisor messages
			log.Panic("Error: Example process does not handle this type")
		}
	}
}

//p.MessageIn = make(chan []byte)
//p.MessageOut = make(chan []byte)
