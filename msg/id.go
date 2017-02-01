package msg

import (
	"math/rand"
)

type ProcessId uint64 //HyperVisor: processId
type TerminalId uint64

var ProcessIdGlobal ProcessId = 1 //sequential

func NextProcessId() ProcessId {
	ProcessIdGlobal += 1
	return ProcessIdGlobal
}

func RandTerminalId() TerminalId {
	return (TerminalId)(rand.Int63())
}
