package main

import (
	"github.com/skycoin/viscript/signal"
)

func main() {
	var k int
	k=0
	signal.InitSignalNode("8001").ListenForSignals()
	for {
		k++
		if (k<0) {
			break
		}
	}
}


