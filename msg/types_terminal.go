package msg

// _ "bytes"
// _ "encoding/binary"
// _ "fmt"

/*
Messages from:
	Process -> HyperVisor
*/

//Note:
// - terminal resource IDs, destroy determinism
// - how do we eliminate or abstract resource ids?

const TypePrefix_Terminal uint16 = 0x02 // terminal message prefix

// Process to Hypervisor, input event
const (
	TypePutChar   = 1 + TypePrefix_Terminal
	TypeSetChar   = 2 + TypePrefix_Terminal
	TypeSetCursor = 3 + TypePrefix_Terminal
)

type MessagePutChar struct {
	TermId uint32 //terminal
	Char   uint32
}

type MessageSetChar struct {
	TermId uint32 //terminal
	X      uint32
	Y      uint32
	Char   uint32
}

type MessageSetCursor struct {
	TermId uint32 //terminal
	X      uint32
	Y      uint32
}
