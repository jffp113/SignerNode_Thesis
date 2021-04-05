package network

import (
	"context"
	"errors"
	"fmt"
	"github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p"
	connmgr "github.com/libp2p/go-libp2p-connmgr"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"sync"
	"time"
)

var logger = log.Logger("network")

const NetworkBufSize = 128

type NetConfig struct {
	RendezvousString string
	BootstrapPeers   []string
	Port             int
	Priv             crypto.PrivKey
	//ListenAddresses  addrList
	//ProtocolID       string
}

type network struct {
	messages chan []byte

	ctx context.Context
	ps  *pubsub.PubSub

	mainGroup Group

	self peer.ID
	host host.Host
	disc *peerDiscovery

	groupsLock sync.Mutex
	groups     map[string]Group
}

type Group struct {
	topic    *pubsub.Topic
	sub      *pubsub.Subscription
	messages chan<- []byte
	ctx      context.Context
	cancel   context.CancelFunc
}

type Network interface {
	Broadcast(msg []byte) error
	BroadcastToGroup(groupId string ,msg []byte) error
	//Send(node string, msg []byte)
	JoinGroup(groupId string) error
	LeaveGroup(groupId string) error
	Receive() []byte
	GetMembership() []peer.AddrInfo
}

func (n *network) Broadcast(msg []byte) error {
	return n.mainGroup.Broadcast(msg,n.self)
}

func (n *network) BroadcastToGroup(groupId string ,msg []byte) error {
	logger.Infof("Broadcasting to group, %v",groupId)
	n.groupsLock.Lock()
	defer n.groupsLock.Unlock()

	g,ok := n.groups[groupId]

	if !ok {
		return errors.New("group does not exist")
	}

	return g.Broadcast(msg,n.self)
}

func (g Group) Broadcast(msg []byte,self peer.ID) error {
	m := networkMessage{
		From:    self,
		Content: msg,
	}

	msgBytes, err := m.MarshalBinary()

	if err != nil {
		return err
	}

	return g.topic.Publish(g.ctx, msgBytes)
}

func (n *network) JoinGroup(groupId string) error {
	logger.Debugf("Joining group %v", groupId)
	n.groupsLock.Lock()
	defer n.groupsLock.Unlock()
	topic, err := n.ps.Join(groupId)

	if err != nil {
		return err
	}

	sub, err := topic.Subscribe()

	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(n.ctx)

	g := Group{
		topic:    topic,
		sub:      sub,
		messages: n.messages,
		ctx:      ctx,
		cancel: cancel,
	}
	n.groups[groupId] = g

	go g.processIncomingMsg(n)

	return nil
}

func (n *network) LeaveGroup(groupId string) error {
	n.groupsLock.Lock()
	defer n.groupsLock.Unlock()

	g,ok := n.groups[groupId]
	delete(n.groups,groupId)

	if !ok {
		return errors.New("group does not exist")
	}

	g.cancel()

	return nil
}

func (n *network) Receive() []byte {
	//TODO add context
	return <-n.messages
}

func (n *network) GetMembership() []peer.AddrInfo {
	return n.disc.GetPeers()
}

func CreateNetwork(ctx context.Context, config NetConfig) (Network, error) {
	logger.Debug("Setting up Network")
	h, err := newPeerHost(config)
	discovery := NewDiscovery(ctx, h, config)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	logger.Infof("Peer will be available at %v/p2p/%s", h.Addrs()[0], h.ID().Pretty())

	disc, err := discovery.SetupDiscovery()

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	ps, err := pubsub.NewGossipSub(ctx, h, pubsub.WithDiscovery(disc))

	topic, err := ps.Join("SignerNodeNetwork")

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	sub, err := topic.Subscribe()

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	msgChan := make(chan []byte, NetworkBufSize)

	network := &network{
		messages: msgChan,
		ctx:      ctx,
		ps:       ps,
		mainGroup: Group{topic,
			sub,
			msgChan,
			ctx,
			nil},
		self: h.ID(),
		host: h,
		disc: discovery,
		groups: make(map[string]Group),
	}

	go network.mainGroup.processIncomingMsg(network)

	//showConnectedListPeers(network)

	return network, nil
}

func (g Group) processIncomingMsg(n *network) {
	for {

		select {
		case _ = <-g.ctx.Done():
			return
		default:
		}

		logger.Debug("Waiting for new message")
		msg, err := g.sub.Next(g.ctx)
		logger.Debugf("New message arrived from", msg.ReceivedFrom)

		if err != nil {
			logger.Debug(err)
			close(n.messages)
			return
		}
		// only forward messages delivered by others

		if msg.ReceivedFrom == n.self {
			continue
		}

		cm := new(networkMessage)
		err = cm.UnmarshalBinary(msg.Data)

		if err != nil {
			logger.Error(err)
			continue
		}

		// send valid messages onto the Messages channel
		n.messages <- cm.Content
	}
}

func newPeerHost(config NetConfig) (host.Host, error) {

	logger.Debug("Creating Peer Host")

	listenAddr := fmt.Sprintf("/ip4/0.0.0.0/tcp/%v", config.Port)
	return libp2p.New(
		context.Background(),
		libp2p.ListenAddrStrings(listenAddr),
		libp2p.ConnectionManager(connmgr.NewConnManager(
			3,           // Lowwater
			10,          // HighWater,
			time.Minute, // GracePeriod
		)),
		libp2p.Identity(config.Priv),
	)

}

func showConnectedListPeers(n *network) {
	go func() {
		for {

			fmt.Printf("PubSub: %v\n", n.ps.ListPeers("SignerNodeNetwork"))
			time.Sleep(10 * time.Second)
		}
	}()
}

func NewBootstrapNode(ctx context.Context, config NetConfig) error {
	h, err := newPeerHost(config)

	if err != nil {
		logger.Error(err)
		return err
	}
	logger.Infof("Bootstrap Node will be available at %v/p2p/%s", h.Addrs()[0], h.ID().Pretty())
	discovery := NewDiscovery(ctx, h, config)
	_, err = discovery.SetupDiscovery()

	if err != nil {
		logger.Error(err)
		return err
	}
	return nil

}
