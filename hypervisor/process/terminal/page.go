package process

import (
	"github.com/skycoin/viscript/msg"
)

func (st *State) makePageOfLog(m msg.MessageVisualInfo) {
	//app.At("process/terminal/msg_action", "makePageOfLog")

	//called when:
	//		* receiving new/changed data via TypeVisualInfo msg/event
	//		* backscrolling (where only unchanged VisualInfo is passed)

	if /* VisualInfo changed */ m != st.VisualInfo {
		st.VisualInfo = m
		st.Cli.BuildRowsFromLogEntryFragments(m)
		println("VisualInfo changed   -   st.VisualInfo.NumRows:", st.VisualInfo.NumRows)
	} else {
		println("VisualInfo UNchanged  -  st.VisualInfo.NumRows:", st.VisualInfo.NumRows)
	}

	st.printVisibleRows(m)

	//so you can see & interact with the command prompt even while backscrolled
	st.Cli.EchoWholeCommand(st.proc.OutChannelId)
}

func (st *State) printVisibleRows(vi msg.MessageVisualInfo) {
	println("printVisibleRows()") //...and indicator if backscrolled

	num := len(st.Cli.VisualRows)
	println("len of st.Cli.VisualRows:", num)

	//(n)umber of (v)isible (r)ows   (that current Terminal height allows)
	nvr := int(vi.NumRows - vi.PromptRows)

	if st.Cli.BackscrollAmount > 0 {
		nvr--
		st.printRowsONLY(nvr, num)
		st.printLnAndMAYBELogIt("^^^^^^^^^^^^^^^^ BACKSCROLLED ^^^^^^^^^^^^^^^^", false)
	} else {
		st.printRowsONLY(nvr, num)
	}
}

func (st *State) printRowsONLY(nvr, num int) {
	end := num - st.Cli.BackscrollAmount
	start := end - nvr

	for i := start; i < end && i < num; i++ {
		if /* index not negative */ i > -1 {
			println("pVL i:", i)
			st.printLnAndMAYBELogIt(st.Cli.VisualRows[i], false)
		}
	}
}
