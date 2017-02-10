package msg

const ChannelCapacity = 64 // FIXME?  might only need capacity of 2?
// .... onChar is always paired with an immediate onKey, making 2 entries at once

//Process should implement
type ProcessInterface interface {
	GetId() ProcessId
	GetIncomingChannel() chan []byte //channel for incoming messages
	// GetOutChannel() chan []byte //channel for outgoing messages
	Tick() //process the messages and emit messages
}
