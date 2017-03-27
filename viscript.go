/*

------- NEXT THINGS TODO: -------

* ExternalProcess:
	attach,
	send to background (https://repl.it/GeGn/1),
	send to foreground,
	list jobs (external processes)

* limit resizing to require at least 16 char columns

* make current command line autoscroll horizontally
	* make it optional (if turned off, always truncate the left)

* Resize terminals
	* i believe we'll change the actual grid size, then get enough data
			from the terminal task to fill the backscroll
	* change flow of text with wrapping, so for example,
			squeezing horizontally would cause more lines

* back buffer scrolling
	* pgup/pgdn hotkeys
	* 1-3 lines with scrollwheel

* Fix getting a resizing pointer outside of focused terminal.
		When you click outside terminal it can land on a background
		terminal which then pops in front.  Blocking the resize




------- OLDER TODO: ------- (everything below was for the text editor)

* KEY-BASED NAVIGATION
	* CTRL-HOME/END - PGUP/DN
* BACKSPACE/DELETE at the ends of lines
	pulls us up to prev line, or pulls up next line
* when auto appending to the end of a terminal, scroll all the way down
		(manual activity in the middle could increase size, so do this only when appending to body)


------- LOWER PRIORITY POLISH: -------

* if cursor movement goes past left/right of screen, auto-horizontal-scroll as you type
* same for when newlines/enters/returns push cursor past the bottom of visible space
* vertical scrollbars could have a smaller rendering of the first ~40 chars?
		however not if we map the whole vertical space (when scrollspace is taller than screen),
		because this requires scaling the text.  and keeping the aspect ratio means ~40 (max)
		would alter the width of the scrollbar
* when there is no scrollbar, should be able to see/interact with text in that area

*/

package main

import (
	"github.com/corpusc/viscript/app"
	"github.com/corpusc/viscript/god"
	"github.com/corpusc/viscript/hypervisor"
	"github.com/corpusc/viscript/rpc/terminalmanager"
)

func main() {
	app.MakeHighlyVisibleLogHeader(app.Name, 15)

	hypervisor.Init()

	god.DebugPrintInputEvents = true
	god.Init() //runtime.LockOSThread()

	// rpc
	go func() {
		rpcInstance := terminalmanager.NewRPC()
		rpcInstance.Serve()
	}()

	println("Start Loop;")
	for god.CloseWindow == false {
		god.DispatchEvents() //event channel

		hypervisor.TickTasks()

		god.PollUiInputEvents()
		god.Tick()
		god.UpdateDrawBuffer()
		god.SwapDrawBuffer() //with new frame
	}

	god.TeardownScreen()
	hypervisor.HypervisorTeardown()
}
