package terminal

import (
	"github.com/skycoin/viscript/msg"
	"github.com/skycoin/viscript/viewport/gl"
	"math"
)

func (t *Terminal) SetCursor(x, y int) {
	if t.posIsValidElsePrint(x, y) {
		t.Cursor.X = x
		t.Cursor.Y = y
	}
}

//2 paradigms for adding chars/strings:
//
//(1) Set___At() ------ full manual control/management.
//			explicitly tell terminal exactly where to place
//			something, without disrupting current position.
//			must make sure there is space for it.
//
//(2) Put___() -------- automated flow control.
//			just tell what char/string to put into the current...
//			... flow.  then Terminal manages it's placement & wrapping
//
func (t *Terminal) SetCharacterAt(x, y int, Char uint32) {
	numOOB = 0

	if t.posIsValidElsePrint(x, y) {
		t.Chars[y][x] = Char
	}
}

func (t *Terminal) PutString(s string) {
	println(".PutString(" + s + ")")

	for _, c := range s {
		t.putCharacter(uint32(c))
	}
}

func (t *Terminal) SetStringAt(x, y int, s string) {
	numOOB = 0

	for col /* column */, c := range s {
		if t.posIsValidElsePrint(x+col, y) {
			t.Chars[y][x+col] = uint32(c)
		}
	}
}

//
//
// private
//
//

func (t *Terminal) putCharacter(char uint32) {
	if t.posIsValidElsePrint(t.CurrFlowPos.X, t.CurrFlowPos.Y) {
		t.SetCharacterAt(t.CurrFlowPos.X, t.CurrFlowPos.Y, char)
		t.MoveRight()
	}
}

func (t *Terminal) move(m msg.MessageMoveTerminal) {
	//println("Terminal.move() - given x,y:", float32(m.X), float32(m.Y))
	w := t.Bounds.Width()
	h := t.Bounds.Height()
	minimumVisibleSpan := float32(math.Min(float64(w), float64(h)))
	minimumVisibleSpan /= 10

	//FIXME?  positioning in desktop space by the size of Terminal's characters
	//seems a bit wonky.  they could be any size for any given terminal, in theory.
	//but i'm assuming the given coords should be similar
	//to 80x25, so that it's like text mode positioning, in such a grid.
	t.Bounds.Left = -gl.CanvasExtents.X + float32(m.X)*t.CharSize.X
	t.Bounds.Top = gl.CanvasExtents.Y - float32(m.Y)*t.CharSize.Y
	t.pushDownByTabHeight(&t.Bounds.Top)

	//make sure at least a corner of the Terminal remains visible in the desktop viewport
	if t.Bounds.Left < -gl.CanvasExtents.X-w+minimumVisibleSpan {
		t.Bounds.Left = -gl.CanvasExtents.X - w + minimumVisibleSpan
	}

	if t.Bounds.Left > gl.CanvasExtents.X-minimumVisibleSpan {
		t.Bounds.Left = gl.CanvasExtents.X - minimumVisibleSpan
	}

	if t.Bounds.Top > gl.CanvasExtents.Y+h-minimumVisibleSpan {
		t.Bounds.Top = gl.CanvasExtents.Y + h - minimumVisibleSpan
		t.pushDownByTabHeight(&t.Bounds.Top)
	}

	if t.Bounds.Top < -gl.CanvasExtents.Y+minimumVisibleSpan {
		t.Bounds.Top = -gl.CanvasExtents.Y + minimumVisibleSpan
		t.pushDownByTabHeight(&t.Bounds.Top)
	}

	//set the bottom right corner, now that upper left has been clamped to valid space
	t.Bounds.Right = t.Bounds.Left + w
	t.Bounds.Bottom = t.Bounds.Top - h
}

func (t *Terminal) pushDownByTabHeight(f *float32) {
	*f -= t.CharSize.Y + t.BorderSize
}
