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
	- delete terminal
	- draw terminal state
	- resize terminal (in pixels or chars)
	- move terminal
*/

type TerminalStack struct {
	FocusedId msg.TerminalId
	Focused   *Terminal
	Terms     map[msg.TerminalId]*Terminal

	// private
	nextRect   app.Rectangle // for next/new terminal spawn
	nextDepth  float32
	nextOffset float32 // how far from previous terminal
}

func (ts *TerminalStack) Init() {
	println("TerminalStack.Init()")
	ts.Terms = make(map[msg.TerminalId]*Terminal)
	ts.nextOffset = gl.CanvasExtents.Y / 3
	ts.nextRect = app.Rectangle{
		gl.CanvasExtents.Y,
		gl.CanvasExtents.X / 2,
		-gl.CanvasExtents.Y / 2,
		-gl.CanvasExtents.X}
}

func (ts *TerminalStack) AddTerminal() msg.TerminalId {
	println("TerminalStack.AddTerminal()")

	ts.nextDepth += ts.nextOffset / 10 // done first, cuz desktop is at 0

	tid := msg.RandTerminalId() //terminal id
	ts.Terms[tid] = &Terminal{
		Depth: ts.nextDepth,
		Bounds: &app.Rectangle{
			ts.nextRect.Top,
			ts.nextRect.Right,
			ts.nextRect.Bottom,
			ts.nextRect.Left}}
	ts.Terms[tid].Init()
	ts.FocusedId = tid
	ts.Focused = ts.Terms[tid]

	ts.nextRect.Top -= ts.nextOffset
	ts.nextRect.Right += ts.nextOffset
	ts.nextRect.Bottom -= ts.nextOffset
	ts.nextRect.Left += ts.nextOffset

	//hook up proccess
	ts.SetupTerminalDbus(tid)

	return tid
}

func (ts *TerminalStack) RemoveTerminal(id msg.TerminalId) {
	println("TerminalStack.RemoveTerminal()")
	// delete(ts.Terms, id)
	// TODO: what should happen here after deleting terminal from the stack?
}

func (ts *TerminalStack) Tick() {
	//println("TerminalStack.Tick()")

	for _, term := range ts.Terms {
		term.Tick()
	}
}

func (ts *TerminalStack) ResizeFocusedTerminalRight(newRightOffset float32) {
	println("TerminalStack.ResizeFocusedTerminalRight()")
	ts.Focused.ResizingRight = true
	ts.Focused.Bounds.Right = newRightOffset
}

func (ts *TerminalStack) ResizeFocusedTerminalBottom(newBottomOffset float32) {
	println("TerminalStack.ResizeFocusedTerminalBottom()")
	ts.Focused.ResizingBottom = true
	ts.Focused.Bounds.Bottom = newBottomOffset
}

func (ts *TerminalStack) MoveFocusedTerminal(offset app.Vec2F) {
	println("TerminalStack.MoveTerminal()")
	bounds := ts.Focused.Bounds
	bounds.Top += offset.Y
	bounds.Bottom += offset.Y
	bounds.Left += offset.X
	bounds.Right += offset.X
}

func (ts *TerminalStack) SetupTerminalDbus(TerminalId msg.TerminalId) {
	println("TerminalStack.SetupTerminalDbus()")

	//create process
	var p *termTask.Process = termTask.NewProcess()
	var pi msg.ProcessInterface = msg.ProcessInterface(p)
	ProcessId := hypervisor.AddProcess(pi)

	//terminal dbus
	rid1 := fmt.Sprintf("dbus.pubsub.terminal-%d", int(TerminalId)) //ResourceIdentifier
	tcid := hypervisor.DbusGlobal.CreatePubsubChannel(              //terminal channel id
		dbus.ResourceId(TerminalId), //owner id
		dbus.ResourceTypeTerminal,   //owner type
		rid1)

	//process dbus
	rid2 := fmt.Sprintf("dbus.pubsub.process-%d", int(ProcessId)) //ResourceIdentifier
	pcid := hypervisor.DbusGlobal.CreatePubsubChannel(            //process channel id
		dbus.ResourceId(ProcessId), //owner id
		dbus.ResourceTypeProcess,   //owner type
		rid2)

	p.OutChannelId = uint32(tcid)
	ts.Terms[TerminalId].OutChannelId = uint32(pcid)
	ts.Terms[TerminalId].AttachedProcess = ProcessId

	//subscribe process to the terminal id
	hypervisor.DbusGlobal.AddPubsubChannelSubscriber(
		tcid,
		dbus.ResourceId(ProcessId),
		dbus.ResourceTypeProcess,
		ts.Terms[TerminalId].InChannel) // (a 2nd call had: p.GetIncomingChannel() as last parameter)

	// fmt.Printf("\nPubSub Channel After Adding Subscriber\n %+v\n",
	// 	hypervisor.DbusGlobal.PubsubChannels[tcid])

	//subscribe terminal to the process id
	hypervisor.DbusGlobal.AddPubsubChannelSubscriber(
		pcid,
		dbus.ResourceId(TerminalId),
		dbus.ResourceTypeTerminal,
		pi.GetIncomingChannel()) // (a 2nd call had: ts.Terms[TerminalId].InChannel) as last parameter)

	// fmt.Printf("\nPubSub Channel After Adding Subscriber\n %+v\n",
	// 	hypervisor.DbusGlobal.PubsubChannels[pcid])
}

func (ts *TerminalStack) SetFocused(topmostId msg.TerminalId) {
	//store which is focused and bring it to top
	newZ := float32(9.9) //FIXME (for all uses of this var, IF you ever want tons of terms)
	ts.FocusedId = topmostId
	ts.Focused = ts.Terms[topmostId]
	ts.Focused.Depth = newZ

	//store the REST of the terms
	theRest := []*Terminal{}

	for id, t := range ts.Terms {
		if id != ts.FocusedId {
			theRest = append(theRest, t)
		}
	}

	//sort them (top/closest at the start of list)
	fullySorted := false
	for !fullySorted {
		fullySorted = true

		for i := 0; i < len(theRest)-1; i++ {
			if theRest[i].Depth < theRest[i+1].Depth {
				theNext := theRest[i+1]
				theRest[i+1] = theRest[i]
				theRest[i] = theNext
				fullySorted = false
			}
		}
	}

	//assign receding z/depth values
	for _, t := range theRest {
		newZ -= 0.2
		t.Depth = newZ
	}
}

func (ts *TerminalStack) Defocus() {
	ts.FocusedId = 0
	ts.Focused = nil
}
