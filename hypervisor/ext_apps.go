package hypervisor

import (
	"errors"
	"strconv"

	"github.com/skycoin/viscript/msg"
)

var ExternalAppListGlobal ExternalAppList

type ExternalAppList struct {
	TaskMap map[msg.ExternalAppId]msg.ExternalAppInterface
}

func initExternalAppList() {
	ExternalAppListGlobal.TaskMap = make(map[msg.ExternalAppId]msg.ExternalAppInterface)
}

func teardownExternalAppList() {
	ExternalAppListGlobal.TaskMap = nil
	// TODO: Further cleanup
}

func ExternalAppIsRunning(id msg.ExternalAppId) bool {
	_, exists := ExternalAppListGlobal.TaskMap[id]
	return exists
}

func AddExternalApp(ea msg.ExternalAppInterface) msg.ExternalAppId {
	id := ea.GetId()

	if !ExternalAppIsRunning(id) {
		ExternalAppListGlobal.TaskMap[id] = ea
	}

	return id
}

func GetExternalApp(id msg.ExternalAppId) (msg.ExternalAppInterface, error) {
	ea, exists := ExternalAppListGlobal.TaskMap[id]
	if exists {
		return ea, nil
	}

	err := errors.New("External app with id " +
		strconv.Itoa(int(id)) + " isn't running")

	return nil, err
}

func RemoveExternalApp(id msg.ExternalAppId) {
	delete(ExternalAppListGlobal.TaskMap, id)
}

func TickExternalApps() {
	// TODO: Read from response channels if they contain any new messages
	// for _, p := range ExternalAppListGlobal.TaskMap {
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
	// 	println("Got the exit (in TickExternalApps)")
	// default:
	// }
	// p.Tick()
	// }

}
