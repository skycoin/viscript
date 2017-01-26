package hypervisor

import (
	"github.com/corpusc/viscript/hypervisor/dbus"
)

var DbusGlobal dbus.DbusInstance

func DbusInit() {
	DbusGlobal.Init()
}

func DbusTeardown() {

}
