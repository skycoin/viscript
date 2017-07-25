package signal

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
	sgmsg "github.com/skycoin/viscript/signal/msg"
)

var Sequence uint32 = 0

var CountApps uint32 = 1

func incrementCountApps() uint32 {
	CountApps++
	return CountApps
}

func GetNextMessageID() uint32 {
	Sequence++
	return Sequence
}

type MonitorServer struct {
	address          string
	lock             *sync.Mutex
	connections      map[uint32]net.Conn
	responseChannels map[uint32]chan []byte
	sequence         uint32
}

var Monitor *MonitorServer

func Init(address string) *MonitorServer {
	Monitor = NewMonitorServer(address)
	return Monitor
}

func NewMonitorServer(address string) *MonitorServer {
	server := &MonitorServer{address: address}
	server.lock = &sync.Mutex{}
	server.responseChannels = make(map[uint32]chan []byte)
	server.connections = make(map[uint32]net.Conn)
	server.sequence = Sequence
	return server
}

func (self *MonitorServer) Run() {
	go func() {
		self.Serve()
	}()
}

func (self *MonitorServer) Serve() {
	address := self.address

	l, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}

	log.Println("Listening for incoming messages on", self.address)

	for {
		appConn, err := l.Accept() // accept a connection which is created by an app
		if err != nil {
			log.Println("Cannot accept client's connection:", err)
			return
		}
		defer appConn.Close()

		remoteAddr := appConn.RemoteAddr().String()
		go func() { // run listening the connection for user command exchange between viscript and app (ping, shutdown, resources request etc.)
			for {
				message := make([]byte, 42)

				_, err := appConn.Read(message)
				if err != nil {
					return
					if err == io.EOF {
						continue
					} else {
						log.Printf("error while reading message from %s: %s\n", remoteAddr, err)
						break
					}
				}
				uc := &sgmsg.MessageUserCommandAck{}
				err = sgmsg.Deserialize(message, uc)
				if err != nil {
					panic(err)
				}
				//log.Println("received message for sequence", uc.Sequence)

				appId := uc.AppId
				//sequence := uc.Sequence

				self.lock.Lock()
				if _, ok := self.connections[appId]; !ok { // if viscript already created an app, this connection is already in the map
					self.connections[appId] = appConn // if no, then add current connection to the map; so next iterations will skip this step
				}
				respChan, ok0 := self.responseChannels[uc.Sequence] // take response channel for responding to it
				self.lock.Unlock()
				if !ok0 {
					log.Println("no response channel", err)
					continue
				}
				respChan <- uc.Payload // respond to it
			}
		}()
	}
}

//func (self *MonitorServer) ReadFrom(appId msg.ExtTaskId) ([]byte, error) {
//	appMessageChannel, exists := self.responseChannels[uint32(appId)]
//	if !exists {
//		errString := fmt.Sprintf("Channel with ID: %d doesn't exist.", appId)
//		err := errors.New(errString)
//		return []byte{}, err
//	}
//
//	select {
//	case data := <-appMessageChannel:
//		return data, nil
//	default:
//	}
//
//	return []byte{}, errors.New(string(appId) + " app channel is empty.")
//}

func (self *MonitorServer) PrintAll() {
	for key, _ := range self.responseChannels {
		println(key)
	}
}

