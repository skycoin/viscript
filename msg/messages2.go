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

const MP2 uint16 = 0x02 //input message prefix

// Process to Hypervisor, input event
const (
	TypeSetTerminal = 1 + MP2
	TypeSetCursor   = 2 + MP2
)

type MessageSetTerminal struct {
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
