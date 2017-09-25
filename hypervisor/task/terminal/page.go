package task

import (
	"github.com/skycoin/viscript/app"
	"github.com/skycoin/viscript/msg"
)

func (st *State) makePageOfLog(m msg.MessageVisualInfo) {
	//app.At("task/terminal/msg_action", "makePageOfLog")

	//called when:
	//		* receiving new/changed data via TypeVisualInfo msg/event
	//		* backscrolling (where only unchanged VisualInfo is passed)

	if /* VisualInfo changed */ m != st.VisualInfo {
		st.VisualInfo = m
		st.Cli.RebuildVisualRowsFromLogEntryFragments(m)
		println("makePageOfLog()   VisualInfo changed   -   .NumRows/Columns:", st.VisualInfo.NumRows, st.VisualInfo.NumColumns)
	} else {
		//println("VisualInfo UNchanged  -  st.VisualInfo.NumRows:", st.VisualInfo.NumRows)
	}

	st.printVisibleRows(m)

	//so you can see & interact with the command prompt even while backscrolled
	st.Cli.EchoWholeCommand(st.task.OutChannelId)
}

func (st *State) printVisibleRows(vi msg.MessageVisualInfo) {
	//println("printVisibleRows()") //...and indicator if backscrolled

	lvr := len(st.Cli.VisualRows)
	//println("len of st.Cli.VisualRows:", lvr)

	//(n)umber of (p)potential visible (r)ows   (that current Terminal height allows)
	npr := int(vi.NumRows - vi.PromptRows)

	if st.Cli.BackscrollAmount > 0 {
		npr--
		st.printRowsONLY(npr, lvr)
		ib /* indicator bar */ := app.GetLabeledBarOfChars(
			" BACKSCROLLED ", "^", st.VisualInfo.NumColumns)
		st.printLnAndMAYBELogIt(ib, false)
	} else {
		st.printRowsONLY(npr, lvr)
	}
}

var prevBackscrollAmount int

func (st *State) printRowsONLY(npr, lvr int) {
	if st.Cli.BackscrollAmount == 0 &&
		prevBackscrollAmount == 0 {
		return
	}

	//
	//
	prevBackscrollAmount = st.Cli.BackscrollAmount
	max := lvr - st.Cli.BackscrollAmount
	start := max - npr

	for i := start; i < max; i++ {
		if /* index is valid */ i >= 0 && i < lvr {
			//println("pVL i:", i)
			st.printLnAndMAYBELogIt(st.Cli.VisualRows[i], false)
		}
	}
}
