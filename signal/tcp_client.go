package signal

import (
	"net"
	"io"
	"log"
	"github.com/skycoin/viscript/msg"
	sgmsg "github.com/skycoin/viscript/signal/msg"
	"runtime"
	"time"
)

type SignalNode struct {
	port         string
	appId        uint32
}

func InitSignalNode(port string, appId uint32) *SignalNode {
	client := &SignalNode{port: port,
		appId: appId,
	}
	return client
}

func (self *SignalNode) ListenForSignals() {
	listenAddress := "0.0.0.0:" + self.port
	l, err := net.Listen("tcp", listenAddress)
	if err != nil {
		panic(err)
	}
	log.Println("Listen for incoming message on port: " + self.port)
	for {
		appConn, err := l.Accept() // create a connection with the user app (e.g. browser)
		if err != nil {
			log.Println("Cannot accept client's connection")
			return
		}
		defer appConn.Close()

		go func() { // run listening the connection for data and sending it through the meshnet to the server
			for {
				sizeMessage := make([]byte, 40)
				_, err := appConn.Read(sizeMessage)
				if err != nil {
					if err == io.EOF {
						continue
					} else {
						log.Println(err)
					}

				}


				switch msg.GetType(sizeMessage) {

				case sgmsg.TypeUserCommand:
					uc := sgmsg.MessageUserCommand{}
					err = msg.Deserialize(sizeMessage, &uc)
					if err != nil {
						log.Println("Incorrect UserCommand:", sizeMessage)
						continue
					}

					self.handleUserCommand(&uc)

				default:
					log.Println("Bad command")
				}
			}
		}()
	}
}

func (self *SignalNode) handleUserCommand(uc *sgmsg.MessageUserCommand) {
	log.Println("command received:", uc)
	sequence := uc.Sequence
	//appId := uc.AppId
	message := uc.Payload

	test := sgmsg.MessageUserCommand{}
	err := msg.Deserialize(uc.Payload, &test)
	if err != nil {
		log.Println("Incorrect UserCommand:", uc.Payload)
	}

	switch msg.GetType(test.Payload) {

	case sgmsg.TypePing:
		log.Println("ping command")
		ack := &sgmsg.MessagePingAck{}
		ackS := msg.Serialize(sgmsg.TypePingAck, ack)
		self.SendAck(ackS, sequence, self.appId)

	case sgmsg.TypeResourceUsage:
		log.Println("res_usage command")
		cpu, memory, err := GetResources()
		if err == nil {
			ack := &sgmsg.MessageResourceUsageAck{
				cpu,
				memory,
			}
			ackS := msg.Serialize(sgmsg.TypeResourceUsageAck, ack)
			self.SendAck(ackS, sequence, self.appId)
		}

	case sgmsg.TypeShutdown:
		log.Println("shutdown command")
		shutdown := sgmsg.MessageShutdown{}
		err = msg.Deserialize(test.Payload, &shutdown)
		if err != nil {
			panic(err)
		}

		switch shutdown.Stage {
			case 1:
				log.Println("app is preparing for shutdown... ", shutdown.Stage)
				ack := &sgmsg.MessageShutdownAck{Stage: 1}
				ackS := msg.Serialize(sgmsg.TypeShutdownAck, ack)
				self.SendAck(ackS, sequence, self.appId)
			case 2:
				log.Println("turn off daemons... ", shutdown.Stage)
				self.TurnOffNodes()
				ack := &sgmsg.MessageShutdownAck{Stage: 2}
				ackS := msg.Serialize(sgmsg.TypeShutdownAck, ack)
				self.SendAck(ackS, sequence, self.appId)
			case 3:
				ack := &sgmsg.MessageShutdownAck{Stage: 3}
				ackS := msg.Serialize(sgmsg.TypeShutdownAck, ack)
				self.SendAck(ackS, sequence, self.appId)
				panic("goodbye")
		}


	default:
		log.Println("Unknown user command:", message)

	}
}

func (self *SignalNode) TurnOffNodes(){
	time.Sleep(2* time.Second)
	log.Println("Daemons turned off")
}

func (self *SignalNode) SendAck(ackS []byte, sequence, appId uint32) {
	ucAck := &msg.MessageUserCommandAck{
		sequence,
		self.appId,
		ackS,
	}
	ucAckS := msg.Serialize(msg.TypeUserCommandAck, ucAck)
	self.send(ucAckS)
}

func (self *SignalNode) send(data []byte) {
	conn, e := net.Dial("tcp", "127.0.0.1:7999")
	if e != nil {
		log.Println("bad conn")
	}
	_, err := conn.Write(data)
	if err != nil {
		log.Println("Unsuccessful sending")
	}
}

func GetResources() (float64, uint64, error) {
	return 0, getMemStats(), nil
}

func getMemStats() uint64 {
	ms := &runtime.MemStats{}
	runtime.ReadMemStats(ms)
	return ms.Alloc
}

func getCPUProfile() {

}