package network

import (
	"context"
	"github.com/libp2p/go-libp2p-core/discovery"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	discoverydht "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	ma "github.com/multiformats/go-multiaddr"
	"sync"
)

//peerDiscovery provides the core functionality to find peers
//in a P2P network. Discovery is supported by LibP2P discovery
//and uses Kademlia DHT to support the discovery.
type peerDiscovery struct {
	ctx context.Context
	//The host represents a P2P node
	host      host.Host
	config    NetConfig
	cancel    context.CancelFunc
	discovery discovery.Discovery
}

//Creates a new discovery service to find new peers
func NewDiscovery(ctx context.Context, host host.Host, config NetConfig) *peerDiscovery {
	ctx, cancel := context.WithCancel(ctx)
	return &peerDiscovery{
		ctx,
		host,
		config,
		cancel,
		nil,
	}
}

//Advertise the peer in the network channel and start connecting to peers
func (d *peerDiscovery) SetupDiscovery() (discovery.Discovery, error) {
	logger.Debugf("Setting up Discovery bootstrap nodes:%v", d.config.BootstrapPeers)

	kademliaDHT, err := dht.New(d.ctx, d.host)
	if err != nil {
		return nil, err
	}

	if err = kademliaDHT.Bootstrap(d.ctx); err != nil {
		return nil, err
	}

	if err = d.connectToBootstrapNodes(); err != nil {
		return nil, err
	}

	logger.Info("Announcing ourselves...")
	routingDiscovery := discoverydht.NewRoutingDiscovery(kademliaDHT)

	discoverydht.Advertise(d.ctx, routingDiscovery, "network")

	d.discovery = routingDiscovery
	return routingDiscovery, nil
}

//Establish a connection to a bootstrap peer and start finding peers
func (d *peerDiscovery) connectToBootstrapNodes() error {
	var wg sync.WaitGroup
	for _, peerAddr := range d.config.BootstrapPeers {
		addr, err := ma.NewMultiaddr(peerAddr)

		if err != nil {
			return err
		}

		peerinfo, _ := peer.AddrInfoFromP2pAddr(addr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := d.host.Connect(d.ctx, *peerinfo); err != nil {
				logger.Warn(err)
			} else {
				logger.Info("Connection established with bootstrap node:", *peerinfo)
			}
		}()
	}
	wg.Wait()

	return nil
}

//Get the list of connected bootstrap peers
func (d *peerDiscovery) GetBootstrapPeers() []ma.Multiaddr {
	bcap := len(d.config.BootstrapPeers)
	result := make([]ma.Multiaddr, 0, bcap)

	for _, v := range d.config.BootstrapPeers {
		addr, _ := ma.NewMultiaddr(v)
		result = append(result, addr)
	}

	return result
}

//Get the list of connected peers (including bootstrap peers)
func (d *peerDiscovery) GetPeers() []peer.AddrInfo {
	peerChan, err := d.discovery.FindPeers(context.Background(), "network")

	var peers []peer.AddrInfo

	logger.Debug("Discovering peers", err)
	for v := range peerChan {
		logger.Debug(v)
		peers = append(peers, v)
	}

	return peers
}

//func (d *peerDiscovery) findPeers(routingDiscovery *discoverydht.RoutingDiscovery) error {
//	discoverydht.Advertise(d.ctx, routingDiscovery, "network")
//	logger.Debug("Successfully announced!")
//
//	logger.Debug("Searching for other peers...")
//	peerChan, err := routingDiscovery.FindPeers(context.Background(), "network")
//
//	if err != nil {
//		logger.Error(err)
//		return err
//	}
//
//	go func() {
//		for addr := range peerChan {
//			handlePeerFound(d.host, addr)
//		}
//		logger.Info("Finished Searching")
//	}()
//
//	return nil
//}

//func handlePeerFound(host host.Host, pi peer.AddrInfo) {
//	logger.Debugf("Discovered new peer %s\n", pi.ID.Pretty())
//
//	err := host.Connect(context.Background(), pi)
//	if err != nil {
//		logger.Debugf("error connecting To peer %s: %s\n", pi.ID.Pretty(), err)
//	}
//}

//func showConnectedPeers(host host.Host) {
//	go func() {
//		for {
//			fmt.Println(host.Peerstore().Peers())
//			time.Sleep(10 * time.Second)
//		}
//	}()
//}
