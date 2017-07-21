package msg

import (
	"math/rand"
)

type TaskId uint64 //HyperVisor: TaskId
type TerminalId uint64
type ExtTaskId uint64

var TaskIdGlobal TaskId = 1 //sequential
var ExtTaskIdGlobal ExtTaskId = 1

func NextTaskId() TaskId {
	TaskIdGlobal += 1
	return TaskIdGlobal
}

func NextExtTaskId() ExtTaskId {
	ExtTaskIdGlobal += 1
	return ExtTaskIdGlobal
}

func RandTerminalId() TerminalId {
	return (TerminalId)(rand.Int63())
}
