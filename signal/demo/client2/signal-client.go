package main

import (
	"github.com/skycoin/viscript/signal"
)

func main() {
	signal.InitSignalNode("8008").ListenForSignals()
}


