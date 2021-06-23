package network

import (
	"bytes"
	"context"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBroadcastToTwoPeers(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	createBootstrapPeer(ctx)

	config := NetConfig{
		BootstrapPeers: []string{"/ip4/127.0.0.1/tcp/55000/p2p/12D3KooWD1yUy23iVGYCYMZdm2fUy65WFaAc2H2i7ycBT3oJdN1B"},
		PeerAddress: "/ip4/127.0.0.1/tcp/",
		Port: 0,
	}
	net1, err := CreateNetwork(ctx, config)
	assert.Nil(t, err)
	net2, err := CreateNetwork(ctx, config)
	assert.Nil(t, err)
	net3, err := CreateNetwork(ctx, config)
	assert.Nil(t, err)

	time.Sleep(2 * time.Second)

	content := []byte{123}
	err = net1.Broadcast(content)
	assert.Nil(t, err)

	failWithTimeOut(t, 5*time.Second, func() {
		r2 := net2.Receive()
		r1 := net3.Receive()
		assert.Equal(t, r2.GetData(), content)
		assert.Equal(t, r1.GetData(), content)
	})
}

func TestBroadcastToGroup(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	createBootstrapPeer(ctx)

	config := NetConfig{
		BootstrapPeers: []string{"/ip4/127.0.0.1/tcp/55000/p2p/12D3KooWD1yUy23iVGYCYMZdm2fUy65WFaAc2H2i7ycBT3oJdN1B"},
		PeerAddress: "/ip4/127.0.0.1/tcp/",
		Port: 0,
	}
	net1, err := CreateNetwork(ctx, config)
	assert.Nil(t, err)
	net2, err := CreateNetwork(ctx, config)
	assert.Nil(t, err)
	net3, err := CreateNetwork(ctx, config)
	assert.Nil(t, err)

	groupName := "Test"
	err = net1.JoinGroup(groupName)
	assert.Nil(t, err)
	err = net2.JoinGroup(groupName)
	assert.Nil(t, err)
	err = net3.JoinGroup(groupName)
	assert.Nil(t, err)

	time.Sleep(2 * time.Second)

	content := []byte{123}
	err = net1.BroadcastToGroup(groupName, content)
	assert.Nil(t, err)

	failWithTimeOut(t, 5*time.Second, func() {
		r2 := net2.Receive()
		r3 := net3.Receive()
		assert.Equal(t, r2.GetData(), content)
		assert.Equal(t, r3.GetData(), content)
	})

	err = net2.LeaveGroup(groupName)
	assert.Nil(t, err)
	err = net1.BroadcastToGroup(groupName, content)
	assert.Nil(t, err)

	failWithTimeOut(t, 5*time.Second, func() {
		r3 := net3.Receive()
		assert.Equal(t, r3.GetData(), content)
	})

	passIfTimeOut(t, 2*time.Second, func() {
		_ = net2.Receive()
	})

}

func createBootstrapPeer(ctx context.Context) {
	b := []byte("123456789012345678901234567890123")
	reader := bytes.NewReader(b)
	priv, _, _ := crypto.GenerateKeyPairWithReader(crypto.Ed25519, -1, reader)

	NewBootstrapNode(ctx, NetConfig{
		RendezvousString: "",
		Port:             55000,
		PeerAddress: 	"/ip4/127.0.0.1/tcp/",
		Priv:             priv,
	})

}

func failWithTimeOut(t *testing.T, to time.Duration, test func()) {
	timeout := time.After(to)
	done := make(chan bool)
	go func() {
		test()
		done <- true
	}()

	select {
	case <-timeout:
		t.Fatal("Test didn't finish in time")
	case <-done:
	}
}

func passIfTimeOut(t *testing.T, to time.Duration, test func()) {
	ctx, _ := context.WithDeadline(
		context.Background(),
		time.Now().Add(to),
	)
	go func() {
		test()
		select {
		case <-ctx.Done():
		default:
			t.Fatal("Timeout didn't finish")
		}

	}()

	select {
	case <-ctx.Done():
	}
}
