package main

import (
	"SignerNode/network"
	"context"
	"github.com/ipfs/go-log"
	"sync"
)

func main() {
	//log.SetAllLoggers(log.LogLevel(logging.DEBUG))
	log.SetLogLevel("rendezvous", "debug")

	var wg sync.WaitGroup
	wg.Add(1)

	network.CreateNetwork(context.Background(),network.NetConfig{
		RendezvousString: "",
		BootstrapPeers:   []string{"/ip4/127.0.0.1/tcp/50900/p2p/QmTqcKo8N5XeCnAwqtDELn8Xs6LCCXN7LE5H11uQ9UrynS"},
		Port:             50870,
	})


	wg.Wait()
}
