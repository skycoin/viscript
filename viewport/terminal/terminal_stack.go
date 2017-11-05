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
	TermMap   map[msg.TerminalId]*Terminal

	//private
	//next/new terminal spawn vars
	nextRect   app.Rectangle
	nextDepth  float32
	nextOffset app.Vec2F //how far from previous terminal
}

func (ts *TerminalStack) Init() {
	top := gl.CanvasExtents.Y
	left := -gl.CanvasExtents.X
	w := gl.CanvasExtents.X * 1.5 //width of terminal window
	h := gl.CanvasExtents.Y * 1.5 //height " " "

	ts.TermMap = make(map[msg.TerminalId]*Terminal)
	ts.nextOffset.X = (gl.CanvasExtents.X*2 - w) / 2
	ts.nextOffset.Y = (gl.CanvasExtents.Y*2 - h) / 2

	ts.nextRect = app.Rectangle{
		top,
		left + w,
		top - h,
		left}

	//initial terminal window
	Terms.Add()
}

//these are visually grouped, because they're used as 1 unit (to allow for a sort of default parameter)
func (ts *TerminalStack) Add() msg.TerminalId {
	println("<TerminalStack>.Add()")
	return ts.AddWithFixedSizeState(false)
}
func (ts *TerminalStack) AddWithFixedSizeState(fixedSize bool) msg.TerminalId { //^^^
	println("<TerminalStack>.AddWithFixedSizeState()")

	ts.nextDepth += ts.nextOffset.X / 10 // done first, cuz desktop is at 0

	tid := msg.RandTerminalId() //terminal id
	ts.TermMap[tid] = &Terminal{
		Depth: ts.nextDepth,
		Bounds: &app.Rectangle{
			ts.nextRect.Top,
			ts.nextRect.Right,
			ts.nextRect.Bottom,
			ts.nextRect.Left},
		TaskBarButton: &app.Rectangle{}}
	ts.TermMap[tid].FixedSize = fixedSize
	ts.TermMap[tid].Init()
	println("AddWithFixed...... - ts.TermMap[tid].TerminalId:", ts.TermMap[tid].TerminalId)
	ts.SetFocused(ts.TermMap[tid].TerminalId)

	ts.nextRect.Top -= ts.nextOffset.Y
	ts.nextRect.Right += ts.nextOffset.X
	ts.nextRect.Bottom -= ts.nextOffset.Y
	ts.nextRect.Left += ts.nextOffset.X

	ts.SetupTerminal(tid)
	ts.setTaskbarButtonRectangles()
	return tid
}

func (ts *TerminalStack) Defocus() {
	ts.FocusedId = 0
}

func (ts *TerminalStack) GetFocusedTerminal() *Terminal {
	for key, t := range ts.TermMap {
		if t.TerminalId == ts.FocusedId {
			return ts.TermMap[key]
		}
	}

	return nil
}

