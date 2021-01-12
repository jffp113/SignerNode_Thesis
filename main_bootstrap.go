package main

import (
	"SignerNode/network"
	"context"
	"fmt"
	"github.com/ipfs/go-log"
	"sync"
	"time"
)

func main() {
	//log.SetAllLoggers(log.LogLevel(logging.DEBUG))
	_ = log.SetLogLevel("network", "debug")
	_ = log.SetLogLevel("pubsub", "debug")
	_ = log.SetLogLevel("connmgr", "debug")
	_ = log.SetLogLevel("dht", "warn")

	var wg sync.WaitGroup
	wg.Add(1)

	n , _ :=network.CreateNetwork(context.Background(),network.NetConfig{
		RendezvousString: "",
		BootstrapPeers:   []string{"/ip4/127.0.0.1/tcp/52539/p2p/QmeTtPHwtjkmYUtjckbwXaMr4SDnyDZzcyWT1n32DE3A1n"},
		Port:             55349,
	})

	go func() {
		for {
			fmt.Println("Waiting")
			fmt.Println(string(n.Receive()))
		}
	}()
	for {
		time.Sleep(10*time.Second)
		fmt.Println(n.Broadcast([]byte("1")))
	}

}
