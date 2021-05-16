package signermanager

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/jffp113/CryptoProviderSDK/crypto"
	"github.com/jffp113/CryptoProviderSDK/keychain"
	ic "github.com/jffp113/SignerNode_Thesis/interconnect"
	"github.com/jffp113/SignerNode_Thesis/signermanager/pb"
	"github.com/jffp113/SignerNode_Thesis/smartcontractengine"
	"sync"
	"time"
)

//permissionedProtocol represents a protocol where the participants
//are known in advance and every key share was previously installed
type permissionedProtocol struct {
	requests         sync.Map //map[string]*request
	crypto           crypto.ContextFactory
	keychain         keychain.KeyChain
	sc               smartcontractengine.SCContextFactory
	interconnect     ic.Interconnect
	broadcastAnswer  bool
	deleteStaleReqCh chan string
}

func (p *permissionedProtocol) Register(inter ic.Interconnect) error {
	inter.RegisterHandler(ic.SignClientRequest, p.sign)
	inter.RegisterHandler(ic.NetworkMessage, p.processMessage)
	p.interconnect = inter
	return nil
}

func (p *permissionedProtocol) addRequest(req *request, uuid string) {
	ctx,cancel := context.WithCancel(context.Background())
	req.ctx = ctx
	req.timer = time.AfterFunc(TimeoutRequestTime, func() {
		cancel()
		p.deleteStaleReqCh <- uuid
	})
	p.requests.Store(uuid, req)
}

func (p *permissionedProtocol) deleteRequest(uuid string) {
	logger.Debugf("Removing request with uuid %v", uuid)
	v, ok := p.requests.LoadAndDelete(uuid)

	if ok {
		req := v.(*request)
		req.timer.Stop()
	}
}

func (p *permissionedProtocol) getRequest(uuid string) (*request, bool) {
	v, ok := p.requests.Load(uuid)

	if !ok {
		return &request{}, ok
	}

	return v.(*request), ok
}

func (r *permissionedProtocol) AddSigAndTestForEnoughShares(sig []byte, uuid string) (*request, bool) {
	//See if the request exists
	v, ok := r.getRequest(uuid)//TODO see if there is a concorrency problem

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

func (p *permissionedProtocol) processMessage(msg ic.ICMessage, ctx ic.P2pContext) ic.HandlerResponse {
	logger.Debug("Received sign Request, processing.")

	req := pb.ProtocolMessage{}
	proto.Unmarshal(msg.GetData(), &req)

	switch req.Type {
	case pb.ProtocolMessage_SIGN_REQUEST:
		p.processMessageSignRequest(req.Content, msg.GetFrom(), ctx)
	case pb.ProtocolMessage_SIGN_RESPONSE:
		p.processMessageSignResponse(req.Content, ctx)
	}
	return ic.CreateOkMessage(msg.GetData())
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
		for {
			if !firstExec {
				select {
				case newSig := <-v.sharesChan:
					v.shares = append(v.shares, newSig)
				case  <-v.ctx.Done():
					//It was canceled due a timeout
					logger.Debug("Signature aggregation cancelled due timeout")
					return
				default:
					//logger.Debug("No new signature shares")
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

func (p *permissionedProtocol) processMessageSignRequest(data []byte, from string, ctx ic.P2pContext) {
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

	respData, err := createSignResponse(reqSign.UUID, sigShare)

	if p.broadcastAnswer {
		ctx.Broadcast(respData)
	} else {
		ctx.Send(respData, from)
	}
}

func (p *permissionedProtocol) sign(msg ic.ICMessage, ctx ic.P2pContext) ic.HandlerResponse {
	logger.Infof("Broadcasting %v", string(msg.GetData()))

	req := pb.ClientSignMessage{}
	err := proto.Unmarshal(msg.GetData(), &req)

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

	respChan := make(chan ic.HandlerResponse, 1)
	request := &request{
		responseChan: respChan,
		shares:       make([][]byte, 0),
		sharesChan:   make(chan []byte, signInfo.N),
		t:            signInfo.T,
		n:            signInfo.N,
		scheme:       signInfo.Scheme,
		digest:       req.Content,
		uuid:         req.UUID,
	}

	p.addRequest(request, req.UUID)

	signReq, err := createProtocolMessage(msg.GetData(), pb.ProtocolMessage_SIGN_REQUEST)

	if err != nil {
		ic.CreateErrorMessage(err)
	}

	ctx.Broadcast(signReq)

	sigShare, err := p.signWithShare(&req, request.scheme, request.n, request.t)

	if err != nil {
		return ic.CreateErrorMessage(err)
	}

	respData, err := createSignResponse(request.uuid, sigShare)
	p.interconnect.EmitEvent(ic.NetworkMessage, ic.NewMessageFromBytes(respData))
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

//aggregateShares receives a request and tries to aggregate the signature shares
//produced by the different signer nodes.
func (p *permissionedProtocol) aggregateShares(req *request) ([]byte, error) {
	keyName := fmt.Sprintf("%v_%v_%v", req.scheme, req.n, req.t)

	pubKey, err := p.keychain.LoadPublicKey(keyName)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return aggregateShares(req, pubKey, p.crypto)
}

//NewPermissionedProtocol creates a new permissioned protocol with a specific crypto context,
//keychain, smartcontract proxy factory and passing a boolean (broadcastAnswer) to indicate
//if the replies should be broadcasted to all signer nodes.
func NewPermissionedProtocol(crypto crypto.ContextFactory, keychain keychain.KeyChain,
	sc smartcontractengine.SCContextFactory, broadcastAnswer bool) Protocol {

	p := permissionedProtocol{
		requests:         sync.Map{}, //make(map[string]*request),
		crypto:           crypto,
		keychain:         keychain,
		sc:               sc,
		broadcastAnswer:  broadcastAnswer,
		deleteStaleReqCh: make(chan string, 1),
	}

	go deleteNoneCompleteRequests(&p.requests, p.deleteStaleReqCh, context.TODO())

	return &p
}

