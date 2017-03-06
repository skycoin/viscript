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

	case msg.TypeSetCharAt:
		var m msg.MessageSetCharAt
		msg.MustDeserialize(message, &m)
		t.SetCharacterAt(int(m.X), int(m.Y), m.Char)

	case msg.TypePutChar:
		var m msg.MessagePutChar
		msg.MustDeserialize(message, &m)
		t.PutCharacter(m.Char)

	case msg.TypeKey:
		var m msg.MessageKey
		msg.MustDeserialize(message, &m)
		t.onKey(m)

	default:
		println("viewport/terminal/events.go ************ UNHANDLED MESSAGE TYPE! ************")
	}

	return message
}

//
//EVENT HANDLERS
//

func (t *Terminal) onKey(m msg.MessageKey) {
	println("viewport/terminal/events.onKey() -----SHOULD NOT HANDLE THESE HERE?")

	switch m.Key {
	case msg.KeyEnter:
		t.NewLine()
	}
}
