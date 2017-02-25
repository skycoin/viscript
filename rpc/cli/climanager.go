package cli

import (
	"github.com/corpusc/viscript/msg"
	tm "github.com/corpusc/viscript/rpc/terminalmanager"
)

type CliManager struct {
	Commands         map[string]func(args []string) error
	ChosenTerminalId msg.TerminalId
	ChosenProcessId  msg.ProcessId
	CachedIds        []msg.TermAndAttachedProcessID

	Client     *tm.RPCClient
	SessionEnd bool
}

func (c *CliManager) Init(port string) {
	c.initCommands()

	c.CachedIds = make([]msg.TermAndAttachedProcessID, 0)
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
		println("Command not found. Type 'help' or 'h'.")
	}
}

func (c *CliManager) initCommands() {
	c.Commands = map[string]func(args []string) error{}
	c.Commands["help"] = c.PrintHelp
	c.Commands["h"] = c.PrintHelp

	c.Commands["lt"] = c.ListTerminalIDs
	c.Commands["ltp"] = c.ListTermIDsWithAttachedProcesses
	c.Commands["lp"] = c.ListProcessIDs
	c.Commands["stp"] = c.StartTerminalWithProcess
	// TODO: c.commands["cinf"] = c.GetChannelInfo

	c.Commands["clear"] = c.ClearTerminal
	c.Commands["quit"] = c.Quit
	c.Commands["q"] = c.Quit
}

func (c *CliManager) TerminalIdNotNil() bool {
	return c.ChosenProcessId != 0
}

func (c *CliManager) ProcessIdNotNil() bool {
	return c.ChosenProcessId != 0
}
