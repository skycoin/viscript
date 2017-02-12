package terminal

import (
	"fmt"
	"math/rand"

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

	//int / character grid space
	Curs     app.Vec2UI32 //current cursor/insert pos
	GridSize app.Vec2I    //number of characters
	Chars    [NumRows][NumColumns]uint32

	//float32 / GL space
	//(mouse pos events, and resize event resolutions are the only things that use pixels)
	Bounds *app.Rectangle
	Depth  float32 //0 for lowest
}

func (t *Terminal) Init() {
	fmt.Printf("Terminal.Init()\n")

	t.TerminalId = msg.RandTerminalId()
	t.InChannel = make(chan []byte, msg.ChannelCapacity)
	t.GridSize = app.Vec2I{NumColumns, NumRows}
	t.Chars = [NumRows][NumColumns]uint32{}
	t.makeRandomChars(20)
	t.PutString("calling PutString() with a very loooooooooooooooooooooooooooooooooooong string")
	t.SetCursor(0, 0)
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
	t.Curs.X--

	if t.Curs.X < 0 {
		t.Curs.X = uint32(t.GridSize.X) - 1
		t.MoveUp()
	}
}

func (t *Terminal) MoveRight() {
	t.Curs.X++

	if t.Curs.X >= uint32(t.GridSize.X) {
		t.Curs.X = 0
		t.MoveDown()
	}
}

func (t *Terminal) MoveUp() {
	t.Curs.Y--

	if t.Curs.Y < 0 {
		t.Curs.Y = uint32(t.GridSize.Y) - 1
	}
}

func (t *Terminal) MoveDown() {
	t.Curs.Y++

	if t.Curs.Y >= uint32(t.GridSize.Y) {
		t.Curs.Y = 0
	}
}

func (t *Terminal) SpanX() float32 { // span across a char
	return t.Bounds.Width() / float32(t.GridSize.X)
}

func (t *Terminal) SpanY() float32 {
	return t.Bounds.Height() / float32(t.GridSize.Y)
}

func (t *Terminal) SetCursor(x, y uint32) {
	if t.posIsValid(x, y) {
		t.Curs.X = x
		t.Curs.Y = y
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
	if t.posIsValid(t.Curs.X, t.Curs.Y) {
		t.SetCharacterAt(t.Curs.X, t.Curs.Y, char)
		t.MoveRight()
	}
}

func (t *Terminal) SetCharacterAt(x, y uint32, Char uint32) {
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

func (t *Terminal) SetStringAt(X, Y uint32, S string) {
	numOOB = 0

	for x, c := range S {
		if t.posIsValid(X+uint32(x), Y) {
			t.Chars[Y][X+uint32(x)] = uint32(c)
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
func (t *Terminal) makeRandomChars(count int) {
	for i := 0; i < count; i++ {
		t.SetCharacterAt(
			uint32(rand.Int31n(NumColumns)),
			uint32(rand.Int31n(NumRows)),
			uint32(rand.Int31n(128)))
	}
}

var numOOB int // number of out of bound characters
func (t *Terminal) posIsValid(X, Y uint32) bool {
	if X < 0 || X >= uint32(t.GridSize.X) ||
		Y < 0 || Y >= uint32(t.GridSize.Y) {
		numOOB++

		if numOOB == 1 {
			println("****** ATTEMPTED OUT OF BOUNDS CHARACTER PLACEMENT! ******")
		}

		return false
	}

	return true
}
