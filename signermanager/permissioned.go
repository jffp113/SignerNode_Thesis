package signermanager

import (
	"github.com/jffp113/SignerNode_Thesis/signermanager/pb"
	"github.com/jffp113/SignerNode_Thesis/smartcontractengine"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/jffp113/CryptoProviderSDK/crypto"
	"github.com/jffp113/CryptoProviderSDK/keychain"
	"sync"
)

type permissionedProtocol struct {
	requestLock sync.Mutex
	requests    map[string]*request
	crypto      crypto.ContextFactory
	keychain    keychain.KeyChain
	sc          smartcontractengine.SCContextFactory
}

func (p *permissionedProtocol) InstallShares(data []byte) error {
	return errors.New("Operation not supported")
	//TODO change the protocols to a event base where it is necessary to register every
	//TODO operation (message) supported. This will improve code readability
}

func (p *permissionedProtocol) addRequest(req *request, uuid string) {
	p.requestLock.Lock()
	defer p.requestLock.Unlock()
	p.requests[uuid] = req
}

func (p *permissionedProtocol) deleteRequest(uuid string) {
	p.requestLock.Lock()
	defer p.requestLock.Unlock()
	logger.Debugf("Removing request with uuid %v", uuid)
	delete(p.requests, uuid)
}

func (r *permissionedProtocol) AddSigAndTestForEnoughShares(sig []byte, uuid string) (*request, bool) {
	//Lock so no one removes things from requests
	r.requestLock.Lock()
	defer r.requestLock.Unlock()
	//See if the request exists
	v, ok := r.requests[uuid]

	if !ok {
		//Does not exist returns false
		//This can mean to things, request already fulfilled or
		//request is not for this signer node
		return nil, false
	}

	enoughShares := v.AddSigAndCheckIfHaveEnoughShares(sig)

	if enoughShares {
		aggregatingInProgress := v.aggregatingInProgress.Swap(true)
		return v, !aggregatingInProgress
	}

	return v, false
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
	if v, ready := p.AddSigAndTestForEnoughShares(signatureMsg.Signature, signatureMsg.UUID); ready {
		firstExec := true
		//Todo add context with timer
		for {
			if !firstExec {
				select {
				case newSig := <-v.sharesChan:
					v.shares = append(v.shares, newSig)
				default:
					logger.Debug("No new signature shares")
					continue
				}
			}
			logger.Debugf("Aggregating request %v", v)

			firstExec = false
			fullSig, err := p.aggregateShares(v)

			if err != nil {
				//sendErrorMessage(v.responseChan, err)
				logger.Error(err)
				continue
			}
			logger.Debugf("Signature was produced: %v", fullSig)

			resp := pb.ClientSignResponse{
				Scheme:    v.scheme,
				Signature: fullSig,
			}
			bytes, err := proto.Marshal(&resp)

			sendOkMessage(v.responseChan, bytes)
			p.deleteRequest(v.uuid)
			close(v.responseChan)
			return
		}
	}
}

func (p *permissionedProtocol) processMessageSignRequest(req *pb.ProtocolMessage, ctx processContext) {
	logger.Debug("Received Sign Request")
	reqSign := pb.ClientSignMessage{}
	err := proto.Unmarshal(req.Content, &reqSign)

	if err != nil {
		logger.Error(err)
		return
	}

	smartContext, closer := p.sc.GetContext(reqSign.SmartContractAddress)
	defer closer.Close()
	signInfo := smartContext.InvokeSmartContract(req.Content)
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

	req := pb.ClientSignMessage{}
	err := proto.Unmarshal(data, &req)

	if err != nil {
		sendErrorMessage(ctx.returnChan, err)
		logger.Error(err)
		return
	}

	smartContext, closer := p.sc.GetContext(req.SmartContractAddress)
	defer closer.Close()
	signInfo := smartContext.InvokeSmartContract(req.Content)
	logger.Debugf("SmartContract Execution Result: %v", signInfo)

	if signInfo.Error {
		sendErrorMessage(ctx.returnChan, errors.New("error executing smartcontract"))
		return
	}

	if !signInfo.Valid {
		sendInvalidTransactionMessage(ctx.returnChan)
		return
	}

	request := &request{
		responseChan: ctx.returnChan,
		shares:       make([][]byte, 0),
		sharesChan:   make(chan []byte, signInfo.N),
		t:            signInfo.T,
		n:            signInfo.N,
		scheme:       signInfo.Scheme,
		digest:       req.Content,
		uuid:         req.UUID,
	}

	p.addRequest(request, req.UUID)

	signReq, err := createProtocolMessage(data, pb.ProtocolMessage_SIGN_REQUEST)

	if err != nil {
		sendErrorMessage(request.responseChan, err)
		logger.Error(err)
		return
	}

	ctx.broadcast(signReq)

	sigShare, err := p.signWithShare(&req, request.scheme, request.n, request.t)

	if err != nil {
		sendErrorMessage(request.responseChan, err)
		logger.Error(err)
		return
	}

	request.AddSigAndCheckIfHaveEnoughShares(sigShare) //TODO concurrency error until now not triggered
}

func (p *permissionedProtocol) signWithShare(req *pb.ClientSignMessage, scheme string, n, t int) ([]byte, error) {
	keyName := fmt.Sprintf("%v_%v_%v", scheme, n, t)

	privShare, err := p.keychain.LoadPrivateKey(keyName)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return signWithShare(req.Content, privShare, p.crypto, scheme, n, t)
}

func (p *permissionedProtocol) aggregateShares(req *request) ([]byte, error) {
	keyName := fmt.Sprintf("%v_%v_%v", req.scheme, req.n, req.t)

	pubKey, err := p.keychain.LoadPublicKey(keyName)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return aggregateShares(req, pubKey, p.crypto)
}

func NewPermissionedProtocol(crypto crypto.ContextFactory, keychain keychain.KeyChain,
	sc smartcontractengine.SCContextFactory) Protocol {

	p := permissionedProtocol{
		requests: make(map[string]*request),
		crypto:   crypto,
		keychain: keychain,
		sc:       sc,
	}

	return &p
}
