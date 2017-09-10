/*

------- NEXT THINGS TODO: -------

* (Red typed this) RPC cli:
	add functionality to print running jobs for a given task id
	that can be retrieved by lp or setting the task id as default
	because that already exists

* (Red typed this) ExternalApp:
	Ctrl + c - detach, delete, kill probably
	Ctrl + z - detach and let it be running or pause it (https://repl.it/GeGn/1)?,
	jobs - list all jobs of current terminal
	fg <id> - send to foreground

* Sideways auto-scroll command line when it doesn't fit the dedicated space for it
		(atm, 2 lines are reserved along the bottom of a full screen)
		* block character at end to indicate continuing on next line
	* make it optional (if turned off, always truncate the left)

* scan and do/fix the most important FIXME/TODO places in the code



------- OLDER TODO: ------- (everything below was for the text editor)

* KEY-BASED NAVIGATION
	* CTRL-HOME/END - PGUP/DN
* BACKSPACE/DELETE at the ends of lines
	pulls us up to prev line, or pulls up next line
* when auto appending to the end of a terminal, scroll all the way down
		(manual activity in the middle could increase size, so do this only when appending to body)


------- LOWER PRIORITY POLISH: -------

* if cursor movement goes past left/right of screen, auto-horizontal-scroll the page as you type
* same for when newlines/enters/returns push cursor past the bottom of visible space
* vertical scrollbars could have a smaller rendering of the first ~40 chars?
		however not if we map the whole vertical space (when scrollspace is taller than screen),
		because this requires scaling the text.  and keeping the aspect ratio means ~40 (max)
		would alter the width of the scrollbar
* when there is no scrollbar needed, should be able to see/interact with text in that area?

*/

package main

import (
	"os"

	"github.com/skycoin/viscript/app"
	"github.com/skycoin/viscript/config"
	"github.com/skycoin/viscript/headless"
	"github.com/skycoin/viscript/hypervisor"
	"github.com/skycoin/viscript/reds_rpc"
	"github.com/skycoin/viscript/signal"
	"github.com/skycoin/viscript/viewport"
)

func main() {
	app.MakeHighlyVisibleLogEntry(app.Name, 13)
	loadConfig()
	handleAnyArguments()
	inits()

	go func() {
		rpcInstance := reds_rpc.NewRPC()
		rpcInstance.Serve()
	}()

	err := signal.Listen("0.0.0.0:7999")
	if err != nil {
		panic(err)
	}

	//start looping
	for viewport.CloseWindow == false {
		viewport.DispatchEvents() //event channel
		hypervisor.TickTasks()
		hypervisor.TickExternalApps()

		if config.Global.Settings.RunHeadless {
			headless.Tick()
		} else {
			viewport.PollUiInputEvents()
			viewport.Tick()
			viewport.UpdateDrawBuffer()
			viewport.SwapDrawBuffer() //with new frame
		}
	}

	viewport.TeardownScreen()
	hypervisor.Teardown()
}

func loadConfig() {
	err := config.Load("config.yaml")
	if err != nil {
		println(err.Error())
		return
	}
}

func handleAnyArguments() {
	args := os.Args[1:]
	if len(args) == 1 {
		if args[0] == "-h" || args[0] == "-run_headless" {
			config.Global.Settings.RunHeadless = true //override defalt
			app.MakeHighlyVisibleLogEntry("Running in HEADLESS MODE", 9)
		}
	}
}

func inits() {
	hypervisor.Init()

	if config.Global.Settings.RunHeadless {
		headless.Init()
	} else {
		viewport.Init() //runtime.LockOSThread()
	}
}
