package process

import (
	"fmt"
	"github.com/corpusc/viscript/msg"
)

func (self *State) UnpackInputEvents(msgType uint16, message []byte) []byte {

	switch msgType {

	case msg.TypeMousePos:
		var msgMousePos msg.MessageMousePos
		msg.MustDeserialize(message, &msgMousePos)

		if self.DebugPrintInputEvents {
			fmt.Print("TypeMousePos")
			showFloat64("X", msgMousePos.X)
			showFloat64("Y", msgMousePos.Y)
		}

		onMouseCursorPos(msgMousePos)

	case msg.TypeMouseScroll:
		var msgScroll msg.MessageMouseScroll
		msg.MustDeserialize(message, &msgScroll)

		if self.DebugPrintInputEvents {
			fmt.Print("TypeMouseScroll")
			showFloat64("X Offset", msgScroll.X)
			showFloat64("Y Offset", msgScroll.Y)
		}

		onMouseScroll(msgScroll)

	case msg.TypeMouseButton:
		var msgBtn msg.MessageMouseButton
		msg.MustDeserialize(message, &msgBtn)

		if self.DebugPrintInputEvents {
			fmt.Print("TypeMouseButton")
			showUInt8("Button", msgBtn.Button)
			showUInt8("Action", msgBtn.Action)
			showUInt8("Mod", msgBtn.Mod)
		}

		onMouseButton(msgBtn)

	case msg.TypeChar:
		var msgChar msg.MessageChar
		msg.MustDeserialize(message, &msgChar)

		if self.DebugPrintInputEvents {
			fmt.Print("TypeChar")
		}

		onChar(msgChar)

	case msg.TypeKey:
		var keyMsg msg.MessageKey
		msg.MustDeserialize(message, &keyMsg)

		if self.DebugPrintInputEvents {
			fmt.Print("TypeKey")
			showUInt32("Key", keyMsg.Key)
			showUInt32("Scan", keyMsg.Scan)
			showUInt8("Action", keyMsg.Action)
			showUInt8("Mod", keyMsg.Mod)
		}

		onKey(keyMsg)

	case msg.TypeFrameBufferSize:
		// FIXME: BRAD SAYS THIS IS NOT INPUT
		var m msg.MessageFrameBufferSize
		msg.MustDeserialize(message, &m)

		if self.DebugPrintInputEvents {
			fmt.Print("TypeFrameBufferSize")
			showUInt32("X", m.X)
			showUInt32("Y", m.Y)
		}

		onFrameBufferSize(m)
	default:
		fmt.Println("UNKNOWN MESSAGE TYPE!")
	}

	if self.DebugPrintInputEvents {
		fmt.Println()
	}

	return message
}

func showUInt8(s string, x uint8) uint8 {
	fmt.Printf("   [%s: %d]", s, x)
	return x
}

// the rest of these funcs are almost identical, just top 2 vars customized (and string format)
func showSInt32(s string, x int32) {
	fmt.Printf("   [%s: %d]", s, x)
}

func showUInt32(s string, x uint32) {
	fmt.Printf("   [%s: %d]", s, x)
}

func showFloat64(s string, f float64) float64 {
	fmt.Printf("   [%s: %.1f]", s, f)
	return f
}

//
//EVENT HANDLERS
//

// triggered both by moving **AND*** by pressing buttons
func onMouseCursorPos(m msg.MessageMousePos) {
	println("hypervisor/process/terminal/events.onMouseCursorPos(m msg.Message)")
}

func onMouseScroll(m msg.MessageMouseScroll) {
	println("hypervisor/process/terminal/events.onMouseScroll(m msg.Message)")
	/*
		var delta float32 = 30

		if eitherControlKeyHeld() { // horizontal ability from 1D scrolling
			ScrollTermThatHasMousePointer(float32(m.Y)*-delta, 0)
		} else { // can handle both x & y for 2D scrolling
			ScrollTermThatHasMousePointer(float32(m.X)*delta, float32(m.Y)*-delta)
		}
	*/
}

func onFrameBufferSize(m msg.MessageFrameBufferSize) {
	println("hypervisor/process/terminal/events.onFrameBufferSize(m msg.Message)")
}

func onChar(m msg.MessageChar) {
	println("hypervisor/process/terminal/events.onChar(m msg.Message)")
	//InsertRuneIntoDocument("Rune", m.Char)
	//script.Process(false)
}

func onKey(m msg.MessageKey) {
	println("hypervisor/process/terminal/events.onKey(m msg.Message)")
}

func onMouseButton(m msg.MessageMouseButton) {
	println("hypervisor/process/terminal/events.onMouseButton(m msg.Message)")
}
