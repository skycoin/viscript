package terminal

import (
	"github.com/corpusc/viscript/app"
	"github.com/corpusc/viscript/msg"
	"github.com/corpusc/viscript/viewport/gl"
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
	nextRect app.Rectangle // for next/new terminal spawn
	nextSpan float32       // how far from previous terminal
}

func (self *TerminalStack) Init() {
	println("TerminalStack.Init()")
	self.Terms = make(map[msg.TerminalId]*Terminal)
	self.nextSpan = .8
	self.nextRect = app.Rectangle{
		gl.DistanceFromOrigin,
		gl.DistanceFromOrigin,
		-gl.DistanceFromOrigin,
		-gl.DistanceFromOrigin}
}

func (self *TerminalStack) AddTerminal() {
	println("TerminalStack.AddTerminal()")

	rand := msg.RandTerminalId()
	self.Terms[rand] = &Terminal{Bounds: &app.Rectangle{
		self.nextRect.Top,
		self.nextRect.Right,
		self.nextRect.Bottom,
		self.nextRect.Left}}
	self.Terms[rand].Init()

	self.nextRect.Top -= self.nextSpan
	self.nextRect.Right += self.nextSpan
	self.nextRect.Bottom -= self.nextSpan
	self.nextRect.Left += self.nextSpan
}

func (self *TerminalStack) RemoveTerminal(id int) {
	println("TerminalStack.RemoveTerminal()")
}

func (self *TerminalStack) ResizeTerminal(id msg.TerminalId, x int, y int) {
	println("TerminalStack.ResizeTerminal()")
}

func (self *TerminalStack) MoveTerminal(id msg.TerminalId, xoff int, yoff int) {
	println("TerminalStack.MoveTerminal()")
}
