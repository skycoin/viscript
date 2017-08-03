package signal

import (
	"flag"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	serverAddress string
	runClient     bool
	clientId      uint
)

func init() {
	flagSet := flag.NewFlagSet("signal", flag.ContinueOnError)
	flagSet.StringVar(&serverAddress, "signal-server-address", "localhost:7999", "address of signal server")
	flagSet.UintVar(&clientId, "signal-client-id", 1, "id of signal client")
	flagSet.BoolVar(&runClient, "signal-client", false, "run signal client")
	index := 0
	for i, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-") {
			index = i
			break
		}
	}
	flagSet.Parse(os.Args[index+1:])

	if runClient {
		go func() {
			c, err := Connect(serverAddress, clientId)
			for {
				if err != nil {
					log.Errorf("connect to viscript failed %v", err)
				}
				c.WaitUntilDisconnected()
				// sleep 30s to reconnect
				time.Sleep(30 * time.Second)
				err = c.Connect(serverAddress)
			}
		}()
	}
}
