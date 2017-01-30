package process

import (
	"github.com/corpusc/viscript/msg"
)

// Process - Example process
type Process struct {
	Id msg.ProcessId

	MessageIn  chan []byte
	MessageOut chan []byte

	State State
}

// NewProcess - Constrcuts and returns a new Process struct object
func NewProcess() *Process {
	var p Process

	p.Id = msg.NextProcessId()

	p.MessageIn = make(chan []byte)
	p.MessageOut = make(chan []byte)

	p.State.InitState(&p)

	return &p
}

// GetProcessInterface - Returns process's msg.ProcessInterface
func (self *Process) GetProcessInterface() msg.ProcessInterface {
	return msg.ProcessInterface(self)
}

// DeleteProcess - Deletes the process
func (self *Process) DeleteProcess() {

}

// Implement ProcessInterface

// GetId - Returns process's sequential msg.ProcessId
func (self *Process) GetId() msg.ProcessId {
	return self.Id
}

// GetIncomingChannel - Returns process's incoming channel
func (self *Process) GetIncomingChannel() chan []byte {
	return self.MessageIn
}

// GetOutgoingChannel - Returns process's outgoing channel
func (self *Process) GetOutgoingChannel() chan []byte {
	return self.MessageOut
}

// Tick - Business logic
func (self *Process) Tick() {
	self.State.HandleMessages()
}
