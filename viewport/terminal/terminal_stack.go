package terminal

import (
	"github.com/corpusc/viscript/app"
	"github.com/corpusc/viscript/gfx"
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

	// private
	nextRect *app.Rectangle // for next/new terminal spawn
	nextSpan float32        // how far from previous terminal
}

func (self *TerminalStack) Init() {
	println("TerminalStack.Init()")
	self.Terms = make(map[msg.TerminalId]*Terminal)
	self.nextSpan = 1.5
	self.nextRect = &app.Rectangle{
		3,
		gfx.DistanceFromOrigin * .7,
		-1,
		-4}
}

func (self *TerminalStack) AddTerminal() {
	println("TerminalStack.AddTerminal()")

	rand := msg.RandTerminalId()
	self.Terms[rand] = &Terminal{Bounds: self.nextRect}
	self.Terms[rand].Init()

	self.nextRect.Top -= self.nextSpan
	self.nextRect.Right += self.nextSpan
	self.nextRect.Bottom -= self.nextSpan
	self.nextRect.Left += self.nextSpan
}

func (self *TerminalStack) RemoveTerminal(id int) {
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
