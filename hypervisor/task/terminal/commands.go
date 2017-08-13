package task

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/skycoin/viscript/app"
	"github.com/skycoin/viscript/config"
	"github.com/skycoin/viscript/hypervisor"
	extAppImport "github.com/skycoin/viscript/hypervisor/ext_app"
	"github.com/skycoin/viscript/msg"
	"github.com/skycoin/viscript/signal"
	"time"
)

const cp = "hypervisor/task/terminal/commands"

func (st *State) commandHelp() {
	st.PrintLn(app.GetBarOfChars("-", int(st.VisualInfo.NumColumns)))
	//st.PrintLn("Current commands:")
	st.PrintLn("------ Terminals ------")
	st.PrintLn("clear:                 Clears currently focused terminal.")
	st.PrintLn("close_term <id>:       Close terminal by id.")
	st.PrintLn("list_terms:            List all terminal ids.")
	st.PrintLn("new_term:              Add new terminal (n for short).")
	st.PrintLn("------ Apps -----------")
	st.PrintLn("apps:                  Display all available apps with descriptions.")
	st.PrintLn("attach    <id>:        Attach external task with given terminal id.")
	st.PrintLn("list_tasks (-f):       List running tasks (-f for full commands).")
	st.PrintLn("ping      <id>:        Ping app with given id.")
	st.PrintLn("res_usage <id>:        See resource usage for app with given id.")
	st.PrintLn("shutdown  <id>:        [TODO] Shutdown external task with given id.")
	st.PrintLn("start (-a) <command>:  Start external task. (-a to also attach).")
	// st.PrintLn("rpc:                   Issues command: \"go run rpc/cli/cli.go\"")
	// st.PrintLn("Current hotkeys:")
	st.PrintLn("CTRL+Z:                Detach currently attached task.")
	// st.PrintLn("    CTRL+C:           ___description goes here___")
	st.PrintLn(app.GetBarOfChars("-", int(st.VisualInfo.NumColumns)))
}

func (st *State) commandApps() {
	apps := config.Global.Apps

	if len(apps) == 0 {
		st.PrintLn("No available apps found.")
		return
	}

	maxAppKeyLength := 0

	for appKey, _ := range apps {
		if len(appKey) > maxAppKeyLength {
			maxAppKeyLength = len(appKey)
		}
	}

	maxAppKeyLength += 4 // Space after max string length app hash
	s := ""

	for appKey, app := range apps {
		s += appKey

		for i := 0; i < maxAppKeyLength-len(appKey); i++ {
			s += " "
		}

		s += fmt.Sprintf("-%s\n", app.Desc)
	}

	st.PrintLn(s)
	st.PrintLn("Type \"help name\" for expected parameters.")
}

func (st *State) commandAppHelp(args []string) {
	appName := args[0]

	if !config.AppExistsWithName(appName) {
		st.PrintError("App with name: " + appName + " doesn't exist. " +
			"Try running 'apps'.")
		return
	}

	st.PrintLn(config.Global.Apps[appName].Help)
}

func (st *State) commandClearTerminal() {
	st.VisualInfo.CurrRow = 0
	st.publishToOut(msg.Serialize(msg.TypeClear, msg.MessageClear{}))
	st.Cli.EchoWholeCommand(st.task.OutChannelId)
}

func (st *State) commandStart(args []string) {
	app.At(cp, "commandStart")

	if len(args) < 1 {
		st.PrintError("Must pass a command into Start!")
		return
	}

	detached := args[0] != "-a"

	if !detached {
		args = args[1:]
	}

	appName := args[0]

	if !config.AppExistsWithName(appName) {
		st.PrintError("App with name: " + appName + " doesn't exist. " +
			"Try running 'apps'.")
		return
	}

	var tokens []string

	//if there are user passed args for the app override defaults set in config
	if len(args) > 1 {
		pathToApp := config.GetPathForApp(appName)
		tokens = append(tokens, pathToApp)

		for _, arg := range args[1:] {
			tokens = append(tokens, strings.ToLower(arg))
		}
	} else {
		tokens = config.GetPathWithDefaultArgsForApp(appName)
	}

	//if the app is daemon not allow to attach to it
	if config.Global.Apps[appName].Daemon {
		detached = true
	}

	newExternalApp, err := extAppImport.MakeNewExternalApp(tokens, detached)
	if err != nil {
		st.PrintError(err.Error())
		return
	}

	err = newExternalApp.Start()
	if err != nil {
		st.PrintError(err.Error())
		return
	}

	eai := newExternalApp.GetExternalAppInterface()
	appId := hypervisor.AddExternalApp(eai)

	if !detached {
		err = st.task.AttachExternalApp(eai)
		if err != nil {
			st.PrintError(err.Error())
		}
	}

	st.PrintLn("Added external app (ID: " +
		strconv.Itoa(int(appId)) + ", Command: " +
		newExternalApp.CommandLine + ")")
}

