package msg

import (
	_ "bytes"
	_ "encoding/binary"
	_ "fmt"
)

/*
Messages from:
	Process -> HyperVisor
*/

//Note:
// - terminal resource IDs, destroy determinism
// - how do we eliminate or abstract resource ids?

const PrefixTerminal uint16 = 0x02 // terminal message prefix

// Process to Hypervisor, input event
const (
	TypeSetChar   = 1 + PrefixTerminal
	TypeSetCursor = 2 + PrefixTerminal
)

type MessageSetChar struct {
	TermId uint32 //id of the terminal
	X      uint32
	Y      uint32
	Rune   uint32
}

type MessageSetCursor struct {
	TermId uint32 //id of the terminal
	X      uint32
	Y      uint32
}
