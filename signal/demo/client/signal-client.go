package main

import (
	"github.com/skycoin/viscript/signal"
)

func main() {
	var k int
	k=0
	signal.ListenForSignals("8001")
	for {
		k++
		if (k<0) {
			break
		}
	}
}


