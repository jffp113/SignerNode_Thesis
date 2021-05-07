package signermanager

import (
	ic "github.com/jffp113/SignerNode_Thesis/interconnect"
	"github.com/jffp113/SignerNode_Thesis/signermanager/pb"
	"github.com/jffp113/SignerNode_Thesis/smartcontractengine"
	"errors"
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
	sc          smartcontractengine.SCContextFactory
	interconnect ic.Interconnect
}

func (p *permissionedProtocol) Register(inter ic.Interconnect) error {
	inter.RegisterHandler(ic.SignClientRequest,p.sign)
	inter.RegisterHandler(ic.NetworkMessage,p.processMessage)
	p.interconnect = inter
	return nil
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

func (p *permissionedProtocol) processMessage(data []byte, ctx ic.P2pContext) ic.HandlerResponse{
	logger.Debug("Received sign Request, processing.")

	req := pb.ProtocolMessage{}
	proto.Unmarshal(data, &req)

	switch req.Type {
	case pb.ProtocolMessage_SIGN_REQUEST:
		p.processMessageSignRequest(req.Content, ctx)
	case pb.ProtocolMessage_SIGN_RESPONSE:
		p.processMessageSignResponse(req.Content, ctx)
	}
	return ic.CreateOkMessage(data)
}

func (p *permissionedProtocol) processMessageSignResponse(data []byte, ctx ic.P2pContext) {
	logger.Debug("Received sign Response")
	signatureMsg := pb.SignResponse{}

	err := proto.Unmarshal(data, &signatureMsg)

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

			ic.SendOkMessage(v.responseChan, bytes)
			p.deleteRequest(v.uuid)
			close(v.responseChan)
			return
		}
	}
}

func (p *permissionedProtocol) processMessageSignRequest(data []byte, ctx ic.P2pContext) {
	logger.Debug("Received sign Request")
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

	respData, err := createSignResponse(reqSign.UUID,sigShare)

	ctx.Broadcast(respData)
}

func (p *permissionedProtocol) sign(data []byte, ctx ic.P2pContext) ic.HandlerResponse {
	logger.Infof("Broadcasting %v", string(data))

	req := pb.ClientSignMessage{}
	err := proto.Unmarshal(data, &req)

	if err != nil {
		return ic.CreateErrorMessage(err)
	}

	smartContext, closer := p.sc.GetContext(req.SmartContractAddress)
	defer closer.Close()
	signInfo := smartContext.InvokeSmartContract(req.Content)
	logger.Debugf("SmartContract Execution Result: %v", signInfo)

	if signInfo.Error {
		return ic.CreateErrorMessage(errors.New("error executing smartcontract"))
	}

	if !signInfo.Valid {
		return ic.CreateInvalidTransactionMessage()
	}

	respChan := make(chan ic.HandlerResponse,1)
	request := &request{
		responseChan: respChan,
		shares:       make([][]byte, 0),
		sharesChan:   make(chan []byte, signInfo.N),
		t:            signInfo.T,
		n:            signInfo.N,
		scheme:       signInfo.Scheme,
		digest:       req.Content,
		uuid:         req.UUID,
		submitTime:   time.Now(),
	}

	p.addRequest(request, req.UUID)

	signReq, err := createProtocolMessage(data, pb.ProtocolMessage_SIGN_REQUEST)

	if err != nil {
		ic.CreateErrorMessage(err)
	}

	ctx.Broadcast(signReq)

	sigShare, err := p.signWithShare(&req, request.scheme, request.n, request.t)

	if err != nil {
		return ic.CreateErrorMessage(err)
	}

	respData, err := createSignResponse(request.uuid,sigShare)
	p.interconnect.EmitEvent(ic.NetworkMessage,respData)
	return <-respChan
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

	//toDO DELETE
	//go func(){
	//	for {
	//		time.Sleep(2 * time.Second)
	//		fmt.Println("Requests waiting:")
	//		for _,v := range p.requests {
	//			fmt.Println(v)
	//			if v.submitTime.Before(time.Now().Add(2*time.Second)){
	//				p.deleteRequest(v.uuid)
	//				ic.SendErrorMessage(v.responseChan,errors.New("timeout"))
	//				fmt.Println("Removing request")
	//			}
	//
	//		}
	//	}
	//}()

	return &p
}
