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
}

func NewDefaultProcess() *msg.ProccesId {
	var p DefaultProcess

	p.Id = msg.RandProcessId()

	p.MessageIn = new(chan []byte)
	p.MessageOut = new(chan []byte)

	return &p
}

func (self *DefaultProcess) DeleteDefaultProces() {

}

func (self *DefaultProcess) GetIncomingChannel() *chan []byte {
	return &self.MessageIn
}

func (self *DefaultProcess) GetOutgoingChannel() *chan []byte {
	return &self.MessageOut
}

// do bussiness logic
func (self *DefaultProcess) Tick() {

}
