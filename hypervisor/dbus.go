package hypervisor

import (
	"github.com/corpusc/viscript/hypervisor/dbus"
)

var DbusGlobal dbus.DbusInstance

func DbusInit() {
	println("(hypervisor/dbus.go).DbusInit()")
	DbusGlobal.Init()
}

func DbusTeardown() {
	println("(hypervisor/dbus.go).DbusTeardown() ---- TO DOne?")
	DbusGlobal.PubsubChannels = nil
	DbusGlobal.Resources = nil
}
