package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/corpusc/viscript/msg"
	tm "github.com/corpusc/viscript/rpc/terminalmanager"
	"os/exec"
	"runtime"
)

type rpcMessage struct {
	Command   string
	Arguments []string
}

// TODO:? status of the terminal more info with id

var clear map[string]func() // cross platform clear command

func init() {
	clear = make(map[string]func())
	clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

	clear["windows"] = func() {
		cmd := exec.Command("cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

	clear["darwin"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func main() {
	port := "7777"
	if len(os.Args) >= 2 {
		port = os.Args[1]
	}
	fmt.Println(port)
	rpcClient := tm.RunClient(":" + port)
	promptCycle(rpcClient)
}

func promptCycle(rpcClient *tm.RPCClient) {
	ended := false
	for !ended {
		ended = commandDispatcher(rpcClient)
	}
}

func commandDispatcher(rpcClient *tm.RPCClient) bool {
	command, _ := cliInput("Enter the command (help(h) for commands list):\n> ")

	if command == "" {
		return false
	}

	command = strings.ToLower(command)

	switch command {

	case "lt":
		listTerminalIDs(rpcClient)

	case "ltp":
		listTermIDsWithAttachedProcesses(rpcClient)

	case "lp":
		listProcessIDs(rpcClient)

	case "ld":
		listDbusObjects(rpcClient)

	case "lsub":
		listSubscribers(rpcClient)

	case "lpub":
		listPublishers(rpcClient)

	case "stp":
		startTerminalWithProcess(rpcClient)

	case "help", "h":
		printHelp()

	case "clear":
		clearTerminal()

	case "quit", "q":
		fmt.Println("\nGoodbye")
		return true

	default:
		fmt.Printf("\nUnknown command: %s, type "+
			"'help' or 'h' for available commands.\n\n", command)
	}
	return false
}

func cliInput(prompt string) (command string, args []string) {
	fmt.Print(prompt)
	command = ""
	args = []string{}
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	splitted := strings.Fields(input)
	if len(splitted) == 0 {
		return
	}
	command = strings.Trim(splitted[0], " ")
	if len(splitted) > 1 {
		args = splitted[1:]
	}
	return
}

func printHelp() {
	fmt.Println("\n<< -[ HELP ]- >>")
	fmt.Println()
	fmt.Println("> lt\t\tList all terminal IDs.")
	fmt.Println("> ltp\t\tList all terminal IDs with Attached Process IDs.")
	fmt.Println("> lp\t\tList all process IDs.")
	fmt.Println("> ld\t\tList all dbus objects. --TODO")
	fmt.Println("> lsub\t\tList all subscribers. --TODO")
	fmt.Println("> lpub\t\tList all publishers. --TODO")
	fmt.Println()
	fmt.Println("> stp\t\tStart a new terminal with process.")
	fmt.Println()
	fmt.Println("> clear\t\tClear the terminal.")
	fmt.Println("> quit(q)\tQuit from cli.")
	fmt.Println()
}

func clearTerminal() {
	value, ok := clear[runtime.GOOS]
	if ok {
		value()
	} else {
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}

func errorOut(err error) {
	fmt.Println("Error. Server says:", err)
}

func listTerminalIDs(client *tm.RPCClient) {
	termIDs, err := getTerminalIDs(client)
	if err != nil {
		errorOut(err)
		return
	}

	fmt.Printf("\nTerminals(%d total):\n\n", len(termIDs))
	fmt.Println("Num\tID")
	fmt.Println()
	for i, termID := range termIDs {
		fmt.Println(i, "\t", termID)
	}
}

func listTermIDsWithAttachedProcesses(client *tm.RPCClient) {
	termsWithProcessIDs, err := getTermIDsWithProcessIDs(client)
	if err != nil {
		errorOut(err)
		return
	}

	fmt.Printf("\nTerminals(%d total):\n\n", len(termsWithProcessIDs))
	fmt.Println("Index\tTerminalID\t\tAttached Process ID")
	fmt.Println()
	for i, term := range termsWithProcessIDs {
		fmt.Printf("%d\t%d\t%d\n", i, term.TerminalId, term.AttachedProcessId)
	}
}

func listProcessIDs(client *tm.RPCClient) {
	processIDs, err := getProcessIDs(client)
	if err != nil {
		errorOut(err)
		return
	}

	fmt.Printf("\nProcesses(%d total):\n\n", len(processIDs))
	fmt.Println("Num\tID")
	fmt.Println()
	for i, processID := range processIDs {
		fmt.Println(i, "\t", processID)
	}
}

func listDbusObjects(client *tm.RPCClient) {
	fmt.Println("listDbusIDs()")

}

func listSubscribers(client *tm.RPCClient) {
	fmt.Println("listSubscribers()")
}

func listPublishers(client *tm.RPCClient) {
	fmt.Println("listPublishers()")
}

func startTerminalWithProcess(client *tm.RPCClient) {
	fmt.Println("startTerminalWithProcess()")
	response, err := client.SendToRPC("StartTerminalWithProcess", []string{})
	if err != nil {
		errorOut(err)
		return
	}

	var newID msg.TerminalId
	err = msg.Deserialize(response, &newID)
	if err != nil {
		errorOut(err)
		return
	}

	fmt.Println("New terminal was created with ID", newID)
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
