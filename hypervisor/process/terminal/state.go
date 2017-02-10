package process

import (
	"github.com/corpusc/viscript/hypervisor"
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
	// var channelId uint16  doesn't compile unused. uncomment for usage
	var msgType uint16
	var msgTypeMask uint16

	for len(c) > 0 {
		println("(process/terminal/state.go).HandleMessages() - ...ONE ITERATION OF CHANNEL")
		m := <-c // read from channel
		//route the message
		// channelId = msg.GetType(m)	and here too.
		// println("PubSubChannelId", channelId)
		msgType = msg.GetType(m[4:])
		// println("MsgType", msgType)
		msgTypeMask = msgType & 0xff00
		// println("MsgTypeMask", msgTypeMask)

		switch msgTypeMask {
		case msg.TypePrefix_Input:
			println("-----------TypePrefix_Input")
			self.UnpackInputEvents(msgType, m[6:])
			hypervisor.DbusGlobal.PublishTo(self.proc.OutChannelId, m[4:])
		case msg.TypePrefix_Terminal: //process to hypervisor messages
			println("-----------TypePrefix_Terminal")
		default:
			println("**************** UNHANDLED TYPE PREFIX! ****************")
		}

		//FIXME
		//at this point, messages have already been filtered, in our current usage.
		//but the lines below will need to be put inside the cases once we start
		//doing any local filtering or processing of the messages

	}
}
