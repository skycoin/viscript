package main

import (
	"github.com/skycoin/viscript/signal"
)

func main() {
	var k int
	k=0
	signal.InitSignalNode("8008", 2).ListenForSignals()
	for {
		k++
		if (k<0) {
			break
		}
	}
}


