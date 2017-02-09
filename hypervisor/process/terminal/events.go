package process

import (
	"fmt"
	"github.com/corpusc/viscript/msg"
)

func (self *State) UnpackInputEvents(msgType uint16, message []byte) []byte {
	println("hypervisor/process/terminal/events.UnpackInputEvents()")

	switch msgType {

	case msg.TypeMousePos:
		var m msg.MessageMousePos
		msg.MustDeserialize(message, &m)
		onMouseCursorPos(m)

	case msg.TypeMouseScroll:
		var m msg.MessageMouseScroll
		msg.MustDeserialize(message, &m)
		onMouseScroll(m)

	case msg.TypeMouseButton:
		var m msg.MessageMouseButton
		msg.MustDeserialize(message, &m)
		onMouseButton(m)

	case msg.TypeChar:
		var m msg.MessageChar
		msg.MustDeserialize(message, &m)
		onChar(m)
		//hypervisor.DbusGlobal.PublishTo(self.proc.MessageOut, message)
		// TODO...messageout needs to be set in SetupTerminalDbus()?

	case msg.TypeKey:
		var m msg.MessageKey
		msg.MustDeserialize(message, &m)
		onKey(m)

	case msg.TypeFrameBufferSize:
		// FIXME: BRAD SAYS THIS IS NOT INPUT
		var m msg.MessageFrameBufferSize
		msg.MustDeserialize(message, &m)
		onFrameBufferSize(m)
	default:
		fmt.Println("UNKNOWN MESSAGE TYPE!")
	}

	if self.DebugPrintInputEvents {
		fmt.Println()
	}

	return message
}

//
//EVENT HANDLERS
//

// triggered both by moving **AND*** by pressing buttons
func onMouseCursorPos(m msg.MessageMousePos) {
	println("hypervisor/process/terminal/events.onMouseCursorPos()")
}

func onMouseScroll(m msg.MessageMouseScroll) {
	println("hypervisor/process/terminal/events.onMouseScroll()")
}

func onFrameBufferSize(m msg.MessageFrameBufferSize) {
	println("hypervisor/process/terminal/events.onFrameBufferSize()")
}

func onChar(m msg.MessageChar) {
	println("hypervisor/process/terminal/events.onChar()")
}

func onKey(m msg.MessageKey) {
	println("hypervisor/process/terminal/events.onKey()")
}

func onMouseButton(m msg.MessageMouseButton) {
	println("hypervisor/process/terminal/events.onMouseButton()")
}
