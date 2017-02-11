package process

import (
	"github.com/corpusc/viscript/hypervisor"
	"github.com/corpusc/viscript/hypervisor/dbus"
	"github.com/corpusc/viscript/msg"
	//"log"
)

type State struct {
	proc                  *Process
	DebugPrintInputEvents bool
}

func (self *State) Init(proc *Process) {
	println("(process/terminal/state.go).Init()")
	self.proc = proc
	self.DebugPrintInputEvents = true
}

func (self *State) HandleMessages() {
	//println("(process/terminal/state.go).HandleMessages()")
	c := self.proc.MessageIn

	for len(c) > 0 {
		m := <-c
		msgType := msg.GetType(m[4:])
		// println("msgType", msgType)
		msgTypeMask := msgType & 0xff00
		// println("msgTypeMask", msgTypeMask)

		switch msgTypeMask {
		case msg.CATEGORY_Input:
			println("-----------CATEGORY_Input")
			self.UnpackEvent(msgType, m[6:])                                               //m
			hypervisor.DbusGlobal.PublishTo(dbus.ChannelId(self.proc.OutChannelId), m[4:]) //m
		case msg.CATEGORY_Terminal:
			println("-----------CATEGORY_Terminal     ----------NOTHING HANDLED HERE ATM")
		default:
			println("**************** UNHANDLED MESSAGE TYPE CATEGORY! ****************")
		}
	}
}
