package process

import (
	"github.com/corpusc/viscript/msg"
)

//non-instanced
func NewProcess() *Process {
	println("(process/example/process.go).NewProcess()")
	var p Process
	p.Id = msg.NextProcessId()
	p.InChannel = make(chan []byte, msg.ChannelCapacity)
	p.State.Init(&p)
	return &p
}

type Process struct {
	Id           msg.ProcessId
	OutChannelId uint32
	InChannel    chan []byte
	State        State
}

func (self *Process) GetProcessInterface() msg.ProcessInterface {
	println("(process/example/process.go).GetProcessInterface()")
	return msg.ProcessInterface(self)
}

func (self *Process) DeleteProcess() {
	println("(process/example/process.go).DeleteProcess()")
	close(self.InChannel)
	self.State.proc = nil
	self = nil
}

// Implement ProcessInterface

func (self *Process) GetId() msg.ProcessId {
	return self.Id
}

func (self *Process) GetIncomingChannel() chan []byte {
	return self.InChannel
}

func (self *Process) Tick() {
	self.State.HandleMessages()
}
