package msg

const ChannelCapacity = 4096 // FIXME?  might only need capacity of 2?
// .... onChar is always paired with an immediate onKey, making 2 entries at once

type TaskInterface interface {
	GetId() TaskId
	GetIncomingChannel() chan []byte
	GetText() string
	GetType() TaskType
	Tick()
}

type ExternalAppInterface interface {
	//shared vars (with task ^^^)
	GetId() ExternalAppId
	GetTaskInChannel() chan []byte
	Tick()
	//unique vars
	Attach() error
	Detach()
	GetFullCommandLine() string
	GetTaskOutChannel() chan []byte
	GetTaskExitChannel() chan struct{}
	Start() error
	TearDown()
}
