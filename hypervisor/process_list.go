package hypervisor

import (
	example "github.com/corpusc/viscript/hypervisor/process/example"
	"github.com/corpusc/viscript/msg"
)

/*
	Process can
	- receive messages from hypervisor
	- can emite message back to hypervisor
	- have a tick method for handling incoming messages

	Incoming messages from process to hypervisor come in anytime
	- on input dispatch
	Messages outgoing to hypervisor are processed in DispatchProcessEvents()
	- we check each process channel for outgoing messages to the hypervisor
	Each process has a tick() method
	- tick method, input messages are processed, output messages created
*/

var ProcessListGlobal ProcessList

type ProcessList struct {
	ProcessMap map[msg.ProcessId]msg.ProcessInterface //process id to interface
}

func HypervisorInitProcessList() {
	println("process_list.HypervisorInitProcessList()")
	ProcessListGlobal.ProcessMap = make(map[msg.ProcessId]msg.ProcessInterface)
}

func HypervisorProcessListTeardown() {
	println("process_list.HypervisorProcessListTeardown()")
	ProcessListGlobal.ProcessMap = nil
}

func AddProcess(p msg.ProcessInterface) msg.ProcessId {
	println("process_list.AddProcess()")

	id := p.GetId()
	//TODO: check to make sure processId is not already in list
	ProcessListGlobal.ProcessMap[id] = p
	return id
}

func GetProcessEvents() {
	println("process_list.GetProcessEvents()")
	//TODO
}

//run the process, creating new events for hypervisor
func ProcessTick() {
	//println("process_list.ProcessTick()")

	for id, p := range ProcessListGlobal.ProcessMap {
		_ = id
		p.Tick() //only do if incoming messages
	}
}

//events from process to hypervisor
func DispatchProcessEvents() {
	//println("process_list.DispatchProcessEvents()")

	for id, task := range ProcessListGlobal.ProcessMap {
		c := task.GetOutgoingChannel()

		for len(c) > 0 {
			m := <-c //read event
			HandleEvent(m, id)
		}
	}
}

func HandleEvent(msg []byte, Id msg.ProcessId) {
	println("process_list.HandleEvent()               ---------------- TODO !!!!!!!!!!!")
}

//Test by adding example Process
func AddTestProcess() {
	println("process_list.AddTestProcess()")
	var p *example.Process = example.NewProcess()
	var pi msg.ProcessInterface = msg.ProcessInterface(p)
	AddProcess(pi)
}
