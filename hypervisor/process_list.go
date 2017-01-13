package hypervisor

import (
	"github.com/corpusc/viscript/msg"
	//"github.com/corpusc/viscript/process/example"
)

/*
	Process can
	- receive messages from hypervisor
	- can emite messages back to hyper visor
	- have a tick method for processing incoming messages

	Incoming messages from process to hypervisor come in anytime
	- on input dispatch
	Messages outgoing to hypervisor are processed in HandleDispatchProcesEvents()
	- we check each process channel for outgoing messages to the hypervisor
	Each process has a tick() method
	- tick method, input messages are processed, output messages created


*/
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

func GetProcessEvents() {

}

//run the process, creating new events for hypervisor
func ProcessTick() {
	for id, p := range _ProcessList.ProcessMap {
		_ = id
		p.Tick() //only do if incoming messages
	}
}

//events from process to hypervisor
func DispatchProcesEvents() {
	for id, p := range _ProcessList.ProcessMap {
		var c chan []byte = p.GetOutgoingChannel() //channel
		//p.Tick() //only do if incoming messages

		if len(c) > 0 {
			m := <-c //read events
			//do events
			ProcessEvent(m, id)
		}
	}
}

//an incoming message from a process
func ProcessEvent(msg []byte, Id msg.ProcessId) {
	//process messages received by process
}
