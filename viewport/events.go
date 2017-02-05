package viewport

import (
	"fmt"
	"github.com/corpusc/viscript/msg"
	_ "strconv"
)

var DebugPrintInputEvents = false

func UnpackInputEvents(message []byte) []byte {
	switch msg.GetType(message) {

	case msg.TypeMousePos:
		var msgMousePos msg.MessageMousePos
		msg.MustDeserialize(message, &msgMousePos)

		if DebugPrintInputEvents {
			fmt.Print("TypeMousePos")
			showFloat64("X", msgMousePos.X)
			showFloat64("Y", msgMousePos.Y)
		}

		onMouseCursorPos(msgMousePos)

	case msg.TypeMouseScroll:
		var m msg.MessageMouseScroll
		msg.MustDeserialize(message, &m)
		onMouseScroll(m)

	case msg.TypeMouseButton:
		var msgBtn msg.MessageMouseButton
		msg.MustDeserialize(message, &msgBtn)

		if DebugPrintInputEvents {
			fmt.Print("TypeMouseButton")
			showUInt8("Button", msgBtn.Button)
			showUInt8("Action", msgBtn.Action)
			showUInt8("Mod", msgBtn.Mod)
		}

		onMouseButton(msgBtn)

	case msg.TypeChar:
		var msgChar msg.MessageChar
		msg.MustDeserialize(message, &msgChar)

		if DebugPrintInputEvents {
			fmt.Print("TypeChar")
		}

		onChar(msgChar)
		Terms.Focused.RelayToTask(message)

	case msg.TypeKey:
		var keyMsg msg.MessageKey
		msg.MustDeserialize(message, &keyMsg)

		//if DebugPrintInputEvents {
		fmt.Print("TypeKey")
		showUInt32("Key", keyMsg.Key)
		showUInt32("Scan", keyMsg.Scan)
		showUInt8("Action", keyMsg.Action)
		showUInt8("Mod", keyMsg.Mod)
		//}

		onKey(keyMsg)
		Terms.Focused.RelayToTask(message)

	case msg.TypeFrameBufferSize:
		// FIXME: BRAD SAYS THIS IS NOT INPUT
		var m msg.MessageFrameBufferSize
		msg.MustDeserialize(message, &m)

		//if DebugPrintInputEvents {
		fmt.Print("TypeFrameBufferSize")
		showUInt32("X", m.X)
		showUInt32("Y", m.Y)
		//}

		onFrameBufferSize(m)
	default:
		fmt.Println("UNKNOWN MESSAGE TYPE!")
	}

	if DebugPrintInputEvents {
		fmt.Println()
	}

	return message
}
