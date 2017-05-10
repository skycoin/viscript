package msg

const CATEGORY_Terminal uint16 = 0x0200 //flag

const (
	TypeVisualInfo      = 1 + CATEGORY_Terminal
	TypeCommandLine     = 2 + CATEGORY_Terminal
	TypeCommand         = 3 + CATEGORY_Terminal
	TypeTerminalIds     = 4 + CATEGORY_Terminal
	TypePutChar         = 5 + CATEGORY_Terminal
	TypeSetCharAt       = 6 + CATEGORY_Terminal
	TypeSetCursor       = 7 + CATEGORY_Terminal
	TypeFrameBufferSize = 8 + CATEGORY_Terminal //start of low level events
)

type MessageVisualInfo struct {
	NumColumns       uint32
	NumRows          uint32
	NumRowsForPrompt uint32
	CurrRow          uint32
}

type MessageCommandLine struct { //updates/replaces current command line on any change
	TermId       uint32
	CommandLine  string
	CursorOffset uint32 //from first character of command line
}

type MessageCommand struct {
	Command string
	Args    []string
}

type MessageTerminalIds struct {
	Focused TerminalId
	TermIds []TerminalId
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

//low level events
type MessageFrameBufferSize struct {
	X uint32
	Y uint32
}
