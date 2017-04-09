package hypervisor

import (
	"errors"
	"strconv"

	"github.com/corpusc/viscript/msg"
)

var ExtProcessListGlobal ExtProcessList

type ExtProcessList struct {
	ProcessMap map[msg.ExtProcessId]msg.ExtProcessInterface
}

func initExtProcessList() {
	ExtProcessListGlobal.ProcessMap = make(map[msg.ExtProcessId]msg.ExtProcessInterface)
}

func teardownExtProcessList() {
	ExtProcessListGlobal.ProcessMap = nil
	// TODO: Further cleanup
}

func AddExtProcess(ep msg.ExtProcessInterface) msg.ExtProcessId {
	id := ep.GetId()

	_, isInTheMap := ExtProcessListGlobal.ProcessMap[id]
	if !isInTheMap {
		ExtProcessListGlobal.ProcessMap[id] = ep
	}

	return id
}

func GetExtProcess(id msg.ExtProcessId) (msg.ExtProcessInterface, error) {
	extProc, exists := ExtProcessListGlobal.ProcessMap[id]
	if exists {
		return extProc, nil
	}
	err := errors.New("External process with id " +
		strconv.Itoa(int(id)) + " doesn't exist!")
	return nil, err
}

func TickExtTasks() {
	for _, p := range ExtProcessListGlobal.ProcessMap {
		if !p.GetRunningInBg() {
			p.Tick()
		}
	}
}
