package task

import (
	"github.com/skycoin/viscript/app"
	"github.com/skycoin/viscript/msg"
)

var prevBackscrollAmount int

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

	nvr := len(st.Cli.VisualRows) //number of visual rows

	//(n)umber of (l)eftover (r)ows
	//(...after dedicating row/s to the prompt
	//		& possibly the backscroll indicator)
	nlr := int(vi.NumRows - vi.PromptRows)

	if st.Cli.BackscrollAmount <= 0 {
		st.printLeftoverRows(nlr, nvr)
	} else {
		nlr--
		st.printLeftoverRows(nlr, nvr)

		//print indicator bar
		ib := app.GetLabeledBarOfChars(" BACKSCROLLED ", "^", st.VisualInfo.NumColumns)
		st.printLnAndMAYBELogIt(ib, false)
	}
}

func (st *State) printLeftoverRows(nlr, nvr int) {
	if st.Cli.BackscrollAmount == 0 &&
		prevBackscrollAmount == 0 {
		return
	}

	//
	//
	prevBackscrollAmount = st.Cli.BackscrollAmount
	max := nvr - st.Cli.BackscrollAmount
	start := max - nlr

	for i := start; i < max; i++ {
		if /* index is valid */ i >= 0 && i < nvr {
			//println("pVL i:", i)
			st.printLnAndMAYBELogIt(st.Cli.VisualRows[i], false)
		}
	}
}
