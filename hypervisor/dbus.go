package hypervisor

import (
	"github.com/corpusc/viscript/app"
	"github.com/corpusc/viscript/hypervisor/dbus"
)

var path = "hypervisor/dbus"
var DbusGlobal dbus.DbusInstance

func DbusInit() {
	app.At(path, "DbusInit")
	DbusGlobal.Init()
}

func DbusTeardown() {
	app.At(path, "DbusTeardown")
	DbusGlobal.PubsubChannels = nil
	DbusGlobal.Resources = nil
}
