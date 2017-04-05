package hypervisor

import (
	"github.com/corpusc/viscript/msg"
)

var ExtProcessListGlobal ExtProcessList

type ExtProcessList struct {
	ProcessMap map[msg.ExtProcessId]msg.ExtProcessInterface
}

func InitExtProcessList() {
	ExtProcessListGlobal.ProcessMap = make(map[msg.ExtProcessId]msg.ExtProcessInterface)
}

func TeardownExtProcessList() {
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

func TickExtTasks() {
	for _, p := range ExtProcessListGlobal.ProcessMap {
		p.Tick()
	}
}
