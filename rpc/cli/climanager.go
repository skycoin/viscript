package cli

import (
	"github.com/corpusc/viscript/msg"
	tm "github.com/corpusc/viscript/rpc/terminalmanager"
)

type CliManager struct {
	Commands         map[string]func(args []string) error
	ChosenTerminalId msg.TerminalId
	ChosenProcessId  msg.ProcessId

	Client     *tm.RPCClient
	SessionEnd bool
}

func (c *CliManager) Init(port string) {
	c.initCommands()

	c.SessionEnd = false
	c.Client = tm.RunClient(":" + port)
}

func (c *CliManager) CommandDispatcher(command string, args []string) {
	runFunc, commandExists := c.Commands[command]
	if commandExists {
		serverError := runFunc(args)
		if serverError != nil {
			c.PrintServerError(serverError)
		}
	} else {
		println("Command not found. Type 'help(h)' for commands.")
	}
}

func (c *CliManager) initCommands() {
	c.Commands = map[string]func(args []string) error{}

	c.Commands["ltp"] = c.ListTermIDsWithAttachedProcesses
	c.Commands["sett"] = c.SetDefaultTerminalId
	c.Commands["setp"] = c.SetDefaultProcessId

	c.Commands["cft"] = c.ShowChosenTermChannelInfo

	c.Commands["stp"] = c.StartTerminalWithProcess

	c.Commands["help"] = c.PrintHelp
	c.Commands["h"] = c.PrintHelp
	c.Commands["clear"] = c.ClearTerminal
	c.Commands["c"] = c.ClearTerminal
	c.Commands["quit"] = c.Quit
	c.Commands["q"] = c.Quit
}

func (c *CliManager) PrintServerError(err error) {
	c.Client.ErrorOut(err)
}

func (c *CliManager) TerminalIdNotNil() bool {
	return c.ChosenProcessId != 0
}

func (c *CliManager) ProcessIdNotNil() bool {
	return c.ChosenProcessId != 0
}
