package process

import (
	"strconv"
	"strings"

	"github.com/corpusc/viscript/app"
	"github.com/corpusc/viscript/hypervisor"
	extTask "github.com/corpusc/viscript/hypervisor/task_ext"
	"github.com/corpusc/viscript/msg"
)

const cp = "hypervisor/process/terminal/commands"

func (st *State) commandHelp() {
	app.At(cp, "PrintHelp")

	st.PrintLn(app.GetBarOfChars("-", int(st.VisualInfo.NumColumns)))
	//st.PrintLn("Current commands:")
	st.PrintLn("help:                  This message ('?' or 'h' for short).")
	st.PrintLn("new_terminal:          Add new terminal.")
	st.PrintLn("start (-a) <command>:  Start external task. (pass -a to also attach).")
	st.PrintLn("attach   <id>:         Attach external task with given id to terminal.")
	st.PrintLn("shutdown <id>:         [TODO] Shutdown external task with given id.")
	st.PrintLn("ls (-f):               List external tasks (pass -f for full commands).")
	st.PrintLn("rpc:                   Issues command: \"go run rpc/cli/cli.go\"")
	//st.PrintLn("Current hotkeys:")
	st.PrintLn("CTRL+Z:                Detach currently attached process.")
	//st.PrintLn("    CTRL+C:           ___description goes here___")
	st.PrintLn(app.GetBarOfChars("-", int(st.VisualInfo.NumColumns)))
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

	newExtProc, err := extTask.MakeNewTaskExternal(args, detached)
	if err != nil {
		st.PrintError(err.Error())
		return
	}

	err = newExtProc.Start()
	if err != nil {
		st.PrintError(err.Error())
		return
	}

	extProcInterface := newExtProc.GetExtProcessInterface()
	procId := hypervisor.AddExtProcess(extProcInterface)

	if !detached {
		err = st.proc.AttachExternalProcess(extProcInterface)
		if err != nil {
			st.PrintError(err.Error())
		}
	}

	st.PrintLn("Added External Process (ID: " +
		strconv.Itoa(int(procId)) + ", Command: " +
		newExtProc.CommandLine + ")")

}

func (st *State) commandShutDown(args []string) {
	app.At(cp, "commandShutDown")
	println("TODO: finish implementing command shutdown")
	println("TODO: finish implementing command shutdown")
	println("TODO: finish implementing command shutdown")
}

func (st *State) commandAttach(args []string) {
	app.At(cp, "commandAttach")

	if len(args) < 1 {
		st.PrintError("No task id passed! eg: attach 1")
		return
	}

	passedId, err := strconv.Atoi(args[0])
	if err != nil {
		st.PrintError("Task id must be an integer.")
		return
	}

	extProcId := msg.ExtProcessId(passedId)

	extProc, err := hypervisor.GetExtProcess(extProcId)
	if err != nil {
		st.PrintError(err.Error())
		return
	}

	st.PrintLn(extProc.GetFullCommandLine())
	err = st.proc.AttachExternalProcess(extProc)
	if err != nil {
		st.PrintError(err.Error())
	}
}

func (st *State) commandListExternalTasks(args []string) {
	app.At(cp, "commandListExternalTasks")

	extTaskMap := hypervisor.ExtProcessListGlobal.ProcessMap
	if len(extTaskMap) == 0 {
		st.PrintLn("No external tasks running.\n" +
			"Try starting one with \"start\" command (\"help\" or \"h\" for help).")
		return
	}

	fullPrint := false

	if len(args) > 0 && args[0] == "-f" {
		fullPrint = true
	}

	for procId, extProc := range extTaskMap {
		procCommand := ""

		if fullPrint {
			procCommand = extProc.GetFullCommandLine()
		} else {
			procCommand = strings.Split(
				extProc.GetFullCommandLine(), " ")[0]
		}

		st.Printf("[ %d ] -> [ %s ]\n", int(procId), procCommand)
	}
}
