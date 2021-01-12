package main

import (
	"SignerNode/network"
	"context"
	"github.com/ipfs/go-log"
)

func main() {
	_ = log.SetLogLevel("network", "debug")
	network.NewBootstrapNode(context.Background(),network.NetConfig{
		RendezvousString: "",
		//BootstrapPeers:   []string{"/ip4/127.0.0.1/tcp/52539/p2p/QmeTtPHwtjkmYUtjckbwXaMr4SDnyDZzcyWT1n32DE3A1n"},
		Port:             55349,
	})
	select {}
}
