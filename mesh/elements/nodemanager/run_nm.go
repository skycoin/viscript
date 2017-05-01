package main

import (
	"log"
	"os"
	"strconv"

	network "github.com/skycoin/skycoin/src/mesh/nodemanager"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		panic("not sufficient number of args")
	}
	seqStr := args[1]
	seqInt, err := strconv.Atoi(seqStr)
	if err != nil {
		panic(err)
	}
	if seqInt < 0 {
		panic("negative sequence")
	}
	sequence := uint32(seqInt)
	nm := network.NewNetwork()
	log.Println("meshnet is running")
	nm.TalkToViscript(sequence)
}
