package terminal

import (
	"github.com/corpusc/viscript/msg"
)

func (t *Terminal) UnpackEvent(message []byte) []byte {
	println("viewport/terminal/events.UnpackEvent()")

	//TODO/FIXME:   cache channel id wherever it may be needed
	message = message[4:] //for now DISCARD the channel id prefix

	switch msg.GetType(message) {

	case msg.TypeCommandLine:
		var m msg.MessageCommandLine
		msg.MustDeserialize(message, &m)
		t.updateCommandLine(m)

	case msg.TypePutChar:
		var m msg.MessagePutChar
		msg.MustDeserialize(message, &m)
		t.PutCharacter(m.Char)

	case msg.TypeKey:
		var m msg.MessageKey
		msg.MustDeserialize(message, &m)
		t.onKey(m)

	// case msg.TypeFrameBufferSize:
	// 	println("viewport/terminal/events ------- case msg.TypeFrameBufferSize --------")
	// 	var m msg.MessageFrameBufferSize
	// 	msg.MustDeserialize(message, &m)
	// 	t.onFrameBufferSize(m)

	default:
		println("viewport/terminal/events.go ************ UNHANDLED MESSAGE TYPE! ************")
	}

	return message
}

//
//EVENT HANDLERS
//

//FIXME? PROBABLY SHOULDN'T PASS THIS AROUND DBUS
//if we ever want to handle this message deeper, it should probably
//only go as far as the terminal_stack?  this is only relevant
//to local screen sizes, so should never be sent over a network.
//and task logic shouldn't care about the display.
//if we want to do anything here, then we should probably just
//call this func directly with no messaging?
//
// func (t *Terminal) onFrameBufferSize(m msg.MessageFrameBufferSize) {
// 	println("viewport/terminal/events.onFrameBufferSize()", m.X, m.Y)
// 	println("......NOT DOING ANYTHING HERE, READ THE COMMENTS HERE")
// 	// we DEFINITELY don't want to import/use anything from gl
// }

func (t *Terminal) onKey(m msg.MessageKey) {
	println("viewport/terminal/events.onKey() -----SHOULD NOT HANDLE THESE HERE?")

	switch m.Key {
	case msg.KeyEnter:
		t.MoveDown()
		// 	t.Curr.X = 0
		// case msg.KeyBackspace:
		// 	t.BackSpace()
	}
}
