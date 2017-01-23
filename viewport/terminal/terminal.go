package terminal

import (
	"fmt"
	"github.com/corpusc/viscript/app"
	"github.com/corpusc/viscript/msg"
	//"math"
)

type Terminal struct {
	TerminalId      msg.TerminalId
	AttachedProcess msg.ProcessId

	//vars for character grid (of cells/units)
	Chars    [][]uint32
	Curs     app.Vec2I //current cursor/insert pos
	GridSize app.Vec2I //number of characters

	//vars for GL space / float
	//(mouse pos events are the only things that use pixels)
	Bounds *app.Rectangle // the whole panel, including chrome (title bar & scroll bars)
	Depth  int            //0 for lowest
}

func (t *Terminal) Init() {
	fmt.Printf("Terminal.Init()\n")

	t.TerminalId = msg.RandTerminalId()
	t.GridSize = app.Vec2I{64, 32}
	t.SetPanelSize()
}

func (t *Terminal) SetCursor(X uint32, Y uint32) {
	fmt.Printf("Terminal.SetCursor()\n")
	t.Curs.X = int(X)
	t.Curs.Y = int(Y)
}

func (t *Terminal) SetCharacter(X uint32, Y uint32, Char uint32) {
	fmt.Printf("Terminal.SetCharacter()\n")
	//do bounds check
	t.Chars[Y][X] = Char
}

func (t *Terminal) SetPanelSize() {
	fmt.Printf("Terminal.SetPanelSize()\n")
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
