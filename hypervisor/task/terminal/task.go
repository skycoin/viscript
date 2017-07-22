package process

import (
	"github.com/skycoin/viscript/app"
	"github.com/skycoin/viscript/hypervisor"
	"github.com/skycoin/viscript/msg"
)

var path = "hypervisor/task/terminal/task"

type Process struct {
	Id           msg.TaskId
	Type         msg.TaskType
	Label        string
	OutChannelId uint32
	InChannel    chan []byte
	State        State

	hasExtProcAttached bool
	attachedExtTask    msg.ExtTaskInterface
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
	pr.State.task = nil
	pr = nil
}

func (pr *Process) HasExtTaskAttached() bool {
	return pr.hasExtProcAttached
}

func (pr *Process) AttachExternalTask(extProc msg.ExtTaskInterface) error {
	app.At(path, "AttachExternalTask")
	err := extProc.Attach()
	if err != nil {
		return err
	}

	pr.attachedExtTask = extProc
	pr.hasExtProcAttached = true

	return nil
}

func (pr *Process) DetachExternalTask() {
	app.At(path, "DetachExternalTask")
	// pr.attachedExtTask.Detach()
	pr.attachedExtTask = nil
	pr.hasExtProcAttached = false
}

func (pr *Process) ExitExtTask() {
	app.At(path, "ExitExtTask")
	pr.hasExtProcAttached = false
	extProcId := pr.attachedExtTask.GetId() //for removing from global list.
	pr.attachedExtTask.TearDown()           //(and cleanup)
	pr.attachedExtTask = nil
	hypervisor.RemoveExtTask(extProcId) //...from ExtTaskListGlobal.TaskMap
}

//implement the interface

func (pr *Process) GetId() msg.TaskId {
	return pr.Id
}

func (pr *Process) GetType() msg.TaskType {
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

	if !pr.HasExtTaskAttached() {
		return
	}

	select {
	//case exit := <-pr.attachedExtTask.GetTaskExitChannel():
	// if exit {
	// 	println("Got the exit in task, task is finished.")
	// 	//TODO: still not working yet. looking for the best way to finish
	// 	//multiple goroutines at the same time to avoid any side effects
	// 	pr.ExitExtTask()
	// }
	case data := <-pr.attachedExtTask.GetTaskOutChannel():
		println("Received data from external task, sending to term.")
		pr.State.PrintLn(string(data))
	default:
	}
}
