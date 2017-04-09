package process

import (
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

	hasExtProcAttached bool
	attachedExtProcess msg.ExtProcessInterface
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
	p.hasExtProcAttached = false

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
	return pr.hasExtProcAttached
}

func (pr *Process) AttachExternalProcess(extProc msg.ExtProcessInterface) {
	app.At(path, "AttachExternalProcess")
	pr.attachedExtProcess = extProc
	pr.hasExtProcAttached = true
	pr.attachedExtProcess.Attach()
}

func (pr *Process) DetachExternalProcess() {
	app.At(path, "DetachExternalProcess")
	pr.attachedExtProcess = nil
	pr.hasExtProcAttached = false
	pr.attachedExtProcess.Detach()
}

// In this case process should:
// 1) Send true to the ProcessShouldEnd channel
// 2) Call external process's TearDown
// 3) Remove it from the GlobalList and further cleanup whatever will be needed

// func (pr *Process) ExitExtProcess() error {
// 	_, err := pr.GetAttachedExtProcess()
// 	if err != nil {
// 		return err
// 	}
// 	pr.DeleteAttachedExtProcess()
// 	return nil
// }

// func (pr *Process) DeleteAttachedExtProcess() error {
// 	app.At(path, "DeleteAttachedExtProcess")

// 	extProc, err := pr.GetAttachedExtProcess()
// 	if err != nil {
// 		return err
// 	}

// 	pr.extProcessId = 0
// 	pr.extProcAttached = false
// 	extProc.ShutDown()
// 	delete(pr.extProcesses, pr.extProcessId)
// 	return nil
// }

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

	if !pr.hasExtProcAttached {
		return
	}

	select {
	case exit := <-pr.attachedExtProcess.GetProcessExitChannel():
		if exit {
			println("Got the exit in task, process is finished.")
			// TODO: Exit here
		}
	case data := <-pr.attachedExtProcess.GetProcessOutChannel():
		println("Received data from external process, sending to term.")
		pr.State.PrintLn(string(data))
	default:
	}
}
