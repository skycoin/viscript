package process

import (
	"github.com/corpusc/viscript/msg"
)

//non-instanced
func NewProcess() *Process {
	println("(process/terminal/process.go).NewProcess()")
	var p Process

	p.Id = msg.NextProcessId()

	p.MessageIn = make(chan []byte, msg.ChannelCapacity)
	p.MessageOut = make(chan []byte, msg.ChannelCapacity)

	p.State.Init(&p)

	return &p
}

//terminal process
type Process struct {
	Id msg.ProcessId

	MessageIn  chan []byte
	MessageOut chan []byte

	State State
}

func (self *Process) GetProcessInterface() msg.ProcessInterface {
	println("(process/terminal/process.go).GetProcessInterface()")
	return msg.ProcessInterface(self)
}

func (self *Process) DeleteProcess() {
	println("(process/terminal/process.go).DeleteProcess()")
	close(self.MessageIn)
	close(self.MessageOut)
	self.State.proc = nil
	self = nil
}

//implement the interface

func (self *Process) GetId() msg.ProcessId {
	println("(process/terminal/process.go).GetId()")
	return self.Id
}

func (self *Process) GetIncomingChannel() chan []byte {
	println("(process/terminal/process.go).GetIncomingChannel()")
	return self.MessageIn
}

func (self *Process) GetOutgoingChannel() chan []byte {
	//println("(process/terminal/process.go).GetOutgoingChannel()")
	return self.MessageOut
}

//Business logic
func (self *Process) Tick() {
	//println("(process/terminal/process.go).Tick()")
	self.State.HandleMessages()
}
