package signermanager

import (
	"SignerNode/signermanager/pb"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/jffp113/CryptoProviderSDK/crypto"
	"github.com/jffp113/CryptoProviderSDK/keychain"
	"sync"
)

type permissionedProtocol struct {
	requests map[string]*request
	crypto   crypto.ContextFactory
	keychain keychain.KeyChain
}

type request struct {
	lock         sync.Mutex
	responseChan chan<- []byte
	shares       [][]byte
	t, n         int
	scheme       string
	digest		[]byte
}

func (r *request) AddSig(sig []byte) {
	r.lock.Lock()
	r.shares = append(r.shares, sig)
	r.lock.Unlock()
}

func (p *permissionedProtocol) ProcessMessage(data []byte, ctx processContext) {
	logger.Debug("Received Sign Request, processing.")

	req := pb.ProtocolMessage{}
	proto.Unmarshal(data, &req)

	switch req.Type {
	case pb.ProtocolMessage_SIGN_REQUEST:
		p.processMessageSignRequest(&req, ctx)
	case pb.ProtocolMessage_SIGN_RESPONSE:
		p.processMessageSignResponse(&req, ctx)
	}

}

func (p *permissionedProtocol) processMessageSignResponse(req *pb.ProtocolMessage, ctx processContext) {
	logger.Debug("Received Sign Response")
	//TODO if reply to sign add to a list when enough shares aggregate and send to the blockchain
	signatureMsg := pb.SignResponse{}

	err := proto.Unmarshal(req.Content, &signatureMsg)

	if err != nil {
		//discard share with error
		logger.Error(err)
		return
	}

	//Verify that the sign context belongs to this node
	v, ok := p.requests[signatureMsg.UUID]

	if !ok {
		logger.Debug("Sign response not for here ignoring")
		return //Discard message
	}

	//Append share and verify if have enough shares
	v.AddSig(signatureMsg.Signature)
	if len(v.shares) >= v.t {
		//TODO aggregate signatures and
		fullSig, err := p.aggregateShares(v)

		if err != nil {
			// TODO send error message to the client
			logger.Error(err)
			return
		}
		logger.Infof("Signature was produced: %v",fullSig)
		//TODO send message to the blockchain proxy
		//TODO when the proxy awnsers send a msg to the client
		v.responseChan <- []byte("ok")
	}
}

func (p *permissionedProtocol) processMessageSignRequest(req *pb.ProtocolMessage, ctx processContext) {
	logger.Debug("Received Sign Request")
	reqSign := pb.ClientMessage{}
	err := proto.Unmarshal(req.Content, &reqSign)

	if err != nil {
		logger.Error(err)
		return
	}

	sigShare, err := p.signWithShare(&reqSign)

	if err != nil {
		logger.Error(err)
		return
	}

	resp := pb.SignResponse{
		UUID:      reqSign.UUID,
		Signature: sigShare,
	}

	data, err := proto.Marshal(&resp)

	if err != nil {
		logger.Error(err)
		return
	}

	data, err = createProtocolMessage(data, pb.ProtocolMessage_SIGN_RESPONSE)

	ctx.broadcast(data)
}

func (p *permissionedProtocol) Sign(data []byte, ctx signContext) {
	logger.Infof("Broadcasting %v", string(data))

	req := pb.ClientMessage{}
	err := proto.Unmarshal(data, &req)

	if err != nil {
		//TODO send message to the client
		logger.Error(err)
		return
	}

	p.requests[req.UUID] = &request{
		responseChan: ctx.returnChan,
		shares:       make([][]byte, 0),
		t:            int(req.T),
		n:            int(req.N),
		scheme:       req.Scheme,
		digest:       req.Content,
	}

	signReq, err := createProtocolMessage(data, pb.ProtocolMessage_SIGN_REQUEST)

	if err != nil {
		//TODO send message to the client
		logger.Error(err)
		return
	}

	ctx.broadcast(signReq)

	sigShare, err := p.signWithShare(&req)

	if err != nil {
		//TODO send message to the client
		logger.Error(err)
		return
	}

	request := p.requests[req.UUID]
	request.AddSig(sigShare)
}

func (p *permissionedProtocol) signWithShare(req *pb.ClientMessage) ([]byte, error) {
	keyName := fmt.Sprintf("%v_%v_%v", req.Scheme, req.N, req.T)

	privShare, err := p.keychain.LoadPrivateKey(keyName)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	context := p.crypto.GetSignerVerifierAggregator(req.Scheme)
	b, err := context.Sign(req.Content, privShare)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return b, nil
}

func (p *permissionedProtocol) aggregateShares(req *request) ([]byte, error) {
	keyName := fmt.Sprintf("%v_%v_%v", req.scheme, req.n, req.t)

	pubKey, err := p.keychain.LoadPublicKey(keyName)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	context := p.crypto.GetSignerVerifierAggregator(req.scheme)
	//b, err := context.Sign(req.Content, privShare)

	fullSig, err := context.Aggregate(req.shares,req.digest,pubKey,req.t,req.n)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return fullSig, nil
}

func NewPermissionedProtocol(crypto crypto.ContextFactory, keychain keychain.KeyChain) Protocol {
	return &permissionedProtocol{
		requests: make(map[string]*request),
		crypto:   crypto,
		keychain: keychain,
	}
}

func createProtocolMessage(msg []byte, messageType pb.ProtocolMessage_Type) ([]byte, error) {
	req := pb.ProtocolMessage{
		Type:    messageType,
		Content: msg,
	}

	b, err := proto.Marshal(&req)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return b, err
}
