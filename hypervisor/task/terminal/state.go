package task

import (
	"github.com/skycoin/viscript/app"
	"github.com/skycoin/viscript/msg"
)

var stPath = "hypervisor/task/terminal/state"

type State struct {
	DebugPrintInputEvents bool
	Cli                   *Cli
	VisualInfo            msg.MessageVisualInfo //dimensions, etc. (Terminal sends/updates)
	task                  *Task
	storedTerminalIds     []msg.TerminalId
}

func (st *State) Init(task *Task) {
	st.task = task
	st.DebugPrintInputEvents = true
	st.Cli = NewCli()
	println("st.VisualInfo.NumColumns", st.VisualInfo.NumColumns)
	st.Cli.AddEntriesForLogAndVisualRowsCache(app.HelpText, 80)
}

func (st *State) NumBackscrollRows() int {
	//the extra one is because of the "^^^backscroll^^^" indicator line
	return int(st.VisualInfo.NumRows) - int(st.VisualInfo.PromptRows) - 1
}

func (st *State) HandleMessages() {
	//called per Tick()
	c := st.task.InChannel

	for len(c) > 0 {
		m := <-c
		//TODO/FIXME:   cache channel id wherever it may be needed
		m = m[4:] //.....for now, DISCARD the chan id prefix
		msgType := msg.GetType(m)
		msgCategory := msgType & 0xff00 // get back masked category

		switch msgCategory {

		case msg.CATEGORY_Input:
			st.UnpackMessage(msgType, m)
		case msg.CATEGORY_Terminal:
			st.UnpackMessage(msgType, m)
		default:
			app.At(stPath, "**************** UNHANDLED MESSAGE TYPE CATEGORY! ****************")

		}
	}
}
