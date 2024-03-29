package protocol

import (
	"context"
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/jffp113/CryptoProviderSDK/crypto"
	"github.com/jffp113/CryptoProviderSDK/keychain"
	ic "github.com/jffp113/SignerNode_Thesis/interconnect"
	"github.com/jffp113/SignerNode_Thesis/network"
	"github.com/jffp113/SignerNode_Thesis/signermanager/pb"
	"github.com/jffp113/SignerNode_Thesis/smartcontractengine"
	"sync"
	"time"
)

const TimeSlack = time.Second * -10

type permissionlessProtocol struct {
	//requestsLock sync.Mutex
	requests sync.Map //map[string]*request

	installedKeys sync.Map //map[string]*keyInfo

	requestToKeyLock sync.RWMutex
	requestToKey     map[string]string

	crypto  crypto.ContextFactory
	sc      smartcontractengine.SCContextFactory
	network network.Network

	interconnect     ic.Interconnect
	broadcastAnswer  bool

	//deleteStaleReqCh signals a request that has timeout
	deleteStaleReqCh chan string

	//deleteStaleKeyCh signals a key that should be garbage collected
	deleteStaleKeyCh chan string
}

func (p *permissionlessProtocol) Register(interconnect ic.Interconnect) error {
	interconnect.RegisterHandler(ic.SignClientRequest, p.sign)
	interconnect.RegisterHandler(ic.InstallClientRequest, p.installShares)
	interconnect.RegisterHandler(ic.NetworkMessage, p.processMessage)
	p.interconnect = interconnect
	return nil
}

type keyInfo struct {
	privShare    crypto.PrivateKey
	pubKey       crypto.PublicKey
	validUntil   time.Time
	isOneTimeKey bool
	used         bool
}

func (k keyInfo) expired() bool {
	if k.used && k.isOneTimeKey {
		return true
	} else if k.validUntil.Before(time.Now().Add(TimeSlack)) {
		return true
	}
	return false
}

func (p *permissionlessProtocol) addRequest(req *request, uuid string) {
	ctx,cancel := context.WithCancel(context.Background())
	req.ctx = ctx
	req.timer = time.AfterFunc(TimeoutRequestTime, func() {
		cancel()
		p.deleteStaleReqCh <- uuid
	})
	p.requests.Store(uuid,req)
}

func (p *permissionlessProtocol) deleteRequest(uuid string) {
	v, ok := p.requests.LoadAndDelete(uuid)

	if ok {
		req := v.(*request)
		req.timer.Stop()
	}
}

func (p *permissionlessProtocol) getRequest(uuid string) (*request,bool){
	v, ok := p.requests.Load(uuid)

	if !ok {
		return &request{}, ok
	}

	return v.(*request), ok
}

func (p *permissionlessProtocol) addNewRequestToKey(requestUUID, keyID string) {
	p.requestToKeyLock.Lock()
	defer p.requestToKeyLock.Unlock()

	p.requestToKey[requestUUID] = keyID
}

func (p *permissionlessProtocol) getRequestKeyId(requestUUID string) (string, bool) {
	p.requestToKeyLock.RLock()
	defer p.requestToKeyLock.RUnlock()
	v, ok := p.requestToKey[requestUUID]
	return v, ok
}

func (p *permissionlessProtocol) deleteRequestToKey(requestUUID string) {
	p.requestToKeyLock.Lock()
	defer p.requestToKeyLock.Unlock()

	delete(p.requestToKey, requestUUID)
}

func (p *permissionlessProtocol) addInstalledKey(keyID string, info *keyInfo) {
	p.installedKeys.Store(keyID, info)
	timeout := time.Until(info.validUntil)
	time.AfterFunc(timeout, func() {
		p.deleteStaleKeyCh <- keyID
	})

}

func (p *permissionlessProtocol) getInstalledKey(keyId string) (*keyInfo, bool) {
	v, ok := p.installedKeys.Load(keyId)
	return v.(*keyInfo), ok
}

func (p *permissionlessProtocol) deleteInstalledKey(keyId string) {
	p.installedKeys.Delete(keyId)
}

func (r *permissionlessProtocol) AddSigAndTestForEnoughShares(sig []byte, uuid string) (*request, bool) {
	//See if the request exists
	v, ok := r.getRequest(uuid) //TOdo se if concorrent problems

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

func (p *permissionlessProtocol) processMessage(msg ic.ICMessage, ctx ic.P2pContext) ic.HandlerResponse {
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

func (p *permissionlessProtocol) processMessageSignRequest(data []byte, from string, ctx ic.P2pContext) {
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
		ctx.BroadcastToGroup(reqSign.KeyId, respData)
	} else {
		ctx.Send(respData, from)
	}
}

func (p *permissionlessProtocol) processMessageSignResponse(data []byte, ctx ic.P2pContext) {
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
				default:
					//logger.Debug("No new signature shares")
					continue
				}
			}
			logger.Debugf("Aggregating request %v", v)

			firstExec = false
			fullSig, keyId, err := p.aggregateShares(v)

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
			p.deleteRequestToKey(keyId)
			p.deleteRequest(v.uuid)
			close(v.responseChan)
			return
		}
	}
}

