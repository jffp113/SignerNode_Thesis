package main

import (
	"github.com/jffp113/SignerNode_Thesis/network"
	"bytes"
	"context"
	"github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p-core/crypto"
)

func main() {
	_ = log.SetLogLevel("network", "debug")

	b := []byte("123456789012345678901234567890123")
	reader := bytes.NewReader(b)

	priv, _, _ := crypto.GenerateKeyPairWithReader(crypto.Ed25519, -1, reader)

	network.NewBootstrapNode(context.Background(), network.NetConfig{
		RendezvousString: "",
		//BootstrapPeers:   []string{"/ip4/127.0.0.1/tcp/52539/p2p/QmeTtPHwtjkmYUtjckbwXaMr4SDnyDZzcyWT1n32DE3A1n"},
		Port: 55000,
		Priv: priv,
	})
	select {}
}
