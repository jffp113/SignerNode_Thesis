package signermanager

import (
	"SignerNode/network"
	"github.com/ipfs/go-log"
	"github.com/jffp113/CryptoProviderSDK/client"
	"github.com/jffp113/CryptoProviderSDK/crypto"
	"github.com/jffp113/CryptoProviderSDK/keychain"
	ctx "golang.org/x/net/context"
)

const DefaultNumberOfWorkers = 5
const MessageWorkerChanSize = 20

var logger = log.Logger("protocol")

type signermanager struct {
	// URL for the bootstrap to connect find
	//other participant signer nodes
	bootstrapNode string

	//Path for the private key share
	//and public key
	keyPath  string
	keychain keychain.KeyChain

	//Protocol for choosing who signs
	protocolName string
	protocol     Protocol

	//P2P Network
	network network.Network

	workPool   chan []byte
	numWorkers int

	context ctx.Context

	//Context to the distributed cryptoProvider
	factory   crypto.ContextFactory
	signerURI string
}

func NewSignerManager() *signermanager {
	return &signermanager{
		bootstrapNode: "",
		keyPath:       "",
		protocolName:  "",
		protocol:      nil,
		network:       nil,
		workPool:      make(chan []byte, MessageWorkerChanSize),
		numWorkers:    DefaultNumberOfWorkers,
		context:       ctx.Background(),
		factory:       nil,
		signerURI:     "",
	}
}

func (s *signermanager) SetBootStrapNode(bootstrap string) {
	s.bootstrapNode = bootstrap
}

func (s *signermanager) SetKeyPath(keyPath string) {
	s.keyPath = keyPath
}

func (s *signermanager) SetProtocol(protocol string) {
	s.protocolName = protocol
}

func (s *signermanager) SetSignerURI(uri string) {
	s.signerURI = uri
}

func (s *signermanager) Init() error {
	net, err := network.CreateNetwork(s.context, network.NetConfig{
		BootstrapPeers: []string{s.bootstrapNode},
	})
	if err != nil {
		return err
	}
	s.network = net

	factory, err := client.NewCryptoFactory(s.signerURI)

	if err != nil {
		return err
	}

	s.factory = factory

	s.keychain = keychain.NewKeyChain(s.keyPath)

	p, err := GetProtocol(s.protocolName, s.factory, s.keychain)

	if err != nil {
		return err
	}

	s.protocol = p

	s.startNetworkReceiver()
	s.startWorkers()

	return nil
}

type signContext struct {
	returnChan chan<- []byte
	broadcast  func(msg []byte) error
}

type processContext struct {
	broadcast func(msg []byte) error
}

func (s *signermanager) Sign(data []byte) <-chan []byte {
	ch := make(chan []byte, 1) //TODO maybe a pool of protocol workers?
	go s.protocol.Sign(data, signContext{
		returnChan: ch,
		broadcast:  s.network.Broadcast,
	})

	return ch
}

func (s *signermanager) startWorkers() {
	for i := 0; i < s.numWorkers; i++ {
		go func() {
			for {
				select {
				case data := <-s.workPool:
					s.protocol.ProcessMessage(data, processContext{
						s.network.Broadcast,
					})
				case _ = <-s.context.Done():
					return
				}
			}
		}()
	}
}

func (s *signermanager) startNetworkReceiver() {
	go func() {
		for {
			logger.Debug("Waiting for messages")
			s.workPool <- s.network.Receive()
			logger.Debug("Received message")
		}
	}()
}
