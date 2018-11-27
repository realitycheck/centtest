package main

import (
	"flag"
	"sync"

	"github.com/centrifugal/centrifuge-go"
	"github.com/realitycheck/centtest/lib"
)

var (
	server      = "ws://localhost:8000/connection/websocket"
	idsource    = "uuid"
	numUsers    = 1
	numChannels = 1
	numClients  = 1
)

func init() {
	flag.StringVar(&server, "s", server, "Server to connect")
	flag.IntVar(&numUsers, "nu", numUsers, "Number of users")
	flag.IntVar(&numClients, "nc", numClients, "Number of clients per user")
	flag.IntVar(&numChannels, "nch", numChannels, "Number of channels per user")
	flag.StringVar(&idsource, "i", idsource, "IDs Source")
}

func main() {
	flag.Parse()

	clientConfig := centrifuge.DefaultConfig()
	g := centtest.NewIDGenerator(idsource)

	wg := &sync.WaitGroup{}

	for i := 0; i < numUsers; i++ {
		u := centtest.NewUser(g)
		for j := 0; j < numClients; j++ {
			c := centtest.NewClient(server, clientConfig)
			for k := 0; k < numChannels; k++ {
				wg.Add(1)
				go centtest.Run(u, c, centtest.NewChannel(k), wg)
			}
		}
	}

	wg.Wait()
}
