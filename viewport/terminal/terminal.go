package terminal

import (
	"fmt"

	"github.com/corpusc/viscript/app"
	"github.com/corpusc/viscript/hypervisor"
	"github.com/corpusc/viscript/hypervisor/dbus"
	"github.com/corpusc/viscript/msg"
)

const (
	// num == count/number of...
	NumColumns = 64
	NumRows    = 32
)

type Terminal struct {
	TerminalId      msg.TerminalId
	AttachedProcess msg.ProcessId
	OutChannelId    dbus.ChannelId //id of pubsub channel
	InChannel       chan []byte

	//int/character grid space
	Curr     app.Vec2I //current insert position
	Cursor   app.Vec2I
	GridSize app.Vec2I //number of characters across
	Chars    [NumRows][NumColumns]uint32

	//float/GL space
	//(mouse pos events & frame buffer sizes are the only things that use pixels)
	BorderSize float32
	CharSize   app.Vec2F
	Bounds     *app.Rectangle
	Depth      float32 //0 for lowest

	ResizingRight  bool
	ResizingBottom bool
}

func (t *Terminal) Init() {
	println("Terminal.Init()")

	t.TerminalId = msg.RandTerminalId()
	t.InChannel = make(chan []byte, msg.ChannelCapacity)
	t.GridSize = app.Vec2I{NumColumns, NumRows}
	t.Chars = [NumRows][NumColumns]uint32{}
	t.BorderSize = 0.013
	t.SetSize()

	t.PutString(">")
	t.SetCursor(1, 0)
	t.ResizingRight = false
	t.ResizingBottom = false
}

func (t *Terminal) IsResizing() bool {
	return t.ResizingRight || t.ResizingBottom
}

func (t *Terminal) SetResizingOff() {
	t.ResizingRight = false
	t.ResizingBottom = false
}

func (t *Terminal) SetSize() {
	println("Terminal.SetSize()        --------FIXME once we allow dragging edges")
	t.CharSize.X = (t.Bounds.Width() - t.BorderSize*2) / float32(t.GridSize.X)
	t.CharSize.Y = (t.Bounds.Height() - t.BorderSize*2) / float32(t.GridSize.Y)
}

func (t *Terminal) BackSpace() {
	// FIXME: should have to look at this more in depth tomorrow
	t.MoveLeft()
	t.PutCharacter(0)
	t.MoveLeft()
}

func (t *Terminal) Tick() {
	for len(t.InChannel) > 0 {
		t.UnpackEvent(<-t.InChannel)
	}
}

func (t *Terminal) Clear() {
	for y := 0; y < t.GridSize.Y; y++ {
		for x := 0; x < t.GridSize.X; x++ {
			t.Chars[y][x] = 0
		}
	}
}

func (t *Terminal) RelayToTask(message []byte) {
	hypervisor.DbusGlobal.PublishTo(t.OutChannelId, message)
}

func (t *Terminal) MoveLeft() {
	t.Curr.X--

	if t.Curr.X < 0 {
		t.Curr.X = t.GridSize.X - 1
		t.MoveUp()
	}
}

func (t *Terminal) MoveRight() {
	t.Curr.X++

	if t.Curr.X >= t.GridSize.X {
		t.Curr.X = 0
		t.MoveDown()
	}
}

func (t *Terminal) MoveUp() {
	t.Curr.Y--

	if t.Curr.Y < 0 {
		t.Curr.Y = t.GridSize.Y - 1
	}
}

func (t *Terminal) MoveDown() {
	t.Curr.Y++

	if t.Curr.Y >= t.GridSize.Y {
		t.Curr.Y = 0
	}
}

func (t *Terminal) SetCursor(x, y int) {
	if t.posIsValid(x, y) {
		t.Cursor.X = x
		t.Cursor.Y = y
	}
}

// there should be 2 paradigms of adding chars/strings:
//
// (1) full manual control/management.  (explicitly tell terminal exactly
//			where to place something, without disrupting cursor position.
//			must make sure there is space for it)
// (2) automated flow control.  (just tell what char/string to put into the current flow
//			and term manages it's placement, wrapping, & eventually word-preserving-wrapping)
func (t *Terminal) PutCharacter(char uint32) {
	if t.posIsValid(t.Curr.X, t.Curr.Y) {
		t.SetCharacterAt(t.Curr.X, t.Curr.Y, char)
		t.MoveRight()
	}
}

func (t *Terminal) SetCharacterAt(x, y int, Char uint32) {
	numOOB = 0

	if t.posIsValid(x, y) {
		t.Chars[y][x] = Char
	}
}

func (t *Terminal) PutString(s string) {
	for _, c := range s {
		t.PutCharacter(uint32(c))
	}
}

func (t *Terminal) SetStringAt(X, Y int, S string) {
	numOOB = 0

	for x, c := range S {
		if t.posIsValid(X+x, Y) {
			t.Chars[Y][X+x] = uint32(c)
		}
	}
}

func (t *Terminal) SetGridSize() {
	fmt.Printf("Terminal.SetGridSize()\n")

	// ERROR IN THIS CODE!
	// .... but I don't see any need for this right now anyways.

	// t.Chars = make([][]uint32, t.GridSize.Y, t.GridSize.Y)

	// for i := range(0, t.GridSize.X):
	// 	t.Chars[i] = make([]uint32, t.GridSize.X, t.GridSize.X)
}

// private
func (t *Terminal) updateCommandLine(m msg.MessageCommandLine) {
	for i := 0; i < t.GridSize.X*2; i++ {
		var char uint32
		x := i % t.GridSize.X
		y := i / t.GridSize.X

		if i == int(m.CursorOffset) {
			t.Cursor.X = x
			t.Cursor.Y = y
		}

		y += int(t.Curr.Y)

		if i < len(m.CommandLine) {
			char = uint32(m.CommandLine[i])
		} else {
			char = 0
		}

		t.SetCharacterAt(x, y, char)
	}
}

var numOOB int // number of out of bound characters
func (t *Terminal) posIsValid(X, Y int) bool {
	if X < 0 || X >= t.GridSize.X ||
		Y < 0 || Y >= t.GridSize.Y {
		numOOB++

		if numOOB == 1 {
			println("****** ATTEMPTED OUT OF BOUNDS CHARACTER PLACEMENT! ******")
		}

		return false
	}

	return true
}
