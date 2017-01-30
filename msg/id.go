package msg

import (
	"math/rand"
)

//ProcessId - HyperVisor: processId
type ProcessId uint64

//TerminalId - Terminal: TerminalId
type TerminalId uint64

//ProcessIdGlobal - Global sequential process id
var ProcessIdGlobal ProcessId = 1

func NextProcessId() ProcessId {
	ProcessIdGlobal += 1
	return ProcessIdGlobal
	//return (ProccesId)(rand.Int63())
}

func RandTerminalId() TerminalId {
	return (TerminalId)(rand.Int63())
}
