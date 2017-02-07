package process

import (
	"github.com/corpusc/viscript/msg"
)

func PutChar(out chan []byte, char uint32) {
	println("(process/terminal/api.go).PutChar()")

	msg.SerializeAndDispatch(
		out,
		msg.TypePutChar,
		msg.MessagePutChar{TermId: 0, Char: char})
}

//sending messages back to hypervisor to set terminal
func SetChar(out chan []byte, x uint32, y uint32, char uint32) {
	println("(process/terminal/api.go).SetChar()")

	msg.SerializeAndDispatch(
		out,
		msg.TypeSetChar,
		msg.MessageSetChar{TermId: 0, X: x, Y: y, Char: char})
}

//sending messages back to hypervisor to set terminal
func SetCursor(out chan []byte, x uint32, y uint32) {
	println("(process/terminal/api.go).SetCursor()")

	msg.SerializeAndDispatch(
		out,
		msg.TypeSetCursor,
		msg.MessageSetCursor{TermId: 0, X: x, Y: y})
}
