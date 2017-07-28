package main

import (
	"github.com/skycoin/viscript/signal"
)

func main() {
	signal.InitSignalNode("8001").ListenForSignals()

}


