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

	var wg sync.WaitGroup
	wg.Add(1)

	n , _ :=network.CreateNetwork(context.Background(),network.NetConfig{
		RendezvousString: "",
		BootstrapPeers:   []string{"/ip4/127.0.0.1/tcp/59889/p2p/Qmb8srjSb9JZqqqZK65fhu4U75gU5Lg9SoSP2PBnoWZyBB"},
		Port:             58869,
	})

	for {
		//fmt.Println("Waiting")
		//fmt.Println(string(n.Receive()))
		time.Sleep(10*time.Second)
		fmt.Println(n.Broadcast([]byte("Hello 1")))
	}


}
