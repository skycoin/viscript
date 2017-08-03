package hypervisor

import (
	"errors"
	"strconv"

	"github.com/skycoin/viscript/msg"
)

var ExtTaskListGlobal ExtTaskList

type ExtTaskList struct {
	TaskMap map[msg.ExtAppId]msg.ExtTaskInterface
}

func initExtTaskList() {
	ExtTaskListGlobal.TaskMap = make(map[msg.ExtAppId]msg.ExtTaskInterface)
}

func teardownExtTaskList() {
	ExtTaskListGlobal.TaskMap = nil
	// TODO: Further cleanup
}

func ExtTaskIsRunning(taskId msg.ExtAppId) bool {
	_, exists := ExtTaskListGlobal.TaskMap[taskId]
	return exists
}

func AddExtTask(ea msg.ExtTaskInterface) msg.ExtAppId {
	id := ea.GetId()

	if !ExtTaskIsRunning(id) {
		ExtTaskListGlobal.TaskMap[id] = ea
	}

	return id
}

func GetExtTask(id msg.ExtAppId) (msg.ExtTaskInterface, error) {
	ea, exists := ExtTaskListGlobal.TaskMap[id]
	if exists {
		return ea, nil
	}

	err := errors.New("External task with id " +
		strconv.Itoa(int(id)) + " doesn't exist!")

	return nil, err
}

func RemoveExtTask(id msg.ExtAppId) {
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
