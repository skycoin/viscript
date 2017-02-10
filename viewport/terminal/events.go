package terminal

import (
	"fmt"
	"github.com/corpusc/viscript/msg"
)

func (t *Terminal) UnpackEvent(message []byte) []byte {
	println("viewport/terminal/events.UnpackEvents()")

	switch msg.GetType(message[4:]) { // look past the channel prefix

	case msg.TypeChar:
		var m msg.MessageChar
		msg.MustDeserialize(message, &m)
		onChar(m)

	case msg.TypeKey:
		var m msg.MessageKey
		msg.MustDeserialize(message, &m)
		onKey(m)

	case msg.TypeFrameBufferSize: //FIXME? SHOULD WE HANDLE THIS MESSAGE HERE???
		// FIXME: BRAD SAYS THIS IS NOT INPUT
		var m msg.MessageFrameBufferSize
		msg.MustDeserialize(message, &m)
		onFrameBufferSize(m)
	default:
		fmt.Println("**************** UNHANDLED MESSAGE TYPE! ****************")
	}

	return message
}

//
//EVENT HANDLERS
//

//FIXME? SHOULD WE HANDLE THIS MESSAGE HERE???
func onFrameBufferSize(m msg.MessageFrameBufferSize) {
	println("viewport/terminal/events.onFrameBufferSize()")
}

func onChar(m msg.MessageChar) {
	println("viewport/terminal/events.onChar()")
}

func onKey(m msg.MessageKey) {
	println("viewport/terminal/events.onKey()")
}
