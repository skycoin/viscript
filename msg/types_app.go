package msg

const CATEGORY_App uint16 = 0x0300 //flag

const (
	MsgUserCommand = 1 + CATEGORY_App
)

type UserCommand struct {
	Sequence uint32
	AppId    uint32
	Payload  []byte
}
