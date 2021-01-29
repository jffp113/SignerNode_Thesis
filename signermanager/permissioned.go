package signermanager

import (
	"SignerNode/signermanager/pb"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/jffp113/CryptoProviderSDK/crypto"
	"github.com/jffp113/CryptoProviderSDK/keychain"
	"sync"
	"time"
)

type permissionedProtocol struct {
	requestLock sync.Mutex
	requests    map[string]*request
	crypto      crypto.ContextFactory
	keychain    keychain.KeyChain
}

func (p *permissionedProtocol) addRequest(req *request, uuid string) {
	p.requestLock.Lock()
	p.requests[uuid] = req
	p.requestLock.Unlock()
}

type request struct {
	lock         sync.Mutex
	responseChan chan<- []byte
	shares       [][]byte
	t, n         int
	scheme       string
	uuid string
	digest		[]byte
}

func (r *request) AddSig(sig []byte) {
	r.lock.Lock()
	r.shares = append(r.shares,sig)
	r.lock.Unlock()
}

func (r *permissionedProtocol) AddSigTestAndRemoveFromRequests(sig []byte, uuid string) (*request,bool){
	//Lock so no one removes things from requests
	r.requestLock.Lock()
	defer r.requestLock.Unlock()
	//See if the request exists
	v,ok := r.requests[uuid]

	if !ok {
		//Does not exist returns false
		//This can mean to things, request already fulfilled or
		//request is not for this signer node
		return nil,false
	}

	//Lock request so no one changes the shares
	v.lock.Lock()
	defer v.lock.Unlock()
	//Add the share
	v.shares = append(v.shares, sig)

	enoughShares := len(v.shares) >= v.t

	if enoughShares {
		logger.Debugf("Removing request with uuid %v",uuid)
		delete(r.requests,uuid)
	}


	return v, enoughShares
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
	signatureMsg := pb.SignResponse{}

	err := proto.Unmarshal(req.Content, &signatureMsg)

	if err != nil {
		//discard share with error
		logger.Error(err)
		return
	}

	//Check if can aggregate if yes start other
	if v, ready := p.AddSigTestAndRemoveFromRequests(signatureMsg.Signature,signatureMsg.UUID); ready {

		logger.Debugf("Aggregating request %v",v)
		fullSig, err := p.aggregateShares(v)

		if err != nil {
			// TODO send error message to the client
			logger.Error(err)
			return
		}
		logger.Debugf("Signature was produced: %v",fullSig)
		//TODO send message to the blockchain proxy
		//TODO when the proxy awnsers send a msg to the client
		v.responseChan <- []byte("ok")
		close(v.responseChan)
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

	request := &request{
		responseChan: ctx.returnChan,
		shares:       make([][]byte, 0),
		t:            int(req.T),
		n:            int(req.N),
		scheme:       req.Scheme,
		digest:       req.Content,
		uuid: req.UUID,
	}
	p.addRequest(request, req.UUID)

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

	request.AddSig(sigShare)
}

func (p *permissionedProtocol) signWithShare(req *pb.ClientMessage) ([]byte, error) {
	keyName := fmt.Sprintf("%v_%v_%v", req.Scheme, req.N, req.T)

	privShare, err := p.keychain.LoadPrivateKey(keyName)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	context, closer := p.crypto.GetSignerVerifierAggregator(req.Scheme)
	defer closer.Close()
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

	context,closer := p.crypto.GetSignerVerifierAggregator(req.scheme)
	defer closer.Close()
	//b, err := context.Sign(req.Content, privShare)

	fullSig, err := context.Aggregate(req.shares,req.digest,pubKey,req.t,req.n)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return fullSig, nil
}

func NewPermissionedProtocol(crypto crypto.ContextFactory, keychain keychain.KeyChain) Protocol {
	p := permissionedProtocol{
		requests: make(map[string]*request),
		crypto:   crypto,
		keychain: keychain,
	}

	go func() {//TODO remove in the final version
		for {
			time.Sleep(10 * time.Second)
			fmt.Println("Printing Requests")
			for _,v := range p.requests {
				fmt.Println(v)
			}
		}
	}()

	return &p
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
