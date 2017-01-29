package terminal

import (
	"fmt"

	"github.com/corpusc/viscript/app"
	"github.com/corpusc/viscript/msg"
	"github.com/corpusc/viscript/viewport/gl"

	"github.com/corpusc/viscript/hypervisor"
	"github.com/corpusc/viscript/hypervisor/dbus"
	example "github.com/corpusc/viscript/hypervisor/process/example"
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
	nextRect  app.Rectangle // for next/new terminal spawn
	nextDepth float32
	nextSpan  float32 // how far from previous terminal
}

func (self *TerminalStack) Init() {
	println("TerminalStack.Init()")
	self.Terms = make(map[msg.TerminalId]*Terminal)
	self.nextSpan = .3
	self.nextRect = app.Rectangle{
		gl.DistanceFromOrigin,
		gl.DistanceFromOrigin,
		-gl.DistanceFromOrigin,
		-gl.DistanceFromOrigin}
}

func (self *TerminalStack) AddTerminal() {
	println("TerminalStack.AddTerminal()")

	self.nextDepth += self.nextSpan / 10 // done first, cuz desktop is at 0

	tid := msg.RandTerminalId() //terminal id
	self.Terms[tid] = &Terminal{Depth: self.nextDepth, Bounds: &app.Rectangle{
		self.nextRect.Top,
		self.nextRect.Right,
		self.nextRect.Bottom,
		self.nextRect.Left}}
	self.Terms[tid].Init()

	self.nextRect.Top -= self.nextSpan
	self.nextRect.Right += self.nextSpan
	self.nextRect.Bottom -= self.nextSpan
	self.nextRect.Left += self.nextSpan

	//hook up proccess
	self.SetupTerminalDbus(tid)
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

func (self *TerminalStack) SetupTerminalDbus(TerminalId msg.TerminalId) {
	//create process

	//self.Terms[rand].AttachedProcess

	//create process
	var p *example.Process = example.NewProcess()
	var pi msg.ProcessInterface = msg.ProcessInterface(p)
	ProcessId := hypervisor.AddProcess(pi)

	self.Terms[TerminalId].AttachedProcess = ProcessId

	//terminalId and process Id
	//setup dbus
	//hypervisor.DbusGlobal.CreatePubsubChannel(Owner, OwnerType, ResourceIdentifier)

	//terminal dbus
	rid1 := fmt.Sprintf("dbus.pubsub.terminal-%d", int(TerminalId))
	tcid := hypervisor.DbusGlobal.CreatePubsubChannel(dbus.ResourceId(TerminalId), dbus.ResourceTypeTerminal, rid1)
	//func (self *DbusInstance) CreatePubsubChannel(Owner ResourceId, OwnerType ResourceType, ResourceIdentifier string) {

	//process dbus
	rid2 := fmt.Sprintf("dbus.pubsub.process-%d", int(ProcessId))
	pcid := hypervisor.DbusGlobal.CreatePubsubChannel(dbus.ResourceId(ProcessId), dbus.ResourceTypeProcess, rid2)

	//AddPubsubChannelSubscriber(ChannelId ChannelId, ResourceId ResourceId, ResourceType ResourceType) {}
	//AddPubsubChannelSubscriber(ChannelId ChannelId, ResourceId ResourceId, ResourceType ResourceType, channelIn chan []byte) {

	var c chan []byte //needs incoming channel
	c = make(chan []byte)

	//subscribe process to the terminal id
	hypervisor.DbusGlobal.AddPubsubChannelSubscriber(tcid, dbus.ResourceId(ProcessId), dbus.ResourceTypeProcess, c)

	//subscribe process to the process id
	hypervisor.DbusGlobal.AddPubsubChannelSubscriber(pcid, dbus.ResourceId(TerminalId), dbus.ResourceTypeTerminal, c)

}
