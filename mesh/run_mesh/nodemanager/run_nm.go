package main

import (
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

	appIdStr := args[2]
	appIdInt, err := strconv.Atoi(appIdStr)
	if err != nil {
		panic(err)
	}
	if appIdInt < 0 {
		panic("negative sequence")
	}
	appId := uint32(appIdInt)

	nm := network.NewNetwork()
	nm.TalkToViscript(sequence, appId)
}
