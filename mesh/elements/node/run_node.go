package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/skycoin/skycoin/src/mesh/messages"
	"github.com/skycoin/skycoin/src/mesh/node"
)

func main() {
	args := os.Args
	if len(args) < 5 {
		panic("not sufficient number of args")
	}
	host, nmAddr, connect, seqStr := args[1], args[2], args[3], args[4]

	seqInt, err := strconv.Atoi(seqStr)
	if err != nil {
		panic(err)
	}
	if seqInt < 0 {
		panic("negative sequence")
	}
	sequence := uint32(seqInt)

	fmt.Println("host:", host)
	fmt.Println("nmAddr:", nmAddr)
	fmt.Println("connect:", connect)

	need_connect := connect == "true"

	var n messages.NodeInterface

	if need_connect {
		n, err = node.CreateAndConnectNode(host, nmAddr)
	} else {
		n, err = node.CreateNode(host, nmAddr)
	}
	if err != nil {
		panic(err)
	}

	n.TalkToViscript(sequence)
}
