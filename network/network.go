package network

import (
	"context"
	"fmt"
	"github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	ma "github.com/multiformats/go-multiaddr"
	"sync"
	"time"
)

type NetConfig struct {
	RendezvousString string
	BootstrapPeers   []string
	Port int
	//ListenAddresses  addrList
	//ProtocolID       string
}

var logger = log.Logger("rendezvous")

func CreateNetwork(ctx context.Context,config NetConfig) error {
	logger.Debug("Setting up Network")
	h,err := newPeerHost(config)

	if err != nil {
		logger.Error(err)
		return err
	}

	logger.Infof("Peer will be available at /ip4/127.0.0.1/tcp/%v/p2p/%s",config.Port,h.ID().Pretty())

	go func() {
		for {
			fmt.Println(h.Peerstore().Peers())
			time.Sleep(10 * time.Second)
		}
	}()

	return setupDiscovery(ctx,h,config)
}

func newPeerHost(config NetConfig) (host.Host, error) {

	logger.Debug("Creating Peer Host")

	listenAddr := fmt.Sprintf("/ip4/0.0.0.0/tcp/%v",config.Port)

	return libp2p.New(
		context.Background(),
		libp2p.ListenAddrStrings(listenAddr),
		//libp2p.Identity(*prvKey),
		//libp2p.Security(libp2ptls.ID, libp2ptls.New),
		//libp2p.DefaultTransports,
		/*libp2p.ConnectionManager(connmgr.NewConnManager(
			10,          // Lowwater
			50,          // HighWater,
			time.Minute, // GracePeriod
		)),
		libp2p.NATPortMap(),*/
		/*libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			idht, err = dht.New(ctx, h)
			return idht, err
		}),
		libp2p.EnableAutoRelay(),*/
	)

}

func setupDiscovery(ctx context.Context, host host.Host,config NetConfig) error {
	logger.Debugf("Setting up Discovery bootstrap nodes:%v",config.BootstrapPeers)

	kademliaDHT, err := dht.New(ctx, host)
	if err != nil {
		return err
	}

	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	for _, peerAddr := range config.BootstrapPeers {
		addr, err := ma.NewMultiaddr(peerAddr)

		if err != nil {
			return err
		}

		peerinfo, _ := peer.AddrInfoFromP2pAddr(addr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := host.Connect(ctx, *peerinfo); err != nil {
				logger.Warn(err)
			} else {
				logger.Info("Connection established with bootstrap node:", *peerinfo)
			}
		}()
	}
	wg.Wait()

	logger.Info("Announcing ourselves...")
	routingDiscovery := discovery.NewRoutingDiscovery(kademliaDHT)
	discovery.Advertise(ctx, routingDiscovery, "network")
	logger.Debug("Successfully announced!")

	if findPeers(host,routingDiscovery) != nil {
		return err
	}

	return nil
}

func findPeers(host host.Host,routingDiscovery *discovery.RoutingDiscovery) error {
	logger.Debug("Searching for other peers...")
	peerChan, err := routingDiscovery.FindPeers(context.Background(), "network")

	if err != nil {
		logger.Error(err)
		return err
	}

	go func() {
		for addr := range peerChan {
			handlePeerFound(host,addr)
		}
	}()

	return nil
}

func handlePeerFound(host host.Host ,pi peer.AddrInfo) {
	logger.Debugf("Discovered new peer %s\n", pi.ID.Pretty())

	err := host.Connect(context.Background(), pi)
	if err != nil {
		logger.Debugf("error connecting to peer %s: %s\n", pi.ID.Pretty(), err)
	}
}
