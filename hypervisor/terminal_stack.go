package hypervisor

import (
	"fmt"
	"github.com/corpusc/viscript/msg"
)

/*
	What operations?
	- create terminal
	- delete terminal
	- draw terminal state
	- change terminal in focus
	- resize terminal (in pixels or chars)
	- move terminal

*/
type TerminalStack struct {
	Focused msg.TerminalId
	Terms   map[msg.TerminalId]*Terminal
}

//var Focused *Terminal //use terminal ID
//var Terms []*Terminal

//var runOutputTerminalFrac = float32(0.4) // TEMPORARY fraction of vertical strip height which is dedicated to running code

func (self *TerminalStack) Init() {

}

func (self *TerminalStack) AddTerminal() {

}

func (self *TerminalStack) RemoveTerminal() {

}

//MOST IMPORTANT
func (self *TerminalStack) Draw() {}

//MOST IMPORTANT
func (self *TerminalStack) ResizeTerminal(id msg.TerminalId, x int, y int) {}

func (self *TerminalStack) MoveTerminal(id msg.TerminalId, xoff int, yoff int) {}

/*
func initTerminals() {

	Terms = append(Terms, &Terminal{})
	Terms = append(Terms, &Terminal{})

	Focused = Terms[0]

	Terms[0].Init()
	Terms[1].Init()
}
*/

// (maybe temporary) refactoring additions
//func SetSize() {
//	for _, t := range Terms {
//		t.SetSize()
//	}
//}
