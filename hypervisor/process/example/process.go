package process

import (
	"github.com/corpusc/viscript/msg"
)

//Example process
type Process struct {
	Id msg.ProcessId

	MessageIn  chan []byte
	MessageOut chan []byte

	State State
}

func NewProcess() *Process {
	var p Process

	p.Id = msg.NextProcessId()

	p.MessageIn = make(chan []byte)
	p.MessageOut = make(chan []byte)

	p.State.InitState(&p)

	return &p
}

func (self *Process) GetProcessInterface() msg.ProcessInterface {
	return msg.ProcessInterface(self)
}

func (self *Process) DeleteProcess() {

}

// Implement ProcessInterface

func (self *Process) GetId() msg.ProcessId {
	return self.Id
}

func (self *Process) GetIncomingChannel() chan []byte {
	return self.MessageIn
}

func (self *Process) GetOutgoingChannel() chan []byte {
	return self.MessageOut
}

//Business logic
func (self *Process) Tick() {
	self.State.HandleMessages()
}
