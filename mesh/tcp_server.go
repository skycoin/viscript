package main

// only for example; we don't need it here, something similar should be in corpusc/viscript

import (
	"io"
	"log"
	"net"

	"github.com/skycoin/skycoin/src/mesh/messages"
)

type Server struct {
	address string
}

const ADDR = "0.0.0.0:7999"

func main() {
	server := Server{address: ADDR}
	server.listen()
}

func (self *Server) listen() {
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
				message := make([]byte, 16384)

				n, err := appConn.Read(message)
				if err != nil {
					return
					if err == io.EOF {
						continue
					} else {
						log.Printf("error while reading message from %s: %s\n", remoteAddr, err)
						break
					}
				}

				self.handleUserCommand(message[:n])
			}
		}()
	}
}

func (self *Server) handleUserCommand(ucS []byte) {
	uc := &messages.UserCommand{}
	err := messages.Deserialize(ucS, uc)
	if err != nil {
		panic(err)
	}
	log.Println("received message for sequence", uc.Sequence)
	//do something e.g. ack handling etc.
}
