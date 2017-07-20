package main

import (
	"github.com/skycoin/viscript/signal"
	"github.com/skycoin/viscript/msg"
	"strings"
	"bufio"
	"os"
	"log"
	"fmt"
)

func main() {
	signal.Init("127.0.0.1:7999").Run()
	showHelp()
	promptCycle()
}

func promptCycle() {
	for {
		newCommand, args := inputFromCli()
		if newCommand == "" {
			continue
		}
		dispatcher(strings.ToLower(newCommand), args)
	}
}

func inputFromCli() (command string, args []string) {
	command = ""
	args = []string{}

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()

	splitInput := strings.Fields(input)
	if len(splitInput) == 0 {
		return
	}

	command = strings.Trim(splitInput[0], " ")
	if len(splitInput) > 1 {
		args = splitInput[1:]
	}
	return
}

func dispatcher(cmd string, args []string) {
	log.Println("new command:" + cmd)


	switch cmd {

	case "help":
		showHelp()

	case "ping":
		msgUserCommand := msg.MessageUserCommand{
			Sequence: 1,
			AppId:    2,
			Payload:  msg.Serialize(msg.TypePing, msg.MessagePing{})}

		serializedCommand := msg.Serialize(msg.TypeUserCommand, msgUserCommand)

		signal.Monitor.Send(2, serializedCommand)

	case "shutdown":
		msgUserCommand := msg.MessageUserCommand{
			Sequence: 1,
			AppId:    2,
			Payload:  msg.Serialize(msg.TypeShutdown, msg.MessagePing{})}

		serializedCommand := msg.Serialize(msg.TypeUserCommand, msgUserCommand)

		signal.Monitor.Send(2, serializedCommand)

	case "res_usage":
		msgUserCommand := msg.MessageUserCommand{
			Sequence: 1,
			AppId:    2,
			Payload:  msg.Serialize(msg.TypeResourceUsage, msg.MessagePing{})}

		serializedCommand := msg.Serialize(msg.TypeUserCommand, msgUserCommand)

		signal.Monitor.Send(2, serializedCommand)


	default:
		log.Println("Unknown user command:")

	}

}

func showHelp() {
	fmt.Printf("> help\t\t\tShow list of commands.\n\n")
	fmt.Printf("> ping <id>\t\tPing app with choosen id.\n\n")
	fmt.Printf("> shutdown <id>\t\tKill app with choosen id.\n\n")
	fmt.Printf("> res_usage <id>\tShow cpu and memory stats.\n\n")
}
