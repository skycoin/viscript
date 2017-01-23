package hypervisor

import (
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

func (self *TerminalStack) Init() {
	println("TerminalStack.Init()")
}

func (self *TerminalStack) AddTerminal() {
	println("TerminalStack.AddTerminal()")
}

func (self *TerminalStack) RemoveTerminal() {
	println("TerminalStack.RemoveTerminal()")
}

//MOST IMPORTANT
func (self *TerminalStack) Draw() {
	println("TerminalStack.Draw()")
}

func (self *TerminalStack) ResizeTerminal(id msg.TerminalId, x int, y int) {
	println("TerminalStack.ResizeTerminal()")
}

func (self *TerminalStack) MoveTerminal(id msg.TerminalId, xoff int, yoff int) {
	println("TerminalStack.MoveTerminal()")
}
