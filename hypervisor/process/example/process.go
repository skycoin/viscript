package process

import (
	"github.com/corpusc/viscript/msg"
)

/*
type ProcessInterface interface {
	GetId() ProcessId
	GetIncomingChannel() *chan []byte //channel for incoming messages
	GetOutgoingChannel() *chan []byte //channel for outgoing messages
	Tick()                            //process the messages and emit messages
}
*/

//Example Process
//implements process interface
type Process struct {
	Id msg.ProcessId

	MessageIn  chan []byte
	MessageOut chan []byte

	State State
}

func NewProcess() *Process {
	var p Process

	p.Id = msg.RandProcessId()

	p.MessageIn = make(chan []byte)
	p.MessageOut = make(chan []byte)

	p.State.InitState(&p)

	return &p
}

func (self *Process) GetProcessInterface() msg.ProcessInterface {
	var m msg.ProcessInterface
	m = msg.ProcessInterface(self)
	return m
}

func (self *Process) DeleteProcess() {

}

//implement the interface

func (self *Process) GetId() msg.ProcessId {
	return self.Id
}

func (self *Process) GetIncomingChannel() chan []byte {
	return self.MessageIn
}

func (self *Process) GetOutgoingChannel() chan []byte {
	return self.MessageOut
}

// do business logic
func (self *Process) Tick() {
	self.State.HandleMessages()
}
