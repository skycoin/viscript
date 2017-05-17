package main

import (
	"os"
	"strconv"

	network "github.com/skycoin/skycoin/src/mesh/nodemanager"
)

func main() {
	args := os.Args
	if len(args) < 4 {
		panic("not sufficient number of args")
	}

	domainName := args[1]

	ctrlAddr := args[2]

	seqStr := args[3]
	seqInt, err := strconv.Atoi(seqStr)
	if err != nil {
		panic(err)
	}
	if seqInt < 0 {
		panic("negative sequence")
	}
	sequence := uint32(seqInt)

	appIdStr := args[4]
	appIdInt, err := strconv.Atoi(appIdStr)
	if err != nil {
		panic(err)
	}
	if appIdInt < 0 {
		panic("negative sequence")
	}
	appId := uint32(appIdInt)

	nm := network.NewNetwork(domainName, ctrlAddr)
	nm.TalkToViscript(sequence, appId)
}
