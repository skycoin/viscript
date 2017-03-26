package process

import (
	"errors"

	"strconv"

	"github.com/corpusc/viscript/msg"
)

type Process struct {
	Id           msg.ProcessId
	Type         msg.ProcessType
	Label        string
	OutChannelId uint32
	InChannel    chan []byte
	State        State

	extProcAttached   bool
	extProcessId      msg.ExtProcessId
	extProcessCounter msg.ExtProcessId
	extProcesses      map[msg.ExtProcessId]*ExternalProcess
}

//non-instanced
func NewProcess() *Process {
	println("(process/terminal/process.go).NewProcess()")
	var p Process
	p.Id = msg.NextProcessId()
	p.Type = 0
	p.Label = "TestLabel"
	p.InChannel = make(chan []byte, msg.ChannelCapacity)
	p.State.Init(&p)

	// means no external process is attached
	p.extProcAttached = false
	p.extProcessId = msg.ExtProcessId(0)
	p.extProcessCounter = 0
	p.extProcesses = make(map[msg.ExtProcessId]*ExternalProcess)

	return &p
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

func (pr *Process) HasExtProcessAttached() bool {
	return pr.extProcAttached && pr.extProcessId != 0
}

func (pr *Process) GetAttachedExtProcess() (*ExternalProcess, error) {
	extProc, ok := pr.extProcesses[pr.extProcessId]
	if ok {
		return extProc, nil
	}

	return nil, errors.New("External process with id " +
		strconv.Itoa(int(pr.extProcessId)) + " doesn't exist.")
}

func (pr *Process) AddExtProcessWithCommand(cmd string) (msg.ExtProcessId, error) {
	// TODO: create new external process and add it to the map
	newExtProc, err := NewExternalProcess(&pr.State, cmd)
	if err != nil {
		return 0, err
	}
	pr.extProcessCounter += 1 // Sequential
	pr.extProcesses[pr.extProcessCounter] = newExtProc
	return pr.extProcessCounter, nil
}

func (pr *Process) AttachExtProcess(pID msg.ExtProcessId) {
	pr.extProcessId = pID
	pr.extProcAttached = true
}

func (pr *Process) AddAndAttach(cmd string) error {
	pID, err := pr.AddExtProcessWithCommand(cmd)
	if err != nil {
		return err
	}
	pr.AttachExtProcess(pID)
	return nil
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