func (p *permissionlessProtocol) sign(msg ic.ICMessage, ctx ic.P2pContext) ic.HandlerResponse {
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
	p.addNewRequestToKey(req.UUID, req.KeyId) //TODO memory leak

	signReq, err := createProtocolMessage(msg.GetData(), pb.ProtocolMessage_SIGN_REQUEST)

	if err != nil {
		return ic.CreateErrorMessage(err)
	}

	err = ctx.BroadcastToGroup(req.KeyId, signReq)

	if err != nil {
		return ic.CreateErrorMessage(err)
	}

	sigShare, err := p.signWithShare(&req, request.scheme, request.n, request.t)

	if err != nil {
		return ic.CreateErrorMessage(err)
	}

	respData, err := createSignResponse(request.uuid, sigShare)
	p.interconnect.EmitEvent(ic.NetworkMessage, ic.NewMessageFromBytes(respData))
	return <-respChan
}

func (p *permissionlessProtocol) signWithShare(req *pb.ClientSignMessage, scheme string, n, t int) ([]byte, error) {
	key, exist := p.getInstalledKey(req.KeyId)

	if !exist {
		return nil, errors.New("key does not exist or previously expired")
	}

	if key.expired() {
		return nil, errors.New("key expired")
	}

	sig, err := signWithShare(req.Content, key.privShare, p.crypto, scheme, n, t)

	if err == nil {
		key.used = true
	}

	return sig, err
}

func (p *permissionlessProtocol) aggregateShares(req *request) ([]byte, string, error) {
	keyId, ok := p.getRequestKeyId(req.uuid)

	if !ok {
		return nil, keyId, errors.New("no key set request")
	}

	keyInfo, ok := p.getInstalledKey(keyId)

	if !ok {
		p.deleteRequestToKey(keyId)
		return nil, keyId, errors.New("key does not exist or expired")
	}
	sig, err := aggregateShares(req, keyInfo.pubKey, p.crypto)
	return sig, keyId, err
}

func (p *permissionlessProtocol) installShares(msg ic.ICMessage, ctx ic.P2pContext) ic.HandlerResponse {
	logger.Info("Installing key share")
	request := pb.ClientInstallShareRequest{}
	err := proto.Unmarshal(msg.GetData(), &request)

	if err != nil {
		return ic.CreateErrorMessage(err)
	}

	keyId := hash(request.PublicKey)

	err = p.network.JoinGroup(keyId)

	if err != nil {
		return ic.CreateErrorMessage(err)
	}

	k := keyInfo{
		privShare:    keychain.ConvertBytesToPrivKey(request.PrivateKey),
		pubKey:       keychain.ConvertBytesToPubKey(request.PublicKey),
		validUntil:   time.Unix(request.ValidUntil, 0),
		isOneTimeKey: false,
	}

	p.addInstalledKey(keyId, &k)

	return ic.CreateOkMessage([]byte{})
}

// ShareGarbageCollector removes stale key shares
func (p *permissionlessProtocol) ShareGarbageCollector(deleteStaleKeyCh chan string) {
	for keyId := range deleteStaleKeyCh {
		logger.Info("Removing key share with Id: ", keyId)
		p.network.LeaveGroup(keyId)
		p.deleteInstalledKey(keyId)
	}
}

func NewPermissionlessProtocol(crypto crypto.ContextFactory, sc smartcontractengine.SCContextFactory,
	network network.Network, broadcastAnswer bool) Protocol {

	p := permissionlessProtocol{
		requests:         sync.Map{},//make(map[string]*request),
		installedKeys:    sync.Map{}, //make(map[string]*keyInfo),
		requestToKeyLock: sync.RWMutex{},
		requestToKey:     make(map[string]string),
		crypto:           crypto,
		sc:               sc,
		network:          network,
		broadcastAnswer:  broadcastAnswer,
		deleteStaleReqCh: make(chan string,1),
		deleteStaleKeyCh: make(chan string,1),
	}

	go deleteNoneCompleteRequests(&p.requests, p.deleteStaleReqCh, context.TODO())
	go p.ShareGarbageCollector(p.deleteStaleKeyCh)
	return &p
}
