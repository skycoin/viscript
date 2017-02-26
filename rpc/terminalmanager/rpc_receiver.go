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

func (receiver *RPCReceiver) ListTIDsWithProcessIDs(_ []string, result *[]byte) error {
	terms := receiver.TerminalManager.terminalStack.Terms
	termsWithProcessIDs := make([]msg.TermAndAttachedProcessID, 0)

	for termID, term := range terms {
		termsWithProcessIDs = append(termsWithProcessIDs,
			msg.TermAndAttachedProcessID{TerminalId: termID, AttachedProcessId: term.AttachedProcess})
	}
	fmt.Printf("Terms with process IDs list:%+v", termsWithProcessIDs)
	*result = msg.Serialize((uint16)(0), termsWithProcessIDs)
	return nil
}

// TODO: here, finish this
// func (receiver *RPCReceiver) GetChannelInfo(args []string, result *[]byte) error {
// 	println("rpc_receiver.go/GetChannelInfo()")
// 	return nil
// }

func (receiver *RPCReceiver) StartTerminalWithProcess(_ []string, result *[]byte) error {
	terms := receiver.TerminalManager.terminalStack
	newTerminalID := terms.AddTerminal()
	fmt.Println("Terminal with ID", newTerminalID, "created!")
	*result = msg.Serialize((uint16)(0), newTerminalID)
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
