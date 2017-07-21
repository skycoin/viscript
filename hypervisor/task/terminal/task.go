package process

import (
	"github.com/skycoin/viscript/app"
	"github.com/skycoin/viscript/hypervisor"
	"github.com/skycoin/viscript/msg"
)

var path = "hypervisor/task/terminal/task"

type Process struct {
	Id           msg.TaskId
	Type         msg.ProcessType
	Label        string
	OutChannelId uint32
	InChannel    chan []byte
	State        State

	hasExtProcAttached bool
	attachedExtProcess msg.ExtTaskInterface
}

//non-instanced
func MakeNewTask() *Process {
	println("<" + path + ">.MakeNewTask()")

	var p Process
	p.Id = msg.NextTaskId()
	p.Type = 0
	p.Label = "TestLabel"
	p.InChannel = make(chan []byte, msg.ChannelCapacity)
	p.State.Init(&p)

	//means no external task is attached
	p.hasExtProcAttached = false

	return &p
}

func (pr *Process) GetTaskInterface() msg.TaskInterface {
	app.At(path, "GetTaskInterface")
	return msg.TaskInterface(pr)
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

func (pr *Process) AttachExternalTask(extProc msg.ExtTaskInterface) error {
	app.At(path, "AttachExternalTask")
	err := extProc.Attach()
	if err != nil {
		return err
	}

	pr.attachedExtProcess = extProc
	pr.hasExtProcAttached = true

	return nil
}

func (pr *Process) DetachExternalTask() {
	app.At(path, "DetachExternalTask")
	// pr.attachedExtProcess.Detach()
	pr.attachedExtProcess = nil
	pr.hasExtProcAttached = false
}

func (pr *Process) ExitExtProcess() {
	app.At(path, "ExitExtProcess")
	pr.hasExtProcAttached = false
	extProcId := pr.attachedExtProcess.GetId() //for removing from global list.
	pr.attachedExtProcess.TearDown()           //(and cleanup)
	pr.attachedExtProcess = nil
	hypervisor.RemoveExtProcess(extProcId) //...from ExtTaskListGlobal.TaskMap
}

//implement the interface

func (pr *Process) GetId() msg.TaskId {
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

	if !pr.HasExtProcessAttached() {
		return
	}

	select {
	//case exit := <-pr.attachedExtProcess.GetTaskExitChannel():
	// if exit {
	// 	println("Got the exit in task, process is finished.")
	// 	//TODO: still not working yet. looking for the best way to finish
	// 	//multiple goroutines at the same time to avoid any side effects
	// 	pr.ExitExtProcess()
	// }
	case data := <-pr.attachedExtProcess.GetTaskOutChannel():
		println("Received data from external process, sending to term.")
		pr.State.PrintLn(string(data))
	default:
	}
}
