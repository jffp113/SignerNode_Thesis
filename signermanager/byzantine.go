package signermanager

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/jffp113/CryptoProviderSDK/crypto"
	"github.com/jffp113/CryptoProviderSDK/keychain"
	ic "github.com/jffp113/SignerNode_Thesis/interconnect"
	"github.com/jffp113/SignerNode_Thesis/signermanager/pb"
	"github.com/jffp113/SignerNode_Thesis/smartcontractengine"
)

type byzantineProtocol struct {
	crypto           crypto.ContextFactory
	keychain         keychain.KeyChain
	sc               smartcontractengine.SCContextFactory
	broadcastAnswer  bool
}

func (p *byzantineProtocol) Register(interconnect ic.Interconnect) error {
	interconnect.RegisterHandler(ic.SignClientRequest, p.Sign)
	interconnect.RegisterHandler(ic.InstallClientRequest, p.InstallShares)
	interconnect.RegisterHandler(ic.NetworkMessage, p.processMessage)
	return nil
}

func (p *byzantineProtocol) InstallShares(data ic.ICMessage, ctx ic.P2pContext) ic.HandlerResponse {
	return ic.CreateOkMessage([]byte{})
}

func (p *byzantineProtocol) processMessage(msg ic.ICMessage, ctx ic.P2pContext) ic.HandlerResponse {
	logger.Debug("Received sign Request, processing.")

	req := pb.ProtocolMessage{}
	proto.Unmarshal(msg.GetData(), &req)

	switch req.Type {
	case pb.ProtocolMessage_SIGN_REQUEST:
		p.processMessageSignRequest(req.Content, msg.GetFrom(), ctx)
	case pb.ProtocolMessage_SIGN_RESPONSE:
		p.processMessageSignResponse(&req, ctx)
	}
	return ic.CreateOkMessage(msg.GetData())
}

func (p *byzantineProtocol) processMessageSignResponse(req *pb.ProtocolMessage, ctx ic.P2pContext) {
	logger.Debug("Byzantine Do nothing")
}

func (p *byzantineProtocol) processMessageSignRequest(data []byte, from string, ctx ic.P2pContext) {
	logger.Debug("Received sign(Byzantine) Request")
	reqSign := pb.ClientSignMessage{}
	err := proto.Unmarshal(data, &reqSign)

	if err != nil {
		logger.Error(err)
		return
	}

	smartContext, closer := p.sc.GetContext(reqSign.SmartContractAddress)
	defer closer.Close()
	signInfo := smartContext.InvokeSmartContract(data)
	logger.Debugf("SmartContract Execution Result: %v", signInfo)

	if !signInfo.Valid {
		logger.Error("Error Executing SmartContract", signInfo)
		return
	}


	sigShare, err := p.signWithShare(&reqSign, signInfo.Scheme, signInfo.N, signInfo.T)

	if err != nil {
		logger.Error(err)
		return
	}

	respData, err := createSignResponse(reqSign.UUID, sigShare)

	if p.broadcastAnswer {
		ctx.Broadcast(respData)
	} else {
		ctx.Send(respData, from)
	}
}

func (p *byzantineProtocol) signWithShare(req *pb.ClientSignMessage, scheme string, n, t int) ([]byte, error) {
	keyName := fmt.Sprintf("%v_%v_%v", scheme, n, t)

	privShare, err := p.keychain.LoadPrivateKey(keyName)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	//Change content to become byz
	return signWithShare([]byte("I am a byzantine"), privShare, p.crypto, scheme, n, t)
}

func (p *byzantineProtocol) Sign(msg ic.ICMessage, ctx ic.P2pContext) ic.HandlerResponse {
	//Do nothing
	return ic.CreateOkMessage([]byte{})
}

func NewByzantineProtocol(crypto crypto.ContextFactory, keychain keychain.KeyChain,
	sc smartcontractengine.SCContextFactory, broadcastAnswer bool) Protocol {

	return &byzantineProtocol{
		crypto:          crypto,
		keychain:        keychain,
		sc:              sc,
		broadcastAnswer: broadcastAnswer,
	}
}
