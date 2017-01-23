/*

------- TODO: -------

* KEY-BASED NAVIGATION
	* CTRL-HOME/END - PGUP/DN
* BACKSPACE/DELETE at the ends of lines
	pulls us up to prev line, or pulls up next line
* when there is no scrollbar, should be able to see/interact with text in that area
* when auto appending to the end of a terminal, scroll all the way down
		(manual activity in the middle could increase size, so do this only when appending to body)


------- LOWER PRIORITY POLISH -------

* if typing goes past right of screen, auto-horizontal-scroll as you type
* same for when newlines/enters/returns push cursor past the bottom of visible space
* scrollbars should have a bottom right corner, and a thickness sized background
		for void space, reserved for only that, so the bar never covers up the rightmost
		character/cursor
* when pressing delete at/after the end of a line, should pull up the line below
* vertical scrollbars could have a smaller rendering of the first ~40 chars?
		however not if we map the whole vertical space (when scrollspace is taller than screen),
		because this requires scaling the text.  and keeping the aspect ratio means ~40 (max)
		would alter the width of the scrollbar

*/

package main

import (
	"fmt"
	"github.com/corpusc/viscript/viewport"
	igl "github.com/corpusc/viscript/viewport/gl" //eliminate
)

func main() {
	viewport.DebugPrintInputEvents = true //print input events

	fmt.Printf("Start\n")

	viewport.ViewportInit() //runtime.LockOSThread()
	viewport.ViewportScreenInit()
	viewport.ViewportInitInputEvents()

	//hypervisor.HypervisorInitProcessList()
	//hypervisor.AddTestProcess()

	fmt.Printf("Start Loop; \n")
	for viewport.CloseWindow == false {
		viewport.DispatchInputEvents(igl.InputEvents) //event channel
		//hypervisor.DispatchProcesEvents()               //viewport handles incoming process events
		//hypervisor.ProcessTick()                        //processes, handle incoming events
		viewport.PollUiInputEvents()
		viewport.Update() //in general
		viewport.UpdateDrawBuffer()
		viewport.SwapDrawBuffer() //with new frame
	}

	fmt.Printf("Closing down viewport \n")
	viewport.ViewportScreenTeardown()
	//viewport.HypervisorProcessListTeardown()

}
