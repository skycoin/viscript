package hypervisor

/*
	Hypervisor
	- routes messages
	- maintains resource lists
	-- tasks
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

import (
	"github.com/skycoin/viscript/hypervisor/dbus"
)

var DbusGlobal dbus.DbusInstance

func Init() {
	println("<hypervisor>.Init()")
	initTaskList()
	initRunningExternalApps()
	DbusGlobal.Init()
}

func Teardown() {
	println("<hypervisor>.Teardown()")
	teardownTaskList()
	teardownRunningExternalApps()
	DbusGlobal.PubsubChannels = nil
	DbusGlobal.Resources = nil
}
