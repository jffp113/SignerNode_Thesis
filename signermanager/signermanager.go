package signermanager

import (
	"github.com/golang/protobuf/proto"
	"github.com/ipfs/go-log/v2"
	"github.com/jffp113/CryptoProviderSDK/client"
	"github.com/jffp113/CryptoProviderSDK/crypto"
	"github.com/jffp113/CryptoProviderSDK/keychain"
	ic "github.com/jffp113/SignerNode_Thesis/interconnect"
	"github.com/jffp113/SignerNode_Thesis/network"
	"github.com/jffp113/SignerNode_Thesis/signermanager/pb"
	"github.com/jffp113/SignerNode_Thesis/smartcontractengine"
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

	//broadcastAnswer controls if the signature share
	//produced by a signer node is broadcaster to the entire group
	//or only to the requester.
	broadcastAnswer bool

	workPool   chan ic.ICMessage
	numWorkers int

	context ctx.Context

	//Context to the distributed cryptoProvider
	cryptoFactory crypto.ContextFactory
	signerURI     string

	//Context to execute smartcontracts
	scFactory   smartcontractengine.SCContextFactory
	scURI       string
	peerPort    int
	peerAddress string

	//Interconnect to talk to a sign manager protocol / verify
	//and membership
	interconnect ic.Interconnect
}

func NewSignerManager(confs ...Config) *signermanager {
	manager := &signermanager{
		workPool:   make(chan ic.ICMessage, MessageWorkerChanSize),
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
		PeerAddress:    s.peerAddress,
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

	p, err := GetProtocol(s.protocolName, s.cryptoFactory, s.keychain, s.scFactory, s.network, s.broadcastAnswer)

	if err != nil {
		return err
	}

	s.protocol = p

	s.interconnect, _ = ic.NewInterconnect(ic.SetContext(ic.P2pContext{
		Broadcast:        s.network.Broadcast,
		BroadcastToGroup: s.network.BroadcastToGroup,
		Send:             s.network.Send,
		JoinGroup:        s.network.JoinGroup,
		LeaveGroup:       s.network.LeaveGroup,
	}))

	s.interconnect.RegisterHandler(ic.VerifyClientRequest, s.verify)
	s.interconnect.RegisterHandler(ic.MembershipClientRequest, s.getMembership)
	s.protocol.Register(s.interconnect)

	s.startNetworkReceiver()
	s.startWorkers()

	return nil
}

//Emit a specific event type (HandlerType) in a defined protocol
//or in the signer manager.
func (s *signermanager) EmitEvent(t ic.HandlerType, content []byte) ic.HandlerResponse {
	return s.interconnect.EmitEvent(t, ic.NewMessageFromBytes(content))
}

func (s *signermanager) verify(m ic.ICMessage, ctx ic.P2pContext) ic.HandlerResponse {
	msg := pb.ClientVerifyMessage{}
	err := proto.Unmarshal(m.GetData(), &msg)

	if err != nil {
		return ic.CreateErrorMessage(err)
	}

	context, c := s.cryptoFactory.GetSignerVerifierAggregator(msg.Scheme)
	defer c.Close()
	pubKey := keychain.ConvertBytesToPubKey(msg.PublicKey)
	err = context.Verify(msg.Signature, msg.Digest, pubKey)

	if err != nil {
		return createInvalidMessageVerifyResponse()
	}

	return createValidMessageVerifyMessages()
}

func (s *signermanager) getMembership(_ ic.ICMessage, _ ic.P2pContext) ic.HandlerResponse {
	return createValidMembershipResponse(s.network.GetMembership())
}

func (s *signermanager) startWorkers() {
	for i := 0; i < s.numWorkers; i++ {
		go func() {
			for {
				select {
				case data := <-s.workPool:
					logger.Debug("Start Processing Worker")
					s.interconnect.EmitEvent(ic.NetworkMessage, data)
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
