package process

import (
	"github.com/corpusc/viscript/msg"
)

//non-instanced
func NewProcess() *Process {
	println("(process/terminal/process.go).NewProcess()")
	var p Process
	p.Id = msg.NextProcessId()
	p.Type = 0
	p.Label = "TestLabel"
	p.InChannel = make(chan []byte, msg.ChannelCapacity)
	p.State.Init(&p)
	return &p
}

type Process struct {
	Id           msg.ProcessId
	Type         msg.ProcessType
	Label        string
	OutChannelId uint32
	InChannel    chan []byte
	State        State
}

func (self *Process) GetProcessInterface() msg.ProcessInterface {
	println("(process/terminal/process.go).GetProcessInterface()")
	return msg.ProcessInterface(self)
}

func (self *Process) DeleteProcess() {
	println("(process/terminal/process.go).DeleteProcess()")
	close(self.InChannel)
	self.State.proc = nil
	self = nil
}

//implement the interface

func (self *Process) GetId() msg.ProcessId {
	return self.Id
}

func (self *Process) GetType() msg.ProcessType {
	return self.Type
}

func (self *Process) GetLabel() string {
	return self.Label
}

func (self *Process) GetIncomingChannel() chan []byte {
	return self.InChannel
}

func (self *Process) Tick() {
	self.State.HandleMessages()
}
