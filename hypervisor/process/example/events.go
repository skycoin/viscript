package process

import (
	"fmt"
	"github.com/corpusc/viscript/msg"
)

func (self *State) UnpackInputEvents(msgType uint16, message []byte) []byte {

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

func showUInt8(s string, x uint8) {
	fmt.Printf("   [%s: %d]", s, x)
}

// the rest of these funcs are almost identical, just top 2 vars customized (and string format)
func showSInt32(s string, x int32) {
	fmt.Printf("   [%s: %d]", s, x)
}

func showUInt32(s string, x uint32) {
	fmt.Printf("   [%s: %d]", s, x)
}

func showFloat64(s string, f float64) {
	fmt.Printf("   [%s: %.1f]", s, f)
}

//
//EVENT HANDLERS
//

// triggered both by moving **AND*** by pressing buttons
func onMouseCursorPos(m msg.MessageMousePos) {
	println("hypervisor/process/example/events.onMouseCursorPos()")
}

func onMouseScroll(m msg.MessageMouseScroll) {
	println("hypervisor/process/example/events.onMouseScroll()")
}

func onFrameBufferSize(m msg.MessageFrameBufferSize) {
	println("hypervisor/process/example/events.onFrameBufferSize()")
}

func onChar(m msg.MessageChar) {
	println("hypervisor/process/example/events.onChar()")
}

func onKey(m msg.MessageKey) {
	println("(hypervisor/process/example/events.go).onKey()")
}

func onMouseButton(m msg.MessageMouseButton) {
	println("(hypervisor/process/example/events.go).onMouseButton()")
}