func (self *MonitorServer) Send(appId uint32, message []byte) ([]byte, string, error) {
	respChan, sequence := self.MakeResponseChannel()

	self.lock.Lock()
	conn, e := self.connections[appId]
	if !e {
		log.Println("bad conn")
	}
	self.lock.Unlock()
	var n uint32 = appId
	str := fmt.Sprint(n)

	uc := &sgmsg.MessageUserCommand{sequence, appId, message}
	ucS := sgmsg.Serialize(sgmsg.TypeUserCommand, uc)
	sendTime := time.Now()
	_, err := conn.Write(ucS)
	if err != nil {
		return nil, str, err
	}
	response, err := self.Wait(respChan, sequence)

	switch sgmsg.GetType(response) {

	case sgmsg.TypeResourceUsageAck:
		answer := sgmsg.MessageResourceUsageAck{}
		err = sgmsg.Deserialize(response, &answer)
		if err != nil {
			panic(err)
		}
		log.Println("cpu: ", answer.CPU, "memory: ", answer.Memory)

	case sgmsg.TypePingAck:
		getTime := time.Now()
		log.Print(getTime.Sub(sendTime).Seconds()*1000, " ms")

	case sgmsg.TypeShutdownAck:
		answer := sgmsg.MessageShutdownAck{}
		err = sgmsg.Deserialize(response, &answer)
		if err != nil {
			panic(err)
		}

		switch answer.Stage {
		case 1:
			log.Println("shutdown stage ", answer.Stage, " is over")
			self.SendShutdownCommand(appId, 2)
		case 2:
			log.Println("shutdown stage ", answer.Stage, " is over")
			self.SendShutdownCommand(appId, 3)
		case 3:
			log.Println("shutdown stage ", answer.Stage, " is over")
			log.Println("app is closed.")
		}

	case sgmsg.TypeStartupAck:
		answer := sgmsg.MessageStartupAck{}
		err = sgmsg.Deserialize(response, &answer)
		if err != nil {
			panic(err)
		}

		switch answer.Stage {
		case 1:
			log.Println("startup stage ", answer.Stage, " is over")
			self.SendStartupCommand(appId, 2)
		case 2:
			log.Println("startup stage ", answer.Stage, " is over")
			self.SendStartupCommand(appId, 3)
		case 3:
			log.Println("startup stage ", answer.Stage, " is over")
			log.Println("app is up.")
		}

	default:
		log.Println("Incorrect command type")
	}

	return response, "end", err
}

func (self *MonitorServer) AddSignalNodeConn(address string, port string) {
	str := address + ":" + port
	conn, e := net.Dial("tcp", str)
	if e != nil {
		log.Println("Can't add this node.")
	}
	self.connections[CountApps] = conn
	self.SendStartupCommand(CountApps, 1)
	incrementCountApps()
}

func (self *MonitorServer) SendPingCommand(appId uint32) {
	msgUserCommand := sgmsg.MessageUserCommand{
		Sequence: 1,
		AppId:    appId,
		Payload:  sgmsg.Serialize(sgmsg.TypePing, sgmsg.MessagePing{})}

	serializedCommand := sgmsg.Serialize(sgmsg.TypeUserCommand, msgUserCommand)

	self.Send(appId, serializedCommand)
}


func (self *MonitorServer) SendShutdownCommand(appId uint32, stage uint32) {
	msgUserCommand := sgmsg.MessageUserCommand{
		Sequence: 1,
		AppId:    appId,
		Payload:  sgmsg.Serialize(sgmsg.TypeShutdown, sgmsg.MessageShutdown{  Stage: stage})}

	serializedCommand := sgmsg.Serialize(sgmsg.TypeUserCommand, msgUserCommand)

	self.Send(appId, serializedCommand)
}

func (self *MonitorServer) SendStartupCommand(appId uint32, stage uint32) {
	msgUserCommand := sgmsg.MessageUserCommand{
		Sequence: 1,
		AppId:    appId,
		Payload:  sgmsg.Serialize(sgmsg.TypeStartup, sgmsg.MessageStartup{  Address: self.address, Stage: stage})}

	serializedCommand := sgmsg.Serialize(sgmsg.TypeUserCommand, msgUserCommand)

	self.Send(appId, serializedCommand)
}



func (self *MonitorServer) SendResUsageCommand(appId uint32) {
	msgUserCommand := sgmsg.MessageUserCommand{
		Sequence: 1,
		AppId:    appId,
		Payload:  sgmsg.Serialize(sgmsg.TypeResourceUsage, sgmsg.MessageResourceUsage{})}

	serializedCommand := sgmsg.Serialize(sgmsg.TypeUserCommand, msgUserCommand)

	self.Send(appId, serializedCommand)
}


func (self *MonitorServer) ListNodes() {
	for i:=1; i<=len(self.connections); i++ {
		log.Println("appId: ", i+1, "remote addres: ", self.connections[uint32(i)].RemoteAddr())
	}
}

func (self *MonitorServer) MakeResponseChannel() (chan []byte, uint32) {
	respChan := make(chan []byte)

	self.lock.Lock()
	sequence := self.sequence
	self.responseChannels[sequence] = respChan
	self.sequence++
	self.lock.Unlock()

	return respChan, sequence
}

func (self *MonitorServer) Wait(respChan chan []byte, sequence uint32) ([]byte, error) {
	select {
	case response := <-respChan:
		return response, nil
	case <-time.After(time.Second * 10):
		return nil, errors.New(fmt.Sprintf("Timeout when receiving response for %d\n", sequence))
	}
}
