package terminal

import (
	"github.com/skycoin/viscript/msg"
	"github.com/skycoin/viscript/viewport/gl"
)

func (t *Terminal) move(m msg.MessageMoveTerminal) {
	//println("Terminal.move() - given x,y:", float32(m.X), float32(m.Y))
	w := t.Bounds.Width()
	h := t.Bounds.Height()

	//FIXME?  positioning in desktop space by the size of Terminal's characters
	//seems a bit wonky.  but i'm assuming the given coords should be similar
	//to 80x25, so that it's like text mode positioning
	t.Bounds.Left = -gl.CanvasExtents.X + float32(m.X)*t.CharSize.X
	t.Bounds.Top = gl.CanvasExtents.Y - float32(m.Y)*t.CharSize.Y
	t.Bounds.Top -= t.CharSize.Y + t.BorderSize

	t.Bounds.Right = t.Bounds.Left + w
	t.Bounds.Bottom = t.Bounds.Top - h
}
