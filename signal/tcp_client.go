package signal

import (
	"net"
	"io"
	"log"
	"github.com/skycoin/viscript/msg"
	"runtime"
)

func ListenForSignals(port string) {
	listenAddress := "0.0.0.0:" + port
	l, err := net.Listen("tcp", listenAddress)
	if err != nil {
		panic(err)
	}
	log.Println("Listen for incoming message on port: " + port)
	for {
		appConn, err := l.Accept() // create a connection with the user app (e.g. browser)
		if err != nil {
			log.Println("Cannot accept client's connection")
			return
		}
		defer appConn.Close()

		go func() { // run listening the connection for data and sending it through the meshnet to the server
			for {
				sizeMessage := make([]byte, 30)
				_, err := appConn.Read(sizeMessage)
				if err != nil {
					if err == io.EOF {
						continue
					} else {
						log.Println(err)
					}

				}


				switch msg.GetType(sizeMessage) {

				case msg.TypeUserCommand:
					uc := msg.MessageUserCommand{}
					err = msg.Deserialize(sizeMessage, &uc)
					if err != nil {
						log.Println("Incorrect UserCommand:", sizeMessage)
						continue
					}

					handleUserCommand(&uc)

				default:
					log.Println("Bad command")
				}
			}
		}()
	}
}

func handleUserCommand(uc *msg.MessageUserCommand) {
	log.Println("command received:", uc)
	//sequence := uc.Sequence
	//appId := uc.AppId
	message := uc.Payload

	test := msg.MessageUserCommand{}
	err := msg.Deserialize(uc.Payload, &test)
	if err != nil {
		log.Println("Incorrect UserCommand:", uc.Payload)
	}

	switch msg.GetType(test.Payload) {

		case msg.TypePing:
			log.Println("ping command")
			ack := &msg.MessagePingAck{}
			ackS := msg.Serialize(msg.TypePingAck, ack)
			SendAck(ackS, 2, 1)

		case msg.TypeResourceUsage:
			log.Println("res_usage command")
			cpu, memory, err := GetResources()
			if err == nil {
				ack := &msg.MessageResourceUsageAck{
					cpu,
					memory,
				}
				ackS := msg.Serialize(msg.TypeResourceUsageAck, ack)
				SendAck(ackS, 2, 1)
			}

		case msg.TypeShutdown:
			log.Println("shutdown command")
			ack := &msg.MessageShutdownAck{}
			ackS := msg.Serialize(msg.TypeShutdownAck, ack)
			SendAck(ackS, 2, 1)
			panic("goodbye")

		default:
			log.Println("Unknown user command:", message)

		}
}

func SendAck(ackS []byte, sequence, appId uint32) {
	ucAck := &msg.MessageUserCommandAck{
		sequence,
		appId,
		ackS,
	}
	ucAckS := msg.Serialize(msg.TypeUserCommandAck, ucAck)
	send(ucAckS)
}

func send(data []byte) {
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