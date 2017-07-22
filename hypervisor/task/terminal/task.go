package process

import (
	"github.com/skycoin/viscript/app"
	"github.com/skycoin/viscript/hypervisor"
	"github.com/skycoin/viscript/msg"
)

var path = "hypervisor/task/terminal/task"

type Task struct {
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
func MakeNewTask() *Task {
	println("<" + path + ">.MakeNewTask()")

	var t Task
	t.Id = msg.NextTaskId()
	t.Type = 0
	t.Label = "TestLabel"
	t.InChannel = make(chan []byte, msg.ChannelCapacity)
	t.State.Init(&t)

	//means no external task is attached
	t.hasExtProcAttached = false

	return &t
}

func (ta *Task) GetTaskInterface() msg.TaskInterface {
	app.At(path, "GetTaskInterface")
	return msg.TaskInterface(ta)
}

func (ta *Task) DeleteProcess() {
	app.At(path, "DeleteProcess")
	close(ta.InChannel)
	ta.State.task = nil
	ta = nil
}

func (ta *Task) HasExtTaskAttached() bool {
	return ta.hasExtProcAttached
}

func (ta *Task) AttachExternalTask(extProc msg.ExtTaskInterface) error {
	app.At(path, "AttachExternalTask")
	err := extProc.Attach()
	if err != nil {
		return err
	}

	ta.attachedExtTask = extProc
	ta.hasExtProcAttached = true

	return nil
}

func (ta *Task) DetachExternalTask() {
	app.At(path, "DetachExternalTask")
	// ta.attachedExtTask.Detach()
	ta.attachedExtTask = nil
	ta.hasExtProcAttached = false
}

func (ta *Task) ExitExtTask() {
	app.At(path, "ExitExtTask")
	ta.hasExtProcAttached = false
	extProcId := ta.attachedExtTask.GetId() //for removing from global list.
	ta.attachedExtTask.TearDown()           //(and cleanup)
	ta.attachedExtTask = nil
	hypervisor.RemoveExtTask(extProcId) //...from ExtTaskListGlobal.TaskMap
}

//implement the interface

func (ta *Task) GetId() msg.TaskId {
	return ta.Id
}

func (ta *Task) GetType() msg.TaskType {
	return ta.Type
}

func (ta *Task) GetLabel() string {
	return ta.Label
}

func (ta *Task) GetIncomingChannel() chan []byte {
	return ta.InChannel
}

func (ta *Task) Tick() {
	ta.State.HandleMessages()

	if !ta.HasExtTaskAttached() {
		return
	}

	select {
	//case exit := <-ta.attachedExtTask.GetTaskExitChannel():
	// if exit {
	// 	println("Got the exit in task, task is finished.")
	// 	//TODO: still not working yet. looking for the best way to finish
	// 	//multiple goroutines at the same time to avoid any side effects
	// 	ta.ExitExtTask()
	// }
	case data := <-ta.attachedExtTask.GetTaskOutChannel():
		println("Received data from external task, sending to term.")
		ta.State.PrintLn(string(data))
	default:
	}
}
