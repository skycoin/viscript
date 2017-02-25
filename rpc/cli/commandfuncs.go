package cli

import (
	"fmt"
	"github.com/corpusc/viscript/msg"
	tm "github.com/corpusc/viscript/rpc/terminalmanager"
	"os"
	"os/exec"
	"runtime"
)

func (c *CliManager) PrintServerError(err error) {
	c.Client.ErrorOut(err)
}

func (c *CliManager) PrintHelp(_ []string) error {
	p := fmt.Printf
	p("\n<< -[ HELP ]- >>\n\n")

	p("> lt\t\tList all terminal IDs.\n")
	p("> ltp\t\tList all terminal IDs with Attached Process IDs.\n")
	p("> lp\t\tList all process IDs.\n")
	// p("> ld\t\tList all dbus objects. --TODO\n")
	p("> cinf <Id>\tGet channel info of terminal with Id. --TODO\n\n")
	// p("> lpub\t\tList all publishers. --TODO\n\n")

	p("> stp\t\tStart a new terminal with process.\n\n")

	p("> clear\t\tClear the terminal.\n")
	p("> quit(q)\tQuit from cli.\n\n")

	return nil
}

func (c *CliManager) Quit(_ []string) error {
	c.SessionEnd = false
	return nil
}

func (c *CliManager) ClearTerminal(_ []string) error {

	runtimeOs := runtime.GOOS

	if runtimeOs == "linux" || runtimeOs == "darwin" {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else if runtimeOs == "windows" {
		cmd := exec.Command("cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		println("Your platform is unsupported! I can't clear terminal screen :(.")
	}

	return nil
}

func (c *CliManager) ListTerminalIDs(_ []string) error {
	termIDs, err := getTerminalIDs(c.Client)
	if err != nil {
		return err
	}

	fmt.Printf("\nTerminals(%d total):\n\n", len(termIDs))
	fmt.Println("Num\tID")
	fmt.Println()
	for i, termID := range termIDs {
		fmt.Println(i, "\t", termID)
	}

	return nil
}

func (c *CliManager) ListTermIDsWithAttachedProcesses(_ []string) error {
	termsWithProcessIDs, err := getTermIDsWithProcessIDs(c.Client)
	if err != nil {
		return err
	}

	fmt.Printf("\nTerminals(%d total):\n\n", len(termsWithProcessIDs))
	fmt.Println("Index\tTerminalID\t\tAttached Process ID")
	fmt.Println()
	for i, term := range termsWithProcessIDs {
		fmt.Printf("%d\t%d\t%d\n", i, term.TerminalId, term.AttachedProcessId)
	}

	return nil
}

func (c *CliManager) ListProcessIDs(_ []string) error {
	processIDs, err := getProcessIDs(c.Client)
	if err != nil {
		return err
	}

	fmt.Printf("\nProcesses(%d total):\n\n", len(processIDs))
	fmt.Println("Num\tID")
	fmt.Println()
	for i, processID := range processIDs {
		fmt.Println(i, "\t", processID)
	}

	return nil
}

// func (c *CliManager) ListDbusObjects(client *tm.RPCClient) {
// 	fmt.Println("listDbusIDs()")

// }

// func (c *CliManager) GetChannelInfo(args []string) error {
// 	fmt.Println("commandfuncs.go/GetChannelInfo()")
// 	// TODO: implement this
// 	if len(args) == 0 {
// 		fmt.Printf("\n\nPass the terminal Id as argument please.")
// 		return
// 	}

// 	termId, err := strconv.Atoi(args[0])
// 	// FIXME: we should save i guess the list of terminal IDS and not send a request if its
// 	// incorrect. So we need a terminal id list and from there should cli user choose the id.
// 	if err != nil || termId < 1 {
// 		fmt.Printf("\nArgument should be a number > 0, not %s\n\n", args[0])
// 	}

// 	response, err := c.client.SendToRPC("GetChannelInfo", termId)
// 	if err != nil {
// 		errorOut(err)
// 		return
// 	}

// 	var channelInfo msg.ChannelInfo
// 	err = msg.Deserialize(response, &channelInfo)
// 	if err != nil {
// 		errorOut(err)
// 		return
// 	}

// 	// TODO: print structured channel info

// }

func (c *CliManager) StartTerminalWithProcess(_ []string) error {
	fmt.Println("startTerminalWithProcess()")
	response, err := c.Client.SendToRPC("StartTerminalWithProcess", []string{})
	if err != nil {
		return err
	}

	var newID msg.TerminalId
	err = msg.Deserialize(response, &newID)
	if err != nil {
		return err
	}

	fmt.Println("New terminal was created with ID", newID)

	return nil
}

func getTerminalIDs(client *tm.RPCClient) ([]msg.TerminalId, error) {
	response, err := client.SendToRPC("ListTerminalIDs", []string{})
	if err != nil {
		return []msg.TerminalId{}, err
	}

	var termIDs []msg.TerminalId
	err = msg.Deserialize(response, &termIDs)
	if err != nil {
		return []msg.TerminalId{}, err
	}
	return termIDs, nil
}

func getTermIDsWithProcessIDs(client *tm.RPCClient) ([]msg.TermAndAttachedProcessID, error) {
	response, err := client.SendToRPC("ListTIDsWithProcessIDs", []string{})
	if err != nil {
		return []msg.TermAndAttachedProcessID{}, err
	}

	var termsAndAttachedProcesses []msg.TermAndAttachedProcessID
	err = msg.Deserialize(response, &termsAndAttachedProcesses)
	if err != nil {
		return []msg.TermAndAttachedProcessID{}, err
	}
	return termsAndAttachedProcesses, nil
}

func getProcessIDs(client *tm.RPCClient) ([]msg.ProcessId, error) {
	response, err := client.SendToRPC("ListProcessIDs", []string{})
	if err != nil {
		return []msg.ProcessId{}, err
	}

	var processIDs []msg.ProcessId
	err = msg.Deserialize(response, &processIDs)
	if err != nil {
		return []msg.ProcessId{}, err
	}
	return processIDs, nil
}
