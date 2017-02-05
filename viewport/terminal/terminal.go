package terminal

import (
	"fmt"
	"math/rand"

	"github.com/corpusc/viscript/app"
	//"github.com/corpusc/viscript/hypervisor"
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
	OutChannelId    uint32 //id of pubsub channel
	InChannel       chan []byte

	//vars for character grid (of cells)
	Chars    [NumRows][NumColumns]uint32
	Curs     app.Vec2UI32 //current cursor/insert pos
	GridSize app.Vec2I    //number of characters

	//vars for GL space / float
	//(mouse pos events, and resize event resolutions are the only things that use pixels)
	Bounds *app.Rectangle
	Depth  float32 //0 for lowest
}

func (t *Terminal) Init() {
	fmt.Printf("Terminal.Init()\n")

	t.TerminalId = msg.RandTerminalId()
	t.InChannel = make(chan []byte)
	t.GridSize = app.Vec2I{NumColumns, NumRows}
	t.Chars = [NumRows][NumColumns]uint32{}
	t.makeRandomChars(20)
	t.SetStringAt(37, 2, "this is text made by SetString()")
}

func (t *Terminal) Clear() {
	for y := 0; y < t.GridSize.Y; y++ {
		for x := 0; x < t.GridSize.X; x++ {
			t.Chars[y][x] = 0
		}
	}
}

func (t *Terminal) RelayToTask(message []byte) {
	println("(viewport/terminal/terminal.go).RelayToTask(message []byte)")
	println("FIX THIS: hypervisor.ProcessListGlobal.ProcessMap[t.AttachedProcess].GetIncomingChannel() <- message")

	//ch := hypervisor.ProcessListGlobal.ProcessMap[t.AttachedProcess].GetIncomingChannel()
	// ch <- message // <------- locks up

	//TODO: send this message to AttachedProcess,
	//and have AttachedProcess send SetChar*/SetCursor/etc. back here
}

func (t *Terminal) PutCharacter(m msg.MessageChar) {
	t.SetCharacter(m.Char)
	t.MoveRight()
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

func (t *Terminal) SetCursor(X uint32, Y uint32) {
	//fmt.Printf("Terminal.SetCursor()\n")
	t.Curs.X = X
	t.Curs.Y = Y
}

func (t *Terminal) SetCharacter(Char uint32) {
	//fmt.Printf("Terminal.SetCharacter()\n")
	t.SetCharacterAt(t.Curs.X, t.Curs.Y, Char)
}

func (t *Terminal) SetCharacterAt(X uint32, Y uint32, Char uint32) {
	//fmt.Printf("Terminal.SetCharacterAt()\n")

	numOOB = 0
	if t.validPos(X, Y) {
		t.Chars[Y][X] = Char
	}
}

func (t *Terminal) SetString(s string) {
	//fmt.Printf("Terminal.SetString()\n")
	t.SetStringAt(t.Curs.X, t.Curs.Y, s)
}

func (t *Terminal) SetStringAt(X uint32, Y uint32, S string) {
	//fmt.Printf("Terminal.SetStringAt()\n")

	numOOB = 0
	for x, c := range S {
		if t.validPos(X+uint32(x), Y) {
			t.Chars[Y][X+uint32(x)] = uint32(c)
		}
	}
}

func (t *Terminal) SetGridSize() {
	fmt.Printf("Terminal.SetGridSize()\n")

	// ERROR IN THIS CODE!
	// .... but I don't see any need for this right now anyways.
	// This USED to be "SetSize", which set the GL/float-based size of
	// the panel/window, rather than number of characters

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
func (t *Terminal) validPos(X, Y uint32) bool {
	if X < 0 || X >= uint32(t.GridSize.X) ||
		Y < 0 || Y >= uint32(t.GridSize.Y) {
		if numOOB < 1 {
			println("****** ATTEMPTED OUT OF BOUNDS CHARACTER PLACEMENT! ******")
		}

		numOOB++
		return false
	}

	return true
}
