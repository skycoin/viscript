package msg

const CATEGORY_App uint16 = 0x0300 //flag

const (
	TypeUserCommand = 1 + CATEGORY_App
)

type MessageUserCommand struct {
	Sequence uint32
	AppId    uint32
	Payload  []byte
}
