package climanager

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/skycoin/viscript/hypervisor/dbus"
	"github.com/skycoin/viscript/msg"
	tm "github.com/skycoin/viscript/rpc/terminalmanager"
)

func (c *CliManager) PrintHelp(_ []string) error {
	p := fmt.Printf
	p("\n<< [- HELP -] >>\n\n")

	p("> stp\t\tStart a new terminal with process.\n\n")

	p("> ltp\t\tList terminal Ids with Attached Process Ids.\n")
	p("> lp\t\tList process Ids with labels.\n\n")

	p("> sett <tId>\tSet given terminal Id as default for all following commands.\n")
	p("> setp <pId>\tSet given process Id as default for all following commands.\n\n")

	p("> cft\t\tGet out channel info of terminal with default Id.\n\n")

	p("> clear(c)\tClear the terminal.\n")
	p("> quit(q)\tQuit from cli.\n\n")

	return nil
}

func (c *CliManager) Quit(_ []string) error {
	println("See ya again! :>")
	c.SessionEnd = true
	return nil
}

func (c *CliManager) ClearTerminal(_ []string) error {

	ros := runtime.GOOS

	if ros == "linux" || ros == "darwin" {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else if ros == "windows" {
		cmd := exec.Command("cmd", "/C", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		println("Your platform is unsupported! I can't clear terminal screen :(.")
	}

	return nil
}

func (c *CliManager) ListTermIDsWithAttachedTasks(_ []string) error {
	termsWithTaskIDs, err := GetTermIDsWithTaskIDs(c.Client)

	if err != nil {
		return err
	}

	fmt.Printf("Terminals (%d) defaults marked with {}:\n", len(termsWithTaskIDs))
	fmt.Println("\nIdx\tTerminal Id\t\tAttached Process Id")
	for index, term := range termsWithTaskIDs {
		fmt.Printf("[ %d ]\t", index)

		// mark selected default terminal id
		if term.TerminalId == c.ChosenTerminalId {
			fmt.Printf("{ %d }\t", term.TerminalId)
		} else {
			fmt.Printf("  %d\t", term.TerminalId)
		}

		// mark selected default process id
		if term.AttachedTaskId == c.ChosenTaskId {
			fmt.Printf("{ %d }\t", term.AttachedTaskId)
		} else {
			fmt.Printf("  %d\t", term.AttachedTaskId)
		}
		fmt.Printf("\n")
	}
	println()

	return nil
}

func (c *CliManager) ListTasks(_ []string) error {
	taskInfos, err := GetTasks(c.Client)
	if err != nil {
		return err
	}

	fmt.Printf("Tasks (%d) default marked with {}:\n", len(taskInfos))
	fmt.Println("\nIdx\t Id\t Type\t\t Label")
	for index, processInfo := range taskInfos {
		if processInfo.Id == c.ChosenTaskId {
			fmt.Printf("[ %d ]\t{ %6d } %6d \t%s\n", index, processInfo.Id, processInfo.Type, processInfo.Label)
		} else {
			fmt.Printf("[ %d ]\t  %6d   %6d \t%s\n", index, processInfo.Id, processInfo.Type, processInfo.Label)
		}
	}
	println()
	return nil
}

func (c *CliManager) SetDefaultTerminalId(args []string) error {
	if len(args) == 0 {
		fmt.Printf("\n\nPass the terminal Id as argument please.")
		return nil
	}

	termId, err := strconv.Atoi(args[0])
	if err != nil || termId < 1 {
		fmt.Printf("\n\nArgument should be a number > 0, not %s\n\n", args[0])
		return nil
	}

	c.ChosenTerminalId = msg.TerminalId(termId)
	return nil
}

func (c *CliManager) SetDefaultTaskId(args []string) error {
	if len(args) == 0 {
		fmt.Printf("\n\nArgument should be a number > 0, not %s\n\n", args[0])
		return nil
	}

	taskId, err := strconv.Atoi(args[0])
	if err != nil || taskId < 1 {
		fmt.Printf("\n\nArgument should be a number > 0, not %s\n\n", args[0])
	}

	c.ChosenTaskId = msg.TaskId(taskId)
	return nil
}

func (c *CliManager) ShowChosenTermChannelInfo(_ []string) error {
	if c.ChosenTerminalId == 0 {
		fmt.Printf("\nDefault terminal Id is not set.\n\n")
		return nil
	}

	response, err := c.Client.SendToRPC("GetTermChannelInfo", []string{fmt.Sprintf("%d", c.ChosenTerminalId)})
	if err != nil {
		return err
	}

	var channelInfo msg.ChannelInfo
	err = msg.Deserialize(response, &channelInfo)
	if err != nil {
		return err
	}

	fmt.Printf("Term (Id: %d) out channel info:\n", c.ChosenTerminalId)

	println("Channel Id:", channelInfo.ChannelId)
	println("Channel Owner:", channelInfo.Owner)
	println("Channel Owner's Type:", dbus.ResourceTypeNames[channelInfo.OwnerType])
	println("Channel ResourceIdentifier:", channelInfo.ResourceIdentifier)

	subCount := len(channelInfo.Subscribers)

	if subCount == 0 {
		fmt.Printf("No subscribers to this channel.\n")
	} else {
		fmt.Printf("Channel's Subscribers (%d total):\n\n", subCount)
		fmt.Println("Index\tResourceId\t\tResource Type")
		for index, subscriber := range channelInfo.Subscribers {
			fmt.Println(index, "\t", subscriber.SubscriberId, "\t\t",
				dbus.ResourceTypeNames[subscriber.SubscriberType])
		}
	}

	return nil
}

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

func GetTerminalIDs(client *tm.RPCClient) ([]msg.TerminalId, error) {
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

func GetTermIDsWithTaskIDs(client *tm.RPCClient) ([]msg.TermAndAttachedTaskId, error) {
	response, err := client.SendToRPC("ListTIDsWithTaskIDs", []string{})
	if err != nil {
		return []msg.TermAndAttachedTaskId{}, err
	}

	var termsAndAttachedTasks []msg.TermAndAttachedTaskId
	err = msg.Deserialize(response, &termsAndAttachedTasks)
	if err != nil {
		return []msg.TermAndAttachedTaskId{}, err
	}
	return termsAndAttachedTasks, nil
}

func GetTasks(client *tm.RPCClient) ([]msg.TaskInfo, error) {
	response, err := client.SendToRPC("ListTasks", []string{})
	if err != nil {
		return []msg.TaskInfo{}, err
	}

	var taskInfos []msg.TaskInfo
	err = msg.Deserialize(response, &taskInfos)
	if err != nil {
		return []msg.TaskInfo{}, err
	}
	return taskInfos, nil
}
