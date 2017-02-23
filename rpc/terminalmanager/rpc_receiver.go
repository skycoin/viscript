package terminalmanager

import (
	"fmt"
	"github.com/corpusc/viscript/msg"
)

type RPCReceiver struct {
	TerminalManager *TerminalManager
}

func (receiver *RPCReceiver) ListTerminalIDs(_ []string, result *[]byte) error {
	terms := receiver.TerminalManager.terminalStack.Terms
	terminalIDs := make([]msg.TerminalId, 0)

	for k, _ := range terms {
		terminalIDs = append(terminalIDs, k)
	}
	fmt.Println("Terminal ID list:", terminalIDs)
	*result = msg.Serialize((uint16)(0), terminalIDs)
	return nil
}

func (receiver *RPCReceiver) ListProcessIDs(_ []string, result *[]byte) error {
	processes := receiver.TerminalManager.processList.ProcessMap
	processIDs := make([]msg.ProcessId, 0)

	for k, _ := range processes {
		processIDs = append(processIDs, k)
	}
	fmt.Println("Process ID list:", processIDs)
	*result = msg.Serialize((uint16)(0), processIDs)
	return nil
}