func (st *State) commandAppPing(args []string) {
	app.At(cp, "commandAppPing")

	if len(args) < 1 {
		st.PrintError("No task id passed! e.g. ping 1")
		return
	}

	passedID, err := strconv.Atoi(args[0])
	if err != nil {
		st.PrintError("Task id must be an integer.")
		return
	}

	client, ok := signal.GetClient(uint(passedID))
	if !ok {
		st.PrintError("Task with given id is not running.")
		return
	}

	start := time.Now()
	_, err = client.Ping()
	if err != nil {
		st.PrintError(err.Error())
		return
	}
	st.PrintLn(fmt.Sprintf("ping time %s", time.Now().Sub(start).String()))
}

func (st *State) commandShutDown(args []string) {
	app.At(cp, "commandShutDown")

	if len(args) < 1 {
		st.PrintError("No task id passed! e.g. shutdown 1")
		return
	}

	passedID, err := strconv.Atoi(args[0])
	if err != nil {
		st.PrintError("Task id must be an integer.")
		return
	}

	client, ok := signal.GetClient(uint(passedID))
	if !ok {
		st.PrintError("Task with given id is not running.")
		return
	}

	resp, err := client.Shutdown()
	if err != nil {
		st.PrintError(err.Error())
		return
	}
	st.PrintLn(fmt.Sprintf("shutdown pid %d", resp.Pid))
}

func (st *State) commandResourceUsage(args []string) {
	app.At(cp, "commandResourceUsage")
	if len(args) < 1 {
		st.PrintError("No task id passed! e.g. res_usage 1")
		return
	}

	passedID, err := strconv.Atoi(args[0])
	if err != nil {
		st.PrintError("Task id must be an integer.")
		return
	}

	client, ok := signal.GetClient(uint(passedID))
	if !ok {
		st.PrintError("Task with given id is not running.")
		return
	}

	resp, err := client.Top()
	if err != nil {
		st.PrintError(err.Error())
		return
	}
	stats := resp.MemStats
	st.PrintLn(fmt.Sprintf("MemStats HeapAlloc:%d HeapSys:%d HeapIdle:%d HeapInuse:%d StackSys:%d StackInuse:%d",
		stats.HeapAlloc, stats.HeapSys, stats.HeapIdle, stats.HeapInuse, stats.StackSys, stats.StackInuse))
}

func (st *State) commandAttach(args []string) {
	app.At(cp, "commandAttach")

	if len(args) < 1 {
		st.PrintError("No task id passed! e.g. attach 1")
		return
	}

	passedID, err := strconv.Atoi(args[0])
	if err != nil {
		st.PrintError("Task id must be an integer.")
		return
	}

	eaId := msg.ExtAppId(passedID)

	ea, err := hypervisor.GetExternalApp(eaId)
	if err != nil {
		st.PrintError(err.Error())
		return
	}

	st.PrintLn(ea.GetFullCommandLine())
	err = st.task.AttachExternalApp(ea)
	if err != nil {
		st.PrintError(err.Error())
	}
}

func (st *State) commandListRunningExternalApps(args []string) {
	app.At(cp, "commandListRunningExternalApps")

	extApps := hypervisor.ExternalAppListGlobal.TaskMap
	if len(extApps) == 0 {
		st.PrintLn("No external apps running.\n" +
			"Try starting one with \"start\" command (\"help\" or \"h\" for help).")
		return
	}

	fullPrint := false

	if len(args) > 0 && args[0] == "-f" {
		fullPrint = true
	}

	for id, extApp := range extApps {
		appCmd := ""

		if fullPrint {
			appCmd = extApp.GetFullCommandLine()
		} else {
			appCmd = strings.Split(extApp.GetFullCommandLine(), " ")[0]
		}

		st.Printf("[ %d ] -> [ %s ]\n", int(id), appCmd)
	}
}

func (st *State) commandCloseTerminalFirstStage(args []string) {
	if len(args) == 1 {
		//handle storedTerminalIds errors
		if /****/ len(st.storedTerminalIds) < 1 {
			st.PrintError("Use 'list_terms' command to see their IDs")
			return
		} else if len(st.storedTerminalIds) == 1 {
			st.PrintError("Shouldn't close when only 1 terminal remains (UNTIL GUI IS MADE)")
			return
		}

		//handle arg conversion errors
		storedId, err := strconv.Atoi(args[0])
		if err != nil {
			st.PrintError("Unable to convert passed index.")
			s := "err.Error(): \"" + err.Error() + "\""
			st.PrintError(s)
			return
		}

		//handle index range errors
		if storedId < 0 ||
			storedId >= len(st.storedTerminalIds) {
			st.PrintError("Index not in range.")
			return
		}

		//everything should be valid here
		st.SendCommand("close_term", []string{strconv.Itoa(int(st.storedTerminalIds[storedId]))})
	} else { //args failure (too many/few passed)
		st.PrintError("Must supply ONE valid ID argument")
		//IDEALLY we'd use ANY VALID index at position 0,
		//but I think you'd want us to prioritize simplicity
		//over nitpickiness like that.
		return
	}
}
