package process

import (
	"fmt"
	"github.com/corpusc/viscript/msg"
)

func Printfc(out chan []byte, format string, vars ...interface{}) {
	println("(process/terminal/api.go).Printfc()")

	b := fmt.Sprintf(format, vars)

	for _, v := range b {
		PutChar(out, uint32(v))
	}
}

func PutChar(out chan []byte, char uint32) {
	println("(process/terminal/api.go).PutChar()")

	msg.SerializeAndDispatch(
		out,
		msg.TypePutChar,
		msg.MessagePutChar{0, char})
}

//sending messages back to hypervisor to set terminal
func SetCharAt(out chan []byte, x uint32, y uint32, char uint32) {
	println("(process/terminal/api.go).SetCharAt()")

	msg.SerializeAndDispatch(
		out,
		msg.TypeSetCharAt,
		msg.MessageSetCharAt{TermId: 0, X: x, Y: y, Char: char})
}

//sending messages back to hypervisor to set terminal
func SetCursor(out chan []byte, x uint32, y uint32) {
	println("(process/terminal/api.go).SetCursor()")

	msg.SerializeAndDispatch(
		out,
		msg.TypeSetCursor,
		msg.MessageSetCursor{TermId: 0, X: x, Y: y})
}