func (ts *TerminalStack) MoveFocusedTerminal(hiResDelta app.Vec2F, mouseDeltaSinceClick *app.Vec2F) {
	d := mouseDeltaSinceClick
	println("MoveFocusedTerminal()   -   ts.FocusedId:", ts.FocusedId)
	foc := ts.GetFocusedTerminal()

	if foc == nil {
		return
	}

	cs := foc.CharSize
	fb := foc.Bounds

	if keyboard.AltKeyIsDown { //smooth, high resolution
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

func (ts *TerminalStack) Remove(id msg.TerminalId) {
	//this is only ever called after finding a valid id to remove
	println("<TerminalStack>.Remove():", id)

	var trash msg.TerminalId
	for key, term := range ts.TermMap {
		println("key:", key)

		if id == term.TerminalId {
			trash = key

			if term == ts.TermMap[ts.FocusedId] {
				ts.FocusedId = 0
			}
		}
	}

	attTask := ts.TermMap[trash].AttachedTask
	outId := hypervisor.GlobalTasks.TaskMap[attTask].GetOutChannelId()

	//remove both subscriptions from dbus
	hypervisor.DbusGlobal.RemoveChannel(dbus.ChannelId(outId))                          //task
	hypervisor.DbusGlobal.RemoveChannel(dbus.ChannelId(ts.TermMap[trash].OutChannelId)) //terminal

	//remove from GlobalTasks list
	//println("len of TaskMap:", len(hypervisor.GlobalTasks.TaskMap))
	delete(hypervisor.GlobalTasks.TaskMap, attTask)
	//println("len of TaskMap:", len(hypervisor.GlobalTasks.TaskMap))

	//remove from terminal stack
	//println("len of TermMap:", len(ts.TermMap))
	delete(ts.TermMap, trash)
	//println("len of TermMap:", len(ts.TermMap))

	ts.setTaskbarButtonRectangles()
}

func (ts *TerminalStack) SetupTerminal(termId msg.TerminalId) {
	//make it's task
	task := termTask.MakeNewTask()
	tskIF := msg.TaskInterface(task)
	tskId := hypervisor.AddTask(tskIF)

	task.State.VisualInfo = ts.TermMap[termId].GetVisualInfo()

	/* the rest is all DBUS related */

	//terminal
	tcid := hypervisor.DbusGlobal.CreatePubsubChannel( //terminal channel id
		dbus.ResourceId(termId),                             //owner id
		dbus.ResourceTypeTerminal,                           //owner type
		fmt.Sprintf("dbus.pubsub.terminal-%d", int(termId))) //ResourceIdentifier

	//task
	pcid := hypervisor.DbusGlobal.CreatePubsubChannel( //(task/)process channel id
		dbus.ResourceId(tskId),                         //owner id
		dbus.ResourceTypeTask,                          //owner type
		fmt.Sprintf("dbus.pubsub.task-%d", int(tskId))) //ResourceIdentifier

	task.OutChannelId = uint32(tcid)
	ts.TermMap[termId].OutChannelId = uint32(pcid)
	ts.TermMap[termId].AttachedTask = tskId

	//subscribe task to the terminal id
	hypervisor.DbusGlobal.AddPubsubChannelSubscriber(
		tcid,
		dbus.ResourceId(tskId),
		dbus.ResourceTypeTask,
		ts.TermMap[termId].InChannel)

	//subscribe terminal to the task id
	hypervisor.DbusGlobal.AddPubsubChannelSubscriber(
		pcid,
		dbus.ResourceId(termId),
		dbus.ResourceTypeTerminal,
		tskIF.GetIncomingChannel())
}

func (ts *TerminalStack) SetFocused(topmostId msg.TerminalId) {
	newZ := float32(9.9) //FIXME (@ all places of this var) IF you ever want more than (about) 50 terms
	ts.FocusedId = topmostId

	//cycle through terminals
	noneFocused := true
	for key, term := range ts.TermMap {
		if term.TerminalId == ts.FocusedId {
			noneFocused = false
			ts.TermMap[key].Depth = newZ
		}
	}

	if noneFocused {
		ts.TermMap[topmostId].Depth = newZ
	}

	//make list of REST of the terms (excluding currently focused term)
	theRest := []*Terminal{}

	for key, t := range ts.TermMap {
		if t.TerminalId != ts.FocusedId {
			theRest = append(theRest, ts.TermMap[key])
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

func (ts *TerminalStack) Tick() {
	for _, term := range ts.TermMap {
		term.Tick()
	}
}

//
//
//private
//
//

func (ts *TerminalStack) setTaskbarButtonRectangles() {
	x := -gl.CanvasExtents.X + app.TaskBarHeight - app.TaskBarBorderSpan
	//TODO!!!
	//SHRINK BUTTONS TO FIT canvas width

	for _, term := range ts.TermMap {
		currButtonWid := app.TaskBarCharWid*float32(len(term.TaskBarButtonText)) + app.TaskBarBorderSpan*2
		term.TaskBarButton.Left = x
		//term.TaskBarButton.Right = x + currButtonWid
		//(actually, we can't know the .TabText width until drabTabId() sets it,
		//...so we set it per frame in drawTerminalButtons() )
		term.TaskBarButton.Bottom = -gl.CanvasExtents.Y + app.TaskBarBorderSpan
		term.TaskBarButton.Top = -gl.CanvasExtents.Y - app.TaskBarBorderSpan + app.TaskBarHeight
		x += currButtonWid
	}
}
