//Package network defines the a wrapper over LibP2P
//library allowing users of this package to easily
//create a publish/subscriber system over a P2P network
//with discovery functionalities
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
	"go.uber.org/atomic"
	"reflect"
	"sync"
	"time"
	"unsafe"
)

var logger = log.Logger("network")

//BufSize defines the max size of message waiting to be processed
//by a go routine
const BufSize = 1024

//NetConfig is used to configure the p2p network
//gives the possibility to:
//choose other RendezvousString (default: network) to advertise a certain peer;
//set a port where the peer is going to be available;
//Choose the private key which will represent the peer id;
//And set the bootstrap peers to find other peers in the P2P network.
type NetConfig struct {
	RendezvousString string
	BootstrapPeers   []string
	Port             int
	Priv             crypto.PrivKey
	PeerAddress      string
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

	//numberOfRoutines is used to own how many
	//go routines are still using messages chan
	//when this variable reaches zero a go routine
	//closes the channel.
	numberOfRoutines atomic.Int32
}

//A group defines a topic in witch a peer
//is listening and publishing
type Group struct {
	topic    *pubsub.Topic
	sub      *pubsub.Subscription
	messages chan<- []byte
	ctx      context.Context
	cancel   context.CancelFunc
}

type Network interface {
	//Broadcast broadcasts to all available peers
	Broadcast(msg []byte) error
	//BroadcastToGroup, broadcasts to all peers subscribed
	//to a certain group
	BroadcastToGroup(groupId string ,msg []byte) error
	//Send(node string, msg []byte)
	//JoinGroup - Joins a certain broadcast group
	JoinGroup(groupId string) error
	//LeaveGroup - Leaves a certain broadcast group
	LeaveGroup(groupId string) error

	//Receive messages from all subscribed groups and
	//from the default group
	Receive() []byte

	//Get a subset of the peers membership available
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

	sub, err := topic.Subscribe(func(sub *pubsub.Subscription) error {
		//What I'm doing here is a a bit crepy.
		//However the LIBP2P API does not allow me to set the size of the "buffer" chan for publishing messages
		//which in my case is bad. I produce a lot of messages when clients grow.
		//Now the size of the topic channel is 128 instead of 32. Enough for a commodity machine being able to run
		//5 signer nodes with 100 concurrent clients.
		SetUnexportedField(reflect.ValueOf(sub).Elem().FieldByName("ch"),make(chan *pubsub.Message, 128))
		return nil
	})

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

func SetUnexportedField(field reflect.Value, value interface{}) {
	reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).
		Elem().
		Set(reflect.ValueOf(value))
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

	ps, err := pubsub.NewGossipSub(ctx, h,
							pubsub.WithDiscovery(disc),
							pubsub.WithMessageSigning(false), //no need for signing
							pubsub.WithPeerOutboundQueueSize(1024)) //Bigger outbound pool

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

	msgChan := make(chan []byte, BufSize)

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
	n.numberOfRoutines.Inc()
	for {
		logger.Debug("Waiting for new message")
		msg, err := g.sub.Next(g.ctx)

		if err != nil {
			logger.Debug(err)
			n.numberOfRoutines.Dec()
			if n.numberOfRoutines.Load() == 0{
				close(n.messages)
			}
			return
		}
		logger.Debugf("New message arrived from", msg.ReceivedFrom)
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

	listenAddr := fmt.Sprintf("%v%v",config.PeerAddress ,config.Port)
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

//func showConnectedListPeers(n *network) { //TODO to be deleted
//	go func() {
//		for {
//
//			fmt.Printf("PubSub: %v\n", n.ps.ListPeers("SignerNodeNetwork"))
//			time.Sleep(10 * time.Second)
//		}
//	}()
//}

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
