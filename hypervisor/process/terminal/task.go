package process

import (
	"errors"

	"strconv"

	"github.com/corpusc/viscript/app"
	"github.com/corpusc/viscript/msg"
)

var path = "hypervisor/process/terminal/task"

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
func MakeNewTask() *Process {
	println("<" + path + ">.MakeNewTask()")

	var p Process
	p.Id = msg.NextProcessId()
	p.Type = 0
	p.Label = "TestLabel"
	p.InChannel = make(chan []byte, msg.ChannelCapacity)
	p.State.Init(&p)

	// means no external task is attached
	p.extProcAttached = false
	p.extProcessId = msg.ExtProcessId(0)
	p.extProcessCounter = 0
	p.extProcesses = make(map[msg.ExtProcessId]*ExternalProcess)

	return &p
}

func (pr *Process) GetProcessInterface() msg.ProcessInterface {
	app.At(path, "GetProcessInterface")
	return msg.ProcessInterface(pr)
}

func (pr *Process) DeleteProcess() {
	app.At(path, "DeleteProcess")
	close(pr.InChannel)
	pr.State.proc = nil
	pr = nil
}

func (pr *Process) HasExtProcessAttached() bool {
	return pr.extProcAttached && pr.extProcessId != 0
}

func (pr *Process) GetAttachedExtProcess() (*ExternalProcess, error) {
	app.At(path, "GetAttachedExtProcess")

	extProc, ok := pr.extProcesses[pr.extProcessId]
	if ok {
		return extProc, nil
	}

	return nil, errors.New("External task with id " +
		strconv.Itoa(int(pr.extProcessId)) + " doesn't exist.")
}

func (pr *Process) SendAttachedToBg() error {
	if pr.HasExtProcessAttached() {
		extProc, err := pr.GetAttachedExtProcess()
		if err != nil {
			return err
		}
		extProc.RunningInBg = true
		pr.extProcAttached = false
		pr.extProcessId = 0
	}
	return nil
}

func (pr *Process) SendExtToFg(extProcId msg.ExtProcessId) error {
	// pr.State.PrintError("Proc ID: " + string(extProcId))
	extProc, err := pr.GetAttachedExtProcess()
	if err != nil {
		return err
	}
	extProc.RunningInBg = false
	pr.extProcAttached = true
	pr.extProcessId = extProcId
	return nil
}

func (pr *Process) ExitExtProcess() error {
	if pr.HasExtProcessAttached() {
		extProc, err := pr.GetAttachedExtProcess()
		if err != nil {
			return err
		}
		extProc.shouldEnd = true
	}
	return nil
}

func (pr *Process) DeleteAttachedExtProcess() error {
	app.At(path, "DetachExtProcess")
	extProc, err := pr.GetAttachedExtProcess()
	if err != nil {
		return err
	}
	pr.extProcessId = 0
	pr.extProcAttached = false
	extProc.TearDown()
	delete(pr.extProcesses, pr.extProcessId)
	return nil
}

func (pr *Process) AddTaskExternal(tokens []string) (msg.ExtProcessId, error) {
	app.At(path, "AddTaskExternal")

	newExtProc, err := MakeNewTaskExternal(&pr.State, tokens)
	if err != nil {
		return 0, err
	}

	pr.extProcessCounter += 1 // Sequential
	pr.extProcesses[pr.extProcessCounter] = newExtProc
	return pr.extProcessCounter, nil
}

func (pr *Process) AttachExtProcess(pID msg.ExtProcessId) {
	app.At(path, "AttachExtProcess")
	pr.extProcessId = pID
	pr.extProcAttached = true
}

func (pr *Process) AddAndAttach(tokens []string) error {
	app.At(path, "AddAndAttach")

	pID, err := pr.AddTaskExternal(tokens)
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
