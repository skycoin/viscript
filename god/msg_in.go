package god

import (
	"github.com/corpusc/viscript/god/gl"
	"github.com/corpusc/viscript/msg"
	_ "strconv"
)

var DebugPrintInputEvents = true

func UnpackMessage(message []byte) []byte {
	switch msg.GetType(message) {

	case msg.TypeMousePos:
		var m msg.MessageMousePos
		msg.MustDeserialize(message, &m)
		onMouseCursorPos(m)

	case msg.TypeMouseScroll:
		var m msg.MessageMouseScroll
		msg.MustDeserialize(message, &m)
		onMouseScroll(m)

		if Terms.Focused != nil {
			Terms.Focused.RelayToTask(message)
		}

	case msg.TypeMouseButton:
		var m msg.MessageMouseButton
		msg.MustDeserialize(message, &m)
		onMouseButton(m)

	case msg.TypeChar:
		var m msg.MessageChar
		msg.MustDeserialize(message, &m)
		onChar(m)

		if Terms.Focused != nil {
			Terms.Focused.RelayToTask(message)
		}

	case msg.TypeKey:
		var m msg.MessageKey
		msg.MustDeserialize(message, &m)
		onKey(m)

		if Terms.Focused != nil {
			Terms.Focused.RelayToTask(message)
		}

	case msg.TypeFrameBufferSize:
		var m msg.MessageFrameBufferSize
		msg.MustDeserialize(message, &m)
		onFrameBufferSize(m)

	default:
		app.At("god/msg_in", "************ UNHANDLED MESSAGE TYPE! ************")
	}

	return message
}

func onFrameBufferSize(m msg.MessageFrameBufferSize) {
	if DebugPrintInputEvents {
		print("TypeFrameBufferSize")
		showUInt32("X", m.X)
		showUInt32("Y", m.Y)
		println()
	}

	gl.SetSize(int32(m.X), int32(m.Y))
}
