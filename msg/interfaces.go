package msg

//Process should implement
type ProcessInterface interface {
	GetId() ProcessId
	GetIncomingChannel() chan []byte //channel for incoming messages
	GetOutgoingChannel() chan []byte //channel for outgoing messages
	Tick()                           //process the messages and emit messages
}
