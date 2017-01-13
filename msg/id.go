package msg

import (
//"math/rand"
)

//HyperVisor: processId
type ProcessId uint64

var ProcessIdGlobal ProcessId = 0 //sequential

func RandProcessId() ProcessId {
	ProcessIdGlobal += 1
	return ProcessIdGlobal
	//return (ProccesId)(rand.Int63())
}
