package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/corpusc/viscript/hypervisor"
)

type rpcMessage struct {
	Command   string
	Arguments []string
}

func main() {
	// port := "7777"
	// if len(os.Args) >= 2 {
	// 	port = os.Args[1]
	// }
	// rpcClient := hypervisor.RunClient(":" + port)
	// promptCycle(rpcClient)
}

func promptCycle(rpcClient *hypervisor.RPCClient) {
	for {
		if commandDispatcher(rpcClient) {
			break
		}
	}
}

func commandDispatcher(rpcClient *hypervisor.RPCClient) bool {
	command, _ := cliInput("Enter the command (help for commands list):\n> ")

	if command == "" {
		return false
	}

	command = strings.ToLower(command)

	switch command {
	case "exit", "quit":
		fmt.Println("\nGoodbye\n")
		return true

	case "help":
		printHelp()

	default:
		fmt.Println("\nUnknown command: %s, type 'help' to get the list of available commands.\n\n", command)
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
	fmt.Println("\n<<[ HELP ]>>\n")
	fmt.Println("> exit (or q)\t\tcloses the terminal.\n")
}

func errorOut(err error) {
	fmt.Println("Error. Server says:", err)
}
