package hypervisor

import (
	"errors"
	"strconv"

	"github.com/skycoin/viscript/msg"
)

var GlobalRunningExternalApps RunningExternalApps

type RunningExternalApps struct {
	TaskMap map[msg.ExternalAppId]msg.ExternalAppInterface
}

func initRunningExternalApps() {
	GlobalRunningExternalApps.TaskMap = make(map[msg.ExternalAppId]msg.ExternalAppInterface)
}

func teardownRunningExternalApps() {
	GlobalRunningExternalApps.TaskMap = nil
	// TODO: Further cleanup
}

func ExternalAppIsRunning(id msg.ExternalAppId) bool {
	_, exists := GlobalRunningExternalApps.TaskMap[id]
	return exists
}

func AddExternalApp(ea msg.ExternalAppInterface) msg.ExternalAppId {
	id := ea.GetId()

	if !ExternalAppIsRunning(id) {
		GlobalRunningExternalApps.TaskMap[id] = ea
	}

	return id
}

func GetExternalApp(id msg.ExternalAppId) (msg.ExternalAppInterface, error) {
	ea, exists := GlobalRunningExternalApps.TaskMap[id]
	if exists {
		return ea, nil
	}

	err := errors.New("External app with id " +
		strconv.Itoa(int(id)) + " isn't running")

	return nil, err
}

func RemoveExternalApp(id msg.ExternalAppId) {
	delete(GlobalRunningExternalApps.TaskMap, id)
}

func TickExternalApps() {
	// TODO: Read from response channels if they contain any new messages
	// for _, p := range GlobalRunningExternalApps.TaskMap {
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
	// case <-p.GetExitChannel():
	// 	println("Got the exit (in TickExternalApps)")
	// default:
	// }
	// p.Tick()
	// }

}
