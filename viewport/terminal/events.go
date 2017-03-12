package terminal

import (
	"fmt"
	"github.com/corpusc/viscript/app"
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

	//lower level messages
	case msg.TypeKey:
		var m msg.MessageKey
		msg.MustDeserialize(message, &m)
		t.onKey(m)

	case msg.TypeMouseScroll:
		var m msg.MessageMouseScroll
		msg.MustDeserialize(message, &m)
		t.onMouseScroll(m)

	default:
		println("viewport/terminal/events.go ************ UNHANDLED MESSAGE TYPE! ************")
	}

	return message
}

//
//EVENT HANDLERS
//

func (t *Terminal) onKey(m msg.MessageKey) {
	println("viewport/terminal/events.onKey()")

	switch m.Key {
	case msg.KeyEnter:
		t.NewLine()
	}
}

func (t *Terminal) onMouseScroll(m msg.MessageMouseScroll) {
	println("viewport/terminal/events.onMouseScroll()")
	fmt.Printf("mouse scroll in Y: %.3f\n", m.Y)
	fmt.Printf("mouse scroll in Y: %.3f\n", m.Y)
	fmt.Printf("mouse scroll in Y: %.3f\n", m.Y)
	fmt.Printf("mouse scroll in Y: %.3f\n", m.Y)

	if m.HoldingControl {
		//only using m.Y because
		//m.X is sideways scrolling (which most(?) mice can't do)
		y := float32(m.Y)
		changeFactor := float32(1 + app.Clamp(y, -1, 1)/10)
		t.CharSize.X *= changeFactor
		t.CharSize.Y *= changeFactor
		t.Bounds.Right = t.Bounds.Left + t.Bounds.Width()*changeFactor
		t.Bounds.Bottom = t.Bounds.Top - t.Bounds.Height()*changeFactor
	}
}
