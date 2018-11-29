package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"sync"
	"syscall"
	"time"

	"github.com/centrifugal/centrifuge-go"
	"github.com/realitycheck/centtest/lib"
)

var (
	server      = "ws://localhost:8000/connection/websocket"
	idsource    = "uuid"
	numUsers    = 1
	numChannels = 1
	numClients  = 1
	quiet       bool
	debug       bool
	chUser      = "user"
	chSystem    = "system:admin"
	timeout     = 5 * time.Second
	cpuprofile  = ""
	memprofile  = ""
)

func init() {
	flag.StringVar(&server, "s", server, "Server to connect")
	flag.IntVar(&numUsers, "nu", numUsers, "Number of users")
	flag.IntVar(&numClients, "nc", numClients, "Number of clients per user")
	flag.IntVar(&numChannels, "nch", numChannels, "Number of channels per user")
	flag.StringVar(&idsource, "i", idsource, "IDs Source")
	flag.BoolVar(&quiet, "quiet", quiet, "Quiet mode, enable to log nothing (default false)")
	flag.BoolVar(&debug, "debug", debug, "Debug mode, enable for verbose logging (default false)")
	flag.StringVar(&chUser, "ch-user", chUser, "Name of user channel")
	flag.StringVar(&chSystem, "ch-system", chSystem, "Name of system channel")
	flag.DurationVar(&timeout, "timeout", timeout, "Timeouts")
	flag.StringVar(&cpuprofile, "cpuprofile", cpuprofile, "write cpu profile to `file`")
	flag.StringVar(&memprofile, "memprofile", memprofile, "write memory profile to `file`")

}

func main() {
	flag.Parse()

	if quiet {
		log.SetOutput(ioutil.Discard)
		log.SetFlags(0)
	}

	clientConfig := centrifuge.DefaultConfig()
	clientConfig.HandshakeTimeout = timeout
	clientConfig.ReadTimeout = timeout
	clientConfig.WriteTimeout = timeout

	wg := sync.WaitGroup{}

	tt := make(chan *centtest.Test)
	go func() {
		for t := range tt {
			go t.Run()
			t := t
			defer func() {
				t.Close()
				wg.Done()
			}()
		}
	}()

	log.Printf("main: initing...")
	g := centtest.NewIDGenerator(idsource)

	for i := 0; i < numUsers; i++ {
		u := centtest.NewUser(g)
		for j := 0; j < numClients; j++ {
			c := centtest.NewClient(server, clientConfig)
			for k := 0; k < numChannels; k++ {
				wg.Add(1)
				name := fmt.Sprintf("%s:%d", chUser, k)
				ch := centtest.NewChannel(name).Attach(u)
				tt <- centtest.NewTest(u, c, ch, debug)
			}
			wg.Add(1)
			ch := centtest.NewChannel(chSystem)
			tt <- centtest.NewTest(u, c, ch, debug)
		}
	}

	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	if memprofile != "" {
		f, err := os.Create(memprofile)
		if err != nil {
			log.Fatal("could not create MEM profile: ", err)
		}
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write MEM profile: ", err)
		}
		f.Close()
	}

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("main: waiting...")
	select {
	case <-exit:
		close(tt)
	}

	log.Printf("main: exiting...")
	wg.Wait()
}
