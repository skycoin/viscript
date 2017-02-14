package terminal

import (
	"github.com/corpusc/viscript/msg"
	"github.com/corpusc/viscript/viewport/gl"
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

	case msg.TypeFrameBufferSize: //FIXME? SHOULD WE HANDLE THIS MESSAGE HERE???
		//(i think it gets consumed and never passed on, probably in viewport)
		println("viewport/terminal/events ------- case msg.TypeFrameBufferSize -------- TODO?")
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

//FIXME? SHOULD WE HANDLE THIS MESSAGE HERE???
func (t *Terminal) onFrameBufferSize(m msg.MessageFrameBufferSize) {
	println("viewport/terminal/events.onFrameBufferSize()", m.X, m.Y)

	// gl.ResizeViewportToFrameBuffer(int32(m.X), int32(m.Y))
}
