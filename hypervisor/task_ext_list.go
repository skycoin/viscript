package hypervisor

import (
	"errors"
	"strconv"

	"github.com/skycoin/viscript/msg"
)

var ExtTaskListGlobal ExtTaskList

type ExtTaskList struct {
	TaskMap map[msg.ExtTaskId]msg.ExtTaskInterface
}

func initExtTaskList() {
	ExtTaskListGlobal.TaskMap = make(map[msg.ExtTaskId]msg.ExtTaskInterface)
}

func teardownExtTaskList() {
	ExtTaskListGlobal.TaskMap = nil
	// TODO: Further cleanup
}

func ExtTaskIsRunning(procId msg.ExtTaskId) bool {
	_, exists := ExtTaskListGlobal.TaskMap[procId]
	return exists
}

func AddExtTask(ep msg.ExtTaskInterface) msg.ExtTaskId {
	id := ep.GetId()

	if !ExtTaskIsRunning(id) {
		ExtTaskListGlobal.TaskMap[id] = ep
	}

	return id
}

func GetExtTask(id msg.ExtTaskId) (msg.ExtTaskInterface, error) {
	extTask, exists := ExtTaskListGlobal.TaskMap[id]
	if exists {
		return extTask, nil
	}

	err := errors.New("External task with id " +
		strconv.Itoa(int(id)) + " doesn't exist!")

	return nil, err
}

func RemoveExtTask(id msg.ExtTaskId) {
	delete(ExtTaskListGlobal.TaskMap, id)
}

func TickExtTasks() {
	// TODO: Read from response channels if they contain any new messages
	// for _, p := range ExtTaskListGlobal.TaskMap {
	// data, err := monitor.Monitor.ReadFrom(p.GetId())
	// if err != nil {
	// 	// println(err.Error())
	// 	// monitor.Monitor.PrintAll()
	// 	continue
	// }

	// ackType := msg.GetType(data)

	// switch ackType {
	// case msg.TypeUserCommandAck:

	// }

	// select {
	// case <-p.GetTaskExitChannel():
	// 	println("Got the exit in task ext list")
	// default:
	// }
	// p.Tick()
	// }

}
