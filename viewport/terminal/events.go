package terminal

import (
	"fmt"
	"github.com/corpusc/viscript/msg"
)

func (t *Terminal) UnpackEvent(message []byte) []byte {
	println("viewport/terminal/events.UnpackEvent()")

	//TODO/FIXME:   cache channel id wherever it may be needed
	message = message[4:] //for now DISCARD the channel id prefix

	switch msg.GetType(message) {

	case msg.TypeChar:
		var m msg.MessageChar
		msg.MustDeserialize(message, &m)
		t.onChar(m)

	case msg.TypeKey:
		var m msg.MessageKey
		msg.MustDeserialize(message, &m)
		t.onKey(m)

	case msg.TypeFrameBufferSize: //FIXME? SHOULD WE HANDLE THIS MESSAGE HERE???
		// FIXME: BRAD SAYS THIS IS NOT INPUT
		var m msg.MessageFrameBufferSize
		msg.MustDeserialize(message, &m)
		t.onFrameBufferSize(m)

	default:
		fmt.Println("viewport/terminal/events.go ************* UNHANDLED MESSAGE TYPE! *************")
	}

	return message
}

//
//EVENT HANDLERS
//

//FIXME? SHOULD WE HANDLE THIS MESSAGE HERE???
func (t *Terminal) onFrameBufferSize(m msg.MessageFrameBufferSize) {
	println("viewport/terminal/events.onFrameBufferSize()")
}

func (t *Terminal) onChar(m msg.MessageChar) {
	println("viewport/terminal/events.onChar()")
	t.PutCharacter(m) // TEMPORARY hack
}

func (t *Terminal) onKey(m msg.MessageKey) {
	println("viewport/terminal/events.onKey()")
}
