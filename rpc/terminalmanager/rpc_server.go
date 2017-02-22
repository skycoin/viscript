package terminalmanager

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

const DEFAULT_PORT = "7777"

type RPC struct {
}

func NewRPC() *RPC {
	newRPC := &RPC{}
	return newRPC
}

func (self *RPC) Serve() {
	port := os.Getenv("TERMINAL_RPC_PORT")
	if port == "" {
		log.Println("No TERMINAL_RPC_PORT environmental variable is found, assignig default port value:", DEFAULT_PORT)
		port = DEFAULT_PORT
	}

	terminalManager := newTerminalManager()
	receiver := new(RPCReceiver)
	receiver.TerminalManager = terminalManager
	err := rpc.Register(receiver)
	if err != nil {
		panic(err)
	}

	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}

	log.Println("Serving RPC on port", port)
	err = http.Serve(l, nil)
	if err != nil {
		panic(err)
	}
}
