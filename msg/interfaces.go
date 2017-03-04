package msg

const ChannelCapacity = 64 // FIXME?  might only need capacity of 2?
// .... onChar is always paired with an immediate onKey, making 2 entries at once

type ProcessInterface interface {
	GetId() ProcessId
	GetType() ProcessType
	GetLabel() string
	GetIncomingChannel() chan []byte //channel for incoming messages
	Tick()                           //process the messages and emit messages
}
