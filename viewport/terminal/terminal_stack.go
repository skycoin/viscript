package terminal

import (
	"fmt"

	"github.com/skycoin/viscript/app"
	"github.com/skycoin/viscript/hypervisor"
	"github.com/skycoin/viscript/hypervisor/dbus"
	"github.com/skycoin/viscript/hypervisor/input/keyboard"
	termTask "github.com/skycoin/viscript/hypervisor/task/terminal"
	"github.com/skycoin/viscript/msg"
	"github.com/skycoin/viscript/viewport/gl"
)

var Terms = TerminalStack{}

type TerminalStack struct {
	FocusedId msg.TerminalId
	Focused   *Terminal
	Terms     map[msg.TerminalId]*Terminal

	//private
	//next/new terminal spawn vars
	nextRect   app.Rectangle
	nextDepth  float32
	nextOffset app.Vec2F // how far from previous terminal
}

func (ts *TerminalStack) Init() {
	w := gl.CanvasExtents.X * 1.5 //width of terminal window
	h := gl.CanvasExtents.Y * 1.5 //height

	ts.Terms = make(map[msg.TerminalId]*Terminal)
	ts.nextOffset.X = (gl.CanvasExtents.X*2 - w) / 2
	ts.nextOffset.Y = (gl.CanvasExtents.Y*2 - h) / 2

	ts.nextRect = app.Rectangle{
		gl.CanvasExtents.Y,
		-gl.CanvasExtents.X + w,
		gl.CanvasExtents.Y - h,
		-gl.CanvasExtents.X}

	//setup a starter terminal window
	Terms.Add()
}

//these are visually grouped, because they're used as 1 unit (to allow for a sort of default argument)
func (ts *TerminalStack) Add() msg.TerminalId {
	println("<TerminalStack>.Add()")
	return ts.AddWithFixedSizeState(false)
}
func (ts *TerminalStack) AddWithFixedSizeState(fixedSize bool) msg.TerminalId { //^^^
	println("<TerminalStack>.AddWithFixedSizeState()")

	ts.nextDepth += ts.nextOffset.X / 10 // done first, cuz desktop is at 0

	tid := msg.RandTerminalId() //terminal id
	ts.Terms[tid] = &Terminal{
		Depth: ts.nextDepth,
		Bounds: &app.Rectangle{
			ts.nextRect.Top,
			ts.nextRect.Right,
			ts.nextRect.Bottom,
			ts.nextRect.Left}}
	ts.Terms[tid].Init()
	ts.Terms[tid].FixedSize = fixedSize
	ts.SetFocused(tid)

	ts.nextRect.Top -= ts.nextOffset.Y
	ts.nextRect.Right += ts.nextOffset.X
	ts.nextRect.Bottom -= ts.nextOffset.Y
	ts.nextRect.Left += ts.nextOffset.X

	ts.SetupTerminal(tid)
	return tid
}

func (ts *TerminalStack) Remove(id msg.TerminalId) {
	println("<TerminalStack>.Remove():", id)
	//TODO: FIXME:
	//what should happen here after deleting terminal from the stack?                       -Red
	//well BEFORE removal, we'll want to unsub/update/cleanup dsub (not sure what else atm) -CC

	for _, term := range ts.Terms {
		if id == term.TerminalId {
			if term == ts.Focused {
				ts.Focused = nil
				ts.FocusedId = 0
			}

			delete(ts.Terms, id)
			//TODO?			//st.SendCommand("list_terms", []string{}) //updates tasks' stored list, on any remaining terminals
		}
	}
}

func (ts *TerminalStack) Tick() {
	for _, term := range ts.Terms {
		term.Tick()
	}
}

func (ts *TerminalStack) MoveFocusedTerminal(hiResDelta app.Vec2F, mouseDeltaSinceClick *app.Vec2F) {
	d := mouseDeltaSinceClick
	cs := ts.Focused.CharSize
	fb := ts.Focused.Bounds

	if keyboard.ControlKeyIsDown { //smooth, high resolution
		fb.MoveBy(hiResDelta)
	} else { //snap movement to char size
		if d.X > cs.X {
			d.X -= cs.X
			fb.MoveBy(app.Vec2F{cs.X, 0})
		} else if d.X < -cs.X {
			d.X += cs.X
			fb.MoveBy(app.Vec2F{-cs.X, 0})
		}

		if d.Y > cs.Y {
			d.Y -= cs.Y
			fb.MoveBy(app.Vec2F{0, cs.Y})
		} else if d.Y < -cs.Y {
			d.Y += cs.Y
			fb.MoveBy(app.Vec2F{0, -cs.Y})
		}
	}
}

func (ts *TerminalStack) SetupTerminal(termId msg.TerminalId) {
	//make it's task
	task := termTask.MakeNewTask()
	tskIF := msg.TaskInterface(task)
	tskId := hypervisor.AddTask(tskIF)

	task.State.VisualInfo = ts.Terms[termId].GetVisualInfo()

	/* the rest is all DBUS related */

	//terminal
	rid1 := fmt.Sprintf("dbus.pubsub.terminal-%d", int(termId)) //ResourceIdentifier
	tcid := hypervisor.DbusGlobal.CreatePubsubChannel(          //terminal channel id
		dbus.ResourceId(termId),   //owner id
		dbus.ResourceTypeTerminal, //owner type
		rid1)

	//task
	rid2 := fmt.Sprintf("dbus.pubsub.task-%d", int(tskId)) //ResourceIdentifier
	pcid := hypervisor.DbusGlobal.CreatePubsubChannel(     //task channel id
		dbus.ResourceId(tskId), //owner id
		dbus.ResourceTypeTask,  //owner type
		rid2)

	task.OutChannelId = uint32(tcid)
	ts.Terms[termId].OutChannelId = uint32(pcid)
	ts.Terms[termId].AttachedTask = tskId

	//subscribe task to the terminal id
	hypervisor.DbusGlobal.AddPubsubChannelSubscriber(
		tcid,
		dbus.ResourceId(tskId),
		dbus.ResourceTypeTask,
		ts.Terms[termId].InChannel)

	//subscribe terminal to the task id
	hypervisor.DbusGlobal.AddPubsubChannelSubscriber(
		pcid,
		dbus.ResourceId(termId),
		dbus.ResourceTypeTerminal,
		tskIF.GetIncomingChannel())
}

func (ts *TerminalStack) SetFocused(topmostId msg.TerminalId) {
	//store which is focused and bring it to top
	newZ := float32(9.9) //FIXME (@ all places of this var) IF you ever want more than (about) 50 terms
	ts.FocusedId = topmostId
	ts.Focused = ts.Terms[topmostId]
	ts.Focused.Depth = newZ

	//REST of the terms
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
