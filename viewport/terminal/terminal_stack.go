package terminal

import (
	"fmt"

	"github.com/corpusc/viscript/app"
	"github.com/corpusc/viscript/hypervisor"
	"github.com/corpusc/viscript/hypervisor/dbus"
	termTask "github.com/corpusc/viscript/hypervisor/process/terminal"
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
	FocusedId msg.TerminalId
	Focused   *Terminal
	Terms     map[msg.TerminalId]*Terminal

	// private
	nextRect  app.Rectangle // for next/new terminal spawn
	nextDepth float32
	nextSpan  float32 // how far from previous terminal
}

func (self *TerminalStack) Init() {
	println("TerminalStack.Init()")
	self.Terms = make(map[msg.TerminalId]*Terminal)
	self.nextSpan = gl.CanvasExtents.Y / 3
	self.nextRect = app.Rectangle{
		gl.CanvasExtents.Y,
		gl.CanvasExtents.X / 2,
		-gl.CanvasExtents.Y / 2,
		-gl.CanvasExtents.X}
}

func (self *TerminalStack) AddTerminal() {
	println("TerminalStack.AddTerminal()")

	self.nextDepth += self.nextSpan / 10 // done first, cuz desktop is at 0

	tid := msg.RandTerminalId() //terminal id
	self.Terms[tid] = &Terminal{
		Depth: self.nextDepth,
		Bounds: &app.Rectangle{
			self.nextRect.Top,
			self.nextRect.Right,
			self.nextRect.Bottom,
			self.nextRect.Left}}
	self.Terms[tid].Init()
	self.FocusedId = tid
	self.Focused = self.Terms[tid]

	self.nextRect.Top -= self.nextSpan
	self.nextRect.Right += self.nextSpan
	self.nextRect.Bottom -= self.nextSpan
	self.nextRect.Left += self.nextSpan

	//hook up proccess
	self.SetupTerminalDbus(tid)
}

func (self *TerminalStack) RemoveTerminal(id msg.TerminalId) {
	println("TerminalStack.RemoveTerminal()")
	// delete(self.Terms, id)
	// TODO: what should happen here after deleting terminal from the stack?
}

func (self *TerminalStack) Tick() {
	//println("TerminalStack.Tick()")

	for _, term := range self.Terms {
		term.Tick()
	}
}

func (self *TerminalStack) ResizeTerminal(id msg.TerminalId, x int, y int) {
	println("TerminalStack.ResizeTerminal()")
}

func (self *TerminalStack) MoveTerminal(id msg.TerminalId, xoff int, yoff int) {
	println("TerminalStack.MoveTerminal()")
}

func (self *TerminalStack) SetupTerminalDbus(TerminalId msg.TerminalId) {
	println("TerminalStack.SetupTerminalDbus()")

	//create process
	var p *termTask.Process = termTask.NewProcess()
	var pi msg.ProcessInterface = msg.ProcessInterface(p)
	ProcessId := hypervisor.AddProcess(pi)

	self.Terms[TerminalId].AttachedProcess = ProcessId

	//terminal dbus, terminal channel id
	rid1 := fmt.Sprintf("dbus.pubsub.terminal-%d", int(TerminalId)) //ResourceIdentifier
	tcid := hypervisor.DbusGlobal.CreatePubsubChannel(
		dbus.ResourceId(TerminalId), //owner id
		dbus.ResourceTypeTerminal,   //owner type
		rid1)

	//process dbus, process channel id
	rid2 := fmt.Sprintf("dbus.pubsub.process-%d", int(ProcessId)) //ResourceIdentifier
	pcid := hypervisor.DbusGlobal.CreatePubsubChannel(
		dbus.ResourceId(ProcessId), //owner id
		dbus.ResourceTypeProcess,   //owner type
		rid2)
	self.Terms[TerminalId].OutChannelId = tcid
	p.PubSubChannelId = pcid

	//subscribe process to the terminal id
	hypervisor.DbusGlobal.AddPubsubChannelSubscriber(
		tcid,
		dbus.ResourceId(ProcessId),
		dbus.ResourceTypeProcess,
		self.Terms[TerminalId].InChannel) // (a 2nd call had: p.GetIncomingChannel() as last parameter)

	fmt.Printf("\nPubSub Channel After Adding Terminal\n %+v\n",
		hypervisor.DbusGlobal.PubsubChannels[tcid])

	//subscribe terminal to the process id
	hypervisor.DbusGlobal.AddPubsubChannelSubscriber(
		pcid,
		dbus.ResourceId(TerminalId),
		dbus.ResourceTypeTerminal,
		pi.GetIncomingChannel()) // (a 2nd call had: self.Terms[TerminalId].InChannel) as last parameter)

	fmt.Printf("\nPubSub Channel After Adding Terminal\n %+v\n",
		hypervisor.DbusGlobal.PubsubChannels[pcid])

	/* we don't need this 2nd pair of add sub calls?
	//subscribe process to the terminal id
	hypervisor.DbusGlobal.AddPubsubChannelSubscriber(
		tcid,
		dbus.ResourceId(ProcessId),
		dbus.ResourceTypeProcess,
		self.Terms[TerminalId].InChannel)

	//subscribe terminal to the process id
	hypervisor.DbusGlobal.AddPubsubChannelSubscriber(
		pcid,
		dbus.ResourceId(TerminalId),
		dbus.ResourceTypeTerminal,
		pi.GetIncomingChannel())
	*/
}
