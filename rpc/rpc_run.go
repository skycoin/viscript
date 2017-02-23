package main

import "github.com/corpusc/viscript/rpc/terminalmanager"

func run() {
	rpcInstance := terminalmanager.NewRPC()
	rpcInstance.Serve()
}
