package network

import (
	"bytes"
	"context"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

//TestNewDiscovery tests if peers connected to a bootstrap peer
//can find each other
func TestNewDiscovery(t *testing.T) {
	b := []byte("123456789012345678901234567890123")
	reader := bytes.NewReader(b)
	priv, _, _ := crypto.GenerateKeyPairWithReader(crypto.Ed25519, -1, reader)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	NewBootstrapNode(ctx, NetConfig{
		RendezvousString: "",
		PeerAddress: "/ip4/127.0.0.1/tcp/",
		Port:             55000,
		Priv:             priv,
	})

	config := NetConfig{
		BootstrapPeers:   []string{"/ip4/127.0.0.1/tcp/55000/p2p/12D3KooWD1yUy23iVGYCYMZdm2fUy65WFaAc2H2i7ycBT3oJdN1B"},
		PeerAddress: "/ip4/127.0.0.1/tcp/",
		Port: 0,
	}

	//Lets first create a new peer discovery
	disc1 := createPeerDiscovery(t, ctx, config)
	time.Sleep(1 * time.Second)               // the time chosen to wait was random
	assert.Equal(t, 2, len(disc1.GetPeers())) //him self and the bootstrap peer

	disc2 := createPeerDiscovery(t, ctx, config)
	time.Sleep(1 * time.Second)
	assert.Equal(t, 3, len(disc2.GetPeers())) //bootstrap + disc1 + disc2
	assert.Equal(t, 3, len(disc1.GetPeers()))

	disc3 := createPeerDiscovery(t, ctx, config)
	time.Sleep(2 * time.Second)
	assert.Equal(t, 4, len(disc2.GetPeers())) //bootstrap + disc1 + disc2 + disc3
	assert.Equal(t, 4, len(disc1.GetPeers()))
	assert.Equal(t, 4, len(disc3.GetPeers()))

	//There is only a bootstrap peer
	assert.Equal(t, 1, len(disc1.GetBootstrapPeers()))
	assert.Equal(t, 1, len(disc2.GetBootstrapPeers()))
	assert.Equal(t, 1, len(disc3.GetBootstrapPeers()))

}

func createPeerDiscovery(t *testing.T, ctx context.Context, config NetConfig) *peerDiscovery {
	h, err := newPeerHost(config)
	assert.Nil(t, err)
	disc := NewDiscovery(ctx, h, config)
	_, err = disc.SetupDiscovery()
	assert.Nil(t, err)

	return disc
}
