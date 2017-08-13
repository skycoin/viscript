package hypervisor

import (
	"github.com/skycoin/viscript/msg"
)

/*
	Task can
	- receive messages from hypervisor
	- can emit message back to hypervisor
	- have a tick method for handling incoming messages

	Incoming messages from tasks to hypervisor come in anytime
	- on input dispatch
	- we check each task channel for outgoing messages to the hypervisor
	Each task has a tick() method
	- tick method, input messages are digested, output messages created
*/

var TaskListGlobal TaskList

type TaskList struct {
	TaskMap map[msg.TaskId]msg.TaskInterface //task id to interface
}

func initTaskList() {
	TaskListGlobal.TaskMap = make(map[msg.TaskId]msg.TaskInterface)
}

func teardownTaskList() {
	TaskListGlobal.TaskMap = nil
	// TODO: actually call teardown methods on all the tasks and also
	// external apps. what about Alt+f4?
	// upon application exit we need to terminate all the running tasks
	// and external apps
}

func AddTask(p msg.TaskInterface) msg.TaskId {
	id := p.GetId()

	_, isInTheMap := TaskListGlobal.TaskMap[id]
	if !isInTheMap {
		TaskListGlobal.TaskMap[id] = p
	}
	return id
}

func GetTaskEvents() {
	println("task_list.GetTaskEvents()   ---------------- TODO !!!!!!!!!!!")
}

func TickTasks() {
	for _, p := range TaskListGlobal.TaskMap {
		p.Tick()
	}
}
