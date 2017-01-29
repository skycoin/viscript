package hypervisor

import (
//"fmt"
)

/*
	Hypervisor
	- routes messages
	- maintains resource lists
	-- processes
	-- network connections?
	-- file system access?

	- routes messages between resources?

	Resouce Type
	Resouce Id
*/

/*
types of messages
- one to one channels, with resource locks
- emits messages
- receives messages
- many to one, pubsub (publication to all subscribers)

- many to one, pubsub (publication to all subscribers)
-- list of people who will receive

- receive message without ACK

- RPC, message with guarnteed return value

- only "owner" can write channel
- anyone can write channel

Can objects create multiple channels?
*/

func HypervisorInit() {
	println("hypervisor.go --- HypervisorInit()")
	HypervisorInitProcessList()
	DbusInit()
	AddTestProcess()
}

func HypervisorTeardown() {
	HypervisorProcessListTeardown()
	DbusTeardown()
}
