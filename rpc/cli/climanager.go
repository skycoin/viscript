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

	c.setCommand("ltp", c.ListTermIDsWithAttachedProcesses)
	c.setCommand("lp", c.ListProcessIDs)
	c.setCommand("sett", c.SetDefaultTerminalId)
	c.setCommand("setp", c.SetDefaultProcessId)

	c.setCommand("cft", c.ShowChosenTermChannelInfo)

	c.setCommand("stp", c.StartTerminalWithProcess)

	c.setCommandWithShortcut("help", c.PrintHelp)
	c.setCommandWithShortcut("clear", c.ClearTerminal)
	c.setCommandWithShortcut("quit", c.Quit)
}

func (c *CliManager) setCommand(command string, f func(args []string) error) {
	if !c.commandExists(command) {
		c.Commands[command] = f
	}
}

func (c *CliManager) setCommandWithShortcut(command string, f func(args []string) error) {
	c.setCommand(command, f)
	c.setCommand(string(command[0]), f) // set first char as shortcut
}

func (c *CliManager) commandExists(key string) bool {
	_, ok := c.Commands[key]
	return ok
}

func (c *CliManager) PrintServerError(err error) {
	c.Client.ErrorOut(err)
}

func (c *CliManager) TerminalIdNotNil() bool {
	return c.ChosenTerminalId != 0
}

func (c *CliManager) ProcessIdNotNil() bool {
	return c.ChosenProcessId != 0
}
