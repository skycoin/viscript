package task

import (
	"github.com/skycoin/viscript/app"
	"github.com/skycoin/viscript/hypervisor"
	"github.com/skycoin/viscript/msg"
)

var path = "hypervisor/task/terminal/task"

type Task struct {
	Id           msg.TaskId
	Type         msg.TaskType
	Text         string
	OutChannelId uint32
	InChannel    chan []byte
	State        State

	hasExternalAppAttached bool
	attachedExternalApp    msg.ExternalAppInterface
}

//non-instanced
func MakeNewTask() *Task {
	println("<" + path + ">.MakeNewTask()")

	var t Task
	t.Id = msg.NextTaskId()
	t.Type = 0
	t.Text = "Initial text"
	t.InChannel = make(chan []byte, msg.ChannelCapacity)
	t.State.Init(&t)
	t.hasExternalAppAttached = false

	return &t
}

func (ta *Task) GetTaskInterface() msg.TaskInterface {
	app.At(path, "GetTaskInterface")
	return msg.TaskInterface(ta)
}

func (ta *Task) DeleteTask() {
	app.At(path, "DeleteTask")
	close(ta.InChannel)
	ta.State.task = nil
	ta = nil
}

func (ta *Task) HasExternalAppAttached() bool {
	return ta.hasExternalAppAttached
}

func (ta *Task) AttachExternalApp(eai msg.ExternalAppInterface) error {
	app.At(path, "AttachExternalApp")
	err := eai.Attach()

	if err != nil {
		return err
	}

	ta.attachedExternalApp = eai
	ta.hasExternalAppAttached = true

	return nil
}

func (ta *Task) DetachExternalApp() {
	app.At(path, "DetachExternalApp")
	// ta.attachedExternalApp.Detach()
	ta.attachedExternalApp = nil
	ta.hasExternalAppAttached = false
}

func (ta *Task) ExitExternalApp() {
	app.At(path, "ExitExternalApp")
	ta.hasExternalAppAttached = false
	id := ta.attachedExternalApp.GetId() //for removing from global list.
	ta.attachedExternalApp.TearDown()    //(and cleanup)
	ta.attachedExternalApp = nil
	hypervisor.RemoveExternalApp(id) //...from GlobalRunningExternalApps.TaskMap
}

//implement the interface

func (ta *Task) GetId() msg.TaskId {
	return ta.Id
}

func (ta *Task) GetType() msg.TaskType {
	return ta.Type
}

func (ta *Task) GetText() string {
	return ta.Text
}

func (ta *Task) GetIncomingChannel() chan []byte {
	return ta.InChannel
}

func (ta *Task) GetOutChannelId() uint32 {
	return ta.OutChannelId
}

func (ta *Task) Tick() {
	ta.State.HandleMessages()

	if !ta.HasExternalAppAttached() {
		return
	}

	select {
	//case exit := <-ta.attachedExternalApp.GetExitChannel():
	// if exit {
	// 	println("Got the exit in task, task is finished.")
	// 	//TODO: still not working yet. looking for the best way to finish
	// 	//multiple goroutines at the same time to avoid any side effects
	// 	ta.ExitExternalApp()
	// }
	case data := <-ta.attachedExternalApp.GetOutputChannel():
		println("Received data from external app, sending to term.")
		ta.State.PrintLn(string(data))
	default:
	}
}
