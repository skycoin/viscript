package api

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

func (pr *Process) GetProcessInterface() msg.ProcessInterface {
	println("(process/terminal/process.go).GetProcessInterface()")
	return msg.ProcessInterface(pr)
}

func (pr *Process) DeleteProcess() {
	println("(process/terminal/process.go).DeleteProcess()")
	close(pr.InChannel)
	pr.State.proc = nil
	pr = nil
}

//implement the interface

func (pr *Process) GetId() msg.ProcessId {
	return pr.Id
}

func (pr *Process) GetType() msg.ProcessType {
	return pr.Type
}

func (pr *Process) GetLabel() string {
	return pr.Label
}

func (pr *Process) GetIncomingChannel() chan []byte {
	return pr.InChannel
}

func (pr *Process) Tick() {
	pr.State.HandleMessages()
}
