package hypervisor

import (
	"github.com/corpusc/viscript/msg"
	//"github.com/corpusc/viscript/process/default"
)

var _ProcessList ProcessList

type ProcessList struct {
	ProcessMap map[msg.ProcessId]msg.ProcessInterface //process id to interface
}

func HypervisorInitProcessList() {
	_ProcessList.ProcessMap = make(map[msg.ProcessId]msg.ProcessInterface)
}

func HypervisorProcessListTeardown() {

}

func AddProcess(p msg.ProcessInterface) {
	id := p.GetId()
	//do check to make sure processId is not already in list
	_ProcessList.ProcessMap[id] = p
}

func TickProcesses() {

	for i, p := range _ProcessList.ProcessMap {
		_ = i
		p.Tick() //only do if incoming messages
	}
}

func GetProcessEvents() {
	for id, p := range _ProcessList.ProcessMap {
		var c chan []byte = p.GetOutgoingChannel() //channel
		//p.Tick() //only do if incoming messages

		if len(c) > 0 {
			m := <-c //read events
			//do events
			ProcessEventDispatch(m, id)
		}
	}

}

func ProcessEventDispatch(msg []byte, Id msg.ProcessId) {
	//process messages received by process
}
