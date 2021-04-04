package signermanager

import (
	"github.com/jffp113/SignerNode_Thesis/network"
	"github.com/jffp113/SignerNode_Thesis/signermanager/pb"
	"github.com/jffp113/SignerNode_Thesis/smartcontractengine"
	"github.com/golang/protobuf/proto"
	"github.com/ipfs/go-log"
	"github.com/jffp113/CryptoProviderSDK/client"
	"github.com/jffp113/CryptoProviderSDK/crypto"
	"github.com/jffp113/CryptoProviderSDK/keychain"
	ctx "golang.org/x/net/context"
)

const DefaultNumberOfWorkers = 10
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
	cryptoFactory crypto.ContextFactory
	signerURI     string

	//Context to execute smartcontracts
	scFactory smartcontractengine.SCContextFactory
	scURI     string
	peerPort  int
}

func NewSignerManager(confs ...Config) *signermanager {
	manager := &signermanager{
		workPool:   make(chan []byte, MessageWorkerChanSize),
		numWorkers: DefaultNumberOfWorkers,
		context:    ctx.Background(),
	}

	for _, v := range confs {
		_ = v(manager)
	}

	return manager

}

func (s *signermanager) Init() error {
	net, err := network.CreateNetwork(s.context, network.NetConfig{
		BootstrapPeers: []string{s.bootstrapNode},
		Port:           s.peerPort,
	})
	if err != nil {
		return err
	}
	s.network = net

	factory, err := client.NewCryptoFactory(s.signerURI)

	if err != nil {
		return err
	}

	s.cryptoFactory = factory

	s.keychain = keychain.NewKeyChain(s.keyPath)

	s.scFactory, err = smartcontractengine.NewSmartContractClientFactory(s.scURI)

	if err != nil {
		return err
	}

	p, err := GetProtocol(s.protocolName, s.cryptoFactory, s.keychain, s.scFactory,s.network)

	if err != nil {
		return err
	}

	s.protocol = p

	s.startNetworkReceiver()
	s.startWorkers()

	return nil
}

type signContext struct {
	returnChan chan<- ManagerResponse
	broadcast  func(msg []byte) error
	broadcastToGroup func(groupId string ,msg []byte) error
	joinGroup func(groupId string) error
	leaveGroup func(groupId string) error
}

type processContext struct {
	broadcast func(msg []byte) error
	broadcastToGroup func(groupId string ,msg []byte) error
	joinGroup func(groupId string) error
	leaveGroup func(groupId string) error
}

func (s *signermanager) Sign(data []byte) <-chan ManagerResponse {
	ch := make(chan ManagerResponse, 1) //TODO maybe a pool of protocol workers?
	go s.protocol.Sign(data, signContext{
		returnChan:       ch,
		broadcast:        s.network.Broadcast,
		broadcastToGroup: s.network.BroadcastToGroup,
		joinGroup:        s.network.JoinGroup,
		leaveGroup:       s.network.LeaveGroup,
	})

	return ch
}

func (s *signermanager) InstallShares(data []byte) <-chan ManagerResponse {
	ch := make(chan ManagerResponse, 1) //TODO maybe a pool of protocol workers?
	go func() {
		err := s.protocol.InstallShares(data)
		if err != nil {
			sendErrorMessage(ch,err)
		}
		sendOkMessage(ch,[]byte{})
	}()

	return ch
}

func (s *signermanager) Verify(data []byte) <-chan ManagerResponse {
	ch := make(chan ManagerResponse, 1) //TODO maybe a pool of protocol workers?

	go func() {
		msg := pb.ClientVerifyMessage{}
		err := proto.Unmarshal(data, &msg)

		if err != nil {
			ch <- ManagerResponse{Error,
				nil,
				err}
		}

		context, c := s.cryptoFactory.GetSignerVerifierAggregator(msg.Scheme)
		defer c.Close()
		pubKey := keychain.ConvertBytesToPubKey(msg.PublicKey)
		err = context.Verify(msg.Signature, msg.Digest, pubKey)

		if err != nil {
			createInvalidMessageVerifyResponse(ch)
			return
		}

		createValidMessageVerifyMessages(ch)
	}()

	return ch
}

func (s *signermanager) GetMembership() <-chan ManagerResponse {
	ch := make(chan ManagerResponse, 1) //TODO maybe a pool of protocol workers?
	//TODO

	go func() {
		createValidMembershipResponse(s.network.GetMembership(), ch)
	}()

	return ch
}

func (s *signermanager) startWorkers() {
	for i := 0; i < s.numWorkers; i++ {
		go func() {
			for {
				select {
				case data := <-s.workPool:
					logger.Debug("Start Processing Worker")
					s.protocol.ProcessMessage(data, processContext{
						broadcast:		  s.network.Broadcast,
						broadcastToGroup: s.network.BroadcastToGroup,
						joinGroup:        s.network.JoinGroup,
						leaveGroup:       s.network.LeaveGroup,
					})
					logger.Debug("Finish Processing Worker")
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
			logger.Debug("Waiting for messages to the protocol")
			s.workPool <- s.network.Receive()
			logger.Debug("Received message to the protocol")
		}
	}()
}
