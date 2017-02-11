package process

import (
	"github.com/corpusc/viscript/hypervisor"
	"github.com/corpusc/viscript/hypervisor/dbus"
	"github.com/corpusc/viscript/msg"
	//"log"
)

//put all your process state here
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
	var msgType uint16
	var msgTypeMask uint16

	for len(c) > 0 {
		println("(process/terminal/state.go).HandleMessages() - ...ONE ITERATION OF CHANNEL")
		m := <-c // read from channel
		//route the message
		msgType = msg.GetType(m[4:])
		// println("msgType", msgType)
		msgTypeMask = msgType & 0xff00
		// println("msgTypeMask", msgTypeMask)

		switch msgTypeMask {
		case msg.CATEGORY_Input:
			println("-----------CATEGORY_Input")
			self.UnpackInputEvents(msgType, m[6:])                                         //m
			hypervisor.DbusGlobal.PublishTo(dbus.ChannelId(self.proc.OutChannelId), m[4:]) //m
		case msg.CATEGORY_Terminal:
			println("-----------CATEGORY_Terminal     ----------NOTHING HANDLED HERE ATM")
		default:
			println("**************** UNHANDLED MESSAGE TYPE CATEGORY! ****************")
		}
	}
}
