package msg

const ChannelCapacity = 4096 // FIXME?  might only need capacity of 2?
// .... onChar is always paired with an immediate onKey, making 2 entries at once

type TaskInterface interface {
	GetId() TaskId
	GetType() TaskType
	GetLabel() string
	GetIncomingChannel() chan []byte //channel for incoming messages
	Tick()                           //digest the messages and emit messages
}

type ExtTaskInterface interface {
	Tick()
	Start() error
	Attach() error
	Detach()
	TearDown()
	GetId() ExtTaskId
	GetFullCommandLine() string
	GetTaskInChannel() chan []byte
	GetTaskOutChannel() chan []byte
	GetTaskExitChannel() chan struct{}
}
