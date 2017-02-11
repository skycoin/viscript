package msg

const CATEGORY_Terminal uint16 = 0x0200 //flag

const (
	TypePutChar   = 1 + CATEGORY_Terminal
	TypeSetChar   = 2 + CATEGORY_Terminal
	TypeSetCursor = 3 + CATEGORY_Terminal
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

// Terminal Driving Messages
type MessageFrameBufferSize struct {
	X uint32
	Y uint32
}
