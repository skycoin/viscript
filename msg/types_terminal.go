package msg

const CATEGORY_Terminal uint16 = 0x0200 //flag

const (
	TypePutString = 1 + CATEGORY_Terminal
	TypePutChar   = 2 + CATEGORY_Terminal
	TypeSetCharAt = 3 + CATEGORY_Terminal
	TypeSetCursor = 4 + CATEGORY_Terminal

	// low level events
	TypeFrameBufferSize = 4 + CATEGORY_Terminal
)

type MessagePutString struct {
	TermId uint32
	String string
}

type MessagePutChar struct {
	TermId uint32
	Char   uint32
}

type MessageSetCharAt struct {
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

// low level events
type MessageFrameBufferSize struct {
	X uint32
	Y uint32
}
