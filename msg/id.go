package msg

import (
	"math/rand"
)

//HyperVisor: processId
type ProcessId uint64

var ProcessIdGlobal ProcessId = 1 //sequential

func RandProcessId() ProcessId {
	ProcessIdGlobal += 1
	return ProcessIdGlobal
	//return (ProccesId)(rand.Int63())
}

//terminate ID
type TerminalId uint64

func RandTerminalId() TerminalId {
	return (TerminalId)(rand.Int63())
}
