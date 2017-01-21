package hypervisor

import (
	"fmt"
	"github.com/corpusc/viscript/msg"
	//"github.com/corpusc/viscript/app"
	//"github.com/corpusc/viscript/tree"
	//"math"
)

type Terminal struct {
	TerminalId      msg.TermininalId
	AttachedProcess msg.ProcessId

	CursX int // current cursor/insert position (in character grid cells/units)
	CursY int
	Xsize int //number of characters in X
	Ysize int //number of characters in Y

	CharacterArray []uint32 //array of characters

	//Whole *app.Rectangle // the whole panel, including chrome (title bar & scroll bars)

	DrawDepth int //0 for lowest
	DrawXsize int //in pixels
	DrawYsize int //in pixels
	DrawXOff  int //in pixels
	DrawYoff  int //in pixels

	//add draw x offset
	//add draw y offset
	//add z or draw
}

func (t *Terminal) Init() {
	fmt.Printf("Terminal.Init()\n")

	self.TerminalId = msg.RandTerminalId()

	Xsize = 32 //default
	Ysize = 64 //default

	t.SetSize()
}

func (t *Terminal) SetCursor(X uint32, Y uint32) {
	//implement
}

func (t *Terminal) SetCharacter(X uint32, Y uint32, Char uint32) {
	//implement
}

func (t *Terminal) SetSize() {
	fmt.Printf("Terminal.SetSize()\n")

}
