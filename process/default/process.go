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

type DefaultProcess struct {
	Id msg.ProcessId

	MessageIn  chan []byte
	MessageOut chan []byte

	State State
}

func NewDefaultProcess() *DefaultProcess {
	var p DefaultProcess

	p.Id = msg.RandProcessId()

	p.MessageIn = make(chan []byte)
	p.MessageOut = make(chan []byte)

	state.InitState(&p)

	return &p
}

func (self *DefaultProcess) GetProcessInterface() msg.ProcessInterface {
	var m msg.ProcessInterface
	m = msg.ProcessInterface(self)
	return m
}

func (self *DefaultProcess) DeleteDefaultProces() {

}

//implement the interface

func (self *DefaultProcess) GetId() msg.ProcessId {
	return self.Id
}

func (self *DefaultProcess) GetIncomingChannel() chan []byte {
	return self.MessageIn
}

func (self *DefaultProcess) GetOutgoingChannel() chan []byte {
	return self.MessageOut
}

// do bussiness logic
func (self *DefaultProcess) Tick() {

}
