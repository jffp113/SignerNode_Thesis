package network

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p"
	connmgr "github.com/libp2p/go-libp2p-connmgr"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"time"
)

var logger = log.Logger("network")

const NetworkBufSize = 128

type NetConfig struct {
	RendezvousString string
	BootstrapPeers   []string
	Port int
	//ListenAddresses  addrList
	//ProtocolID       string
}

type networkMessage struct {
	To      peer.ID
	From    peer.ID
	Content []byte
}

type network struct {
	messages chan []byte

	ctx   context.Context
	ps    *pubsub.PubSub
	topic *pubsub.Topic
	sub   *pubsub.Subscription
	self     peer.ID
}

type Network interface {
	Broadcast(msg []byte) error
	Send(node string, msg []byte)
	Receive() []byte
}

func (n *network) Broadcast(msg []byte) error {

	m := networkMessage{
		To:      "",
		Content: msg,
	}
	msgBytes, err := json.Marshal(m)
	if err != nil {
		return err
	}

	return n.topic.Publish(n.ctx, msgBytes)
}

func (n *network) Send(node string, msg []byte) {

}

func (n *network) Receive() []byte {
	//TODO add context
	return <-n.messages
}


func CreateNetwork(ctx context.Context,config NetConfig) (Network,error) {
	logger.Debug("Setting up Network")
	h,err := newPeerHost(config)
	discovery := NewDiscovery(ctx,h,config)

	if err != nil {
		logger.Error(err)
		return nil,err
	}

	logger.Infof("Peer will be available at %v/p2p/%s",h.Addrs()[0],h.ID().Pretty())

	disc,err := discovery.SetupDiscovery()

	if err != nil {
		logger.Error(err)
		return nil,err
	}

	ps,err := pubsub.NewGossipSub(ctx,h,pubsub.WithDiscovery(disc))

	topic, err := ps.Join("SignerNodeNetwork")

	if err != nil {
		logger.Error(err)
		return nil,err
	}

	sub,err := topic.Subscribe()

	if err != nil {
		logger.Error(err)
		return nil,err
	}

	network := &network{
		messages: make(chan []byte,NetworkBufSize),
		ctx:      ctx,
		ps:       ps,
		topic:    topic,
		sub:      sub,
		self: h.ID(),
	}

	go processIncomingMsg(network)

	showConnectedListPeers(network)

	return network,nil
}

func processIncomingMsg(n *network){
	for{
		logger.Debug("Waiting for new message")
		msg,err := n.sub.Next(n.ctx)
		logger.Debugf("New message arrived from",msg.ReceivedFrom)

		if err != nil {
			close(n.messages)
			return
		}
		// only forward messages delivered by others

		if msg.ReceivedFrom == n.self {
			continue
		}

		cm := new(networkMessage)
		err = json.Unmarshal(msg.Data, &cm)
		if err != nil {
			continue
		}

		// send valid messages onto the Messages channel
		n.messages <- cm.Content
	}
}

func newPeerHost(config NetConfig) (host.Host, error) {

	logger.Debug("Creating Peer Host")

	//listenAddr := fmt.Sprintf("/ip4/0.0.0.0/tcp/%v",config.Port)
	listenAddr := fmt.Sprintf("/ip4/0.0.0.0/tcp/%v",0)

	return libp2p.New(
		context.Background(),
		libp2p.ListenAddrStrings(listenAddr),
		libp2p.ConnectionManager(connmgr.NewConnManager(
			1,          // Lowwater
			3,          // HighWater,
			time.Minute, // GracePeriod
		)),
		//libp2p.Identity(*prvKey),
		//libp2p.Security(libp2ptls.ID, libp2ptls.New),
		//libp2p.DefaultTransports,
		//libp2p.NATPortMap(),
		/*libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			idht, err = dht.New(ctx, h)
			return idht, err
		}),
		/*libp2p.EnableAutoRelay(),*/
	)

}

func showConnectedListPeers(n *network){
	go func() {
		for {

			fmt.Printf("PubSub: %v\n",n.ps.ListPeers("SignerNodeNetwork"))
			time.Sleep(10 * time.Second)
		}
	}()
}

func NewBootstrapNode(ctx context.Context,config NetConfig) error {
	h,err := newPeerHost(config)

	if err != nil {
		logger.Error(err)
		return err
	}
	logger.Infof("Bootstrap Node will be available at %v/p2p/%s",h.Addrs()[0],h.ID().Pretty())
	discovery := NewDiscovery(ctx,h,config)
	_,err = discovery.SetupDiscovery()

	if err != nil {
		logger.Error(err)
		return err
	}
	return nil

}