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

var GlobalTasks Tasks

type Tasks struct {
	TaskMap map[msg.TaskId]msg.TaskInterface //task id to interface
}

func AddTask(ti msg.TaskInterface) msg.TaskId {
	id := ti.GetId()

	_, inTheMap := GlobalTasks.TaskMap[id]
	if !inTheMap {
		GlobalTasks.TaskMap[id] = ti
	}

	return id
}

func GetTaskEvents() {
	println("tasks.GetTaskEvents()   ---------------- TODO !!!!!!!!!!!")
}

func TickTasks() {
	for _, t := range GlobalTasks.TaskMap {
		t.Tick()
	}
}

//
//
//private
//
//

func initTasks() {
	GlobalTasks.TaskMap = make(map[msg.TaskId]msg.TaskInterface)
}

func teardownTasks() {
	GlobalTasks.TaskMap = nil
	// TODO: actually call teardown methods on all the tasks and also
	// external apps. what about Alt+f4?
	// upon application exit we need to terminate all the running tasks
	// and external apps
}
