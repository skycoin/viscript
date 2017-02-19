package terminal

import (
	"github.com/corpusc/viscript/msg"
)

func (t *Terminal) UnpackEvent(message []byte) []byte {
	println("viewport/terminal/events.UnpackEvent()")

	//TODO/FIXME:   cache channel id wherever it may be needed
	message = message[4:] //for now DISCARD the channel id prefix

	switch msg.GetType(message) {

	case msg.TypePutChar:
		println("viewport/terminal/events <<< msg.TypePutChar >>>")
		var m msg.MessagePutChar
		msg.MustDeserialize(message, &m)
		t.PutCharacter(m.Char)

	case msg.TypeFrameBufferSize:
		println("viewport/terminal/events ------- case msg.TypeFrameBufferSize --------")
		var m msg.MessageFrameBufferSize
		msg.MustDeserialize(message, &m)
		t.onFrameBufferSize(m)

	default:
		println("viewport/terminal/events.go ************ UNHANDLED MESSAGE TYPE! ************")
	}

	return message
}

//
//EVENT HANDLERS
//

//FIXME? PROBABLY SHOULDN'T PASS THIS AROUND DBUS
//(i think it gets consumed and never passed on, probably in viewport)
//if we ever want to handle this message deeper, it should probably
//only go as far as the terminal_stack?  this is only relevant
//to local screen sizes, so should never be sent over a network.
//and task logic shouldn't care about the display.
//if we want to do anything here, then we should probably just send msg directly
func (t *Terminal) onFrameBufferSize(m msg.MessageFrameBufferSize) {
	println("viewport/terminal/events.onFrameBufferSize()", m.X, m.Y)
	println("......NOT DOING ANYTHING HERE, READ THE COMMENTS HERE")
	// gl.ResizeViewportToFrameBuffer(int32(m.X), int32(m.Y))
	// we DEFINITELY don't want to import/use anything from gl
}
