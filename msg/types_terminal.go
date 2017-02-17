package msg

const CATEGORY_Terminal uint16 = 0x0200 //flag

const (
	TypePutChar   = 1 + CATEGORY_Terminal
	TypeSetChar   = 2 + CATEGORY_Terminal
	TypeSetCursor = 3 + CATEGORY_Terminal
	TypePutKey    = 4 + CATEGORY_Terminal // FIXME:

	// low level events
	TypeFrameBufferSize = 5 + CATEGORY_Terminal
)

type MessagePutChar struct {
	TermId uint32
	Char   uint32
}

type MessageSetChar struct {
	TermId uint32
	X      uint32
	Y      uint32
	Char   uint32
}

type MessageSetCursor struct {
	TermId uint32
	X      uint32
	Y      uint32
}

// FIXME: I didn't know what to name this
type MessagePutKey struct {
	TermId uint32
	Key    uint32
	Scan   uint32
	Action uint8
	Mod    uint8
}

// low level events
type MessageFrameBufferSize struct {
	X uint32
	Y uint32
}
