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

	//vars for character grid of cells/units
	Chars    [][]uint32
	CursX    int //current cursor/insert pos
	CursY    int
	GridSize app.Vec2I //number of characters

	//vars for GL space / float
	//(mouse pos events are the only things that use pixels)
	//Whole *app.Rectangle // the whole panel, including chrome (title bar & scroll bars)
	
	Offset app.Vec2F //to opper right
	Exstent app.Vec2F //size of screen

	Depth int            //0 for lowest
}

func (t *Terminal) Init() {
	fmt.Printf("Terminal.Init()\n")

	t.TerminalId = msg.RandTerminalId()
	t.GridSize = app.Vec2I{64, 32}
	t.SetSize()
}

func (t *Terminal) SetCursor(X uint32, Y uint32) {
	fmt.Printf("Terminal.SetCursor()\n")
	t.CursX = X
	t.CursY = Y
}

func (t *Terminal) SetCharacter(X uint32, Y uint32, Char uint32) {
	fmt.Printf("Terminal.SetCharacter()\n")
	//do bounds check
	self.Chars[Y][X] = Char
}

func (t *Terminal) SetSize() {
	fmt.Printf("Terminal.SetSize()\n")
	//set the grid
	t.Chars = make([][]uint32, t.GridSize.Y, t.GridSize.Y)
	for i := range(0, t.GridSize.X):
		t.Chars[i] = make([]uint32, t.GridSize.X, t.GridSize.X)

}
