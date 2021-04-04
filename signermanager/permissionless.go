package signermanager

import (
	"github.com/jffp113/SignerNode_Thesis/network"
	"github.com/jffp113/SignerNode_Thesis/signermanager/pb"
	"github.com/jffp113/SignerNode_Thesis/smartcontractengine"
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/jffp113/CryptoProviderSDK/crypto"
	"github.com/jffp113/CryptoProviderSDK/keychain"
	"sync"
	"time"
)

const TimeSlack = time.Second * -10

type permissionlessProtocol struct {
	requestsLock sync.Mutex //TODO see if sync.map is better for lock contention
	requests     map[string]*request

	installedKeys     sync.Map//map[string]*keyInfo

	requestToKeyLock sync.RWMutex
	requestToKey     map[string]string

	crypto  crypto.ContextFactory
	sc      smartcontractengine.SCContextFactory
	network network.Network
}


type keyInfo struct {
	privShare    crypto.PrivateKey
	pubKey       crypto.PublicKey
	validUntil   time.Time
	isOneTimeKey bool
	used         bool
}

func (k keyInfo) expired() bool {
	if k.used && k.isOneTimeKey{
		return true
	} else if  k.validUntil.Before(time.Now().Add(TimeSlack)) {
		return true
	}
	return false
}

func (p *permissionlessProtocol) addRequest(req *request, uuid string) {
	p.requestsLock.Lock()
	defer p.requestsLock.Unlock()
	p.requests[uuid] = req
}

func (p *permissionlessProtocol) deleteRequest(uuid string) {
	p.requestsLock.Lock()
	defer p.requestsLock.Unlock()
	logger.Debugf("Removing request with uuid %v", uuid)
	delete(p.requests, uuid)
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
	p.installedKeys.Store(keyID,info)
}

func (p *permissionlessProtocol) getInstalledKey(keyId string) (*keyInfo, bool) {
	v,ok := p.installedKeys.Load(keyId)
	return  v.(*keyInfo), ok
}

func (p *permissionlessProtocol) deleteInstalledKey(keyId string) {
	p.installedKeys.Delete(keyId)
}

func (r *permissionlessProtocol) AddSigAndTestForEnoughShares(sig []byte, uuid string) (*request, bool) {
	//Lock so no one removes things from requests
	r.requestsLock.Lock()
	defer r.requestsLock.Unlock()
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

func (p *permissionlessProtocol) ProcessMessage(data []byte, ctx processContext) {
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

func (p *permissionlessProtocol) processMessageSignRequest(req *pb.ProtocolMessage, ctx processContext) {
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

	ctx.broadcastToGroup(reqSign.KeyId,data)
}

func (p *permissionlessProtocol) processMessageSignResponse(req *pb.ProtocolMessage, ctx processContext) {
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
			fullSig,keyId, err := p.aggregateShares(v)

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
			p.deleteRequestToKey(keyId)
			p.deleteRequest(v.uuid)
			close(v.responseChan)
			return
		}
	}
}

func (p *permissionlessProtocol) Sign(data []byte, ctx signContext) {
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
	p.addNewRequestToKey(req.UUID, req.KeyId) //TODO memory leak

	signReq, err := createProtocolMessage(data, pb.ProtocolMessage_SIGN_REQUEST)

	if err != nil {
		sendErrorMessage(request.responseChan, err)
		logger.Error(err)
		return
	}

	err = ctx.broadcastToGroup(req.KeyId,signReq) //TODO broadcast to the group

	if err != nil {
		sendErrorMessage(request.responseChan, err)
		logger.Error(err)
		return
	}

	sigShare, err := p.signWithShare(&req, request.scheme, request.n, request.t)

	if err != nil {
		sendErrorMessage(request.responseChan, err)
		logger.Error(err)
		return
	}

	request.AddSigAndCheckIfHaveEnoughShares(sigShare) //TODO concurrency error until now not triggered
}

func (p *permissionlessProtocol) signWithShare(req *pb.ClientSignMessage, scheme string, n, t int) ([]byte, error) {
	key, exist := p.getInstalledKey(req.KeyId)

	if exist {
		return nil, errors.New("key does not exist or previously expired")
	}

	if key.expired() {
		return nil, errors.New("key expired")
	}


	sig, err := signWithShare(req.Content, key.privShare, p.crypto, scheme, n, t)

	if err == nil {
		key.used = true
	}

	return sig,err
}

func (p *permissionlessProtocol) aggregateShares(req *request) ([]byte,string, error) {
	keyId,ok := p.getRequestKeyId(req.uuid)

	if !ok {
		return nil,keyId, errors.New("no key set request")
	}

	keyInfo, ok := p.getInstalledKey(keyId)

	if !ok {
		p.deleteRequestToKey(keyId)
		return nil,keyId, errors.New("key does not exist or expired")
	}
	sig,err := aggregateShares(req, keyInfo.pubKey, p.crypto)
	return sig,keyId,err
}

func (p *permissionlessProtocol) InstallShares(data []byte) error {
	request := pb.ClientInstallShareRequest{}
	err := proto.Unmarshal(data,&request)

	if err != nil {
		return err
	}

	keyId := hash(request.PublicKey)

	err = p.network.JoinGroup(keyId)

	if err != nil {
		return err
	}

	k := keyInfo{
		privShare:    keychain.ConvertBytesToPrivKey(request.PrivateKey),
		pubKey:       keychain.ConvertBytesToPubKey(request.PublicKey),
		validUntil:   time.Unix(request.ValidUntil,0),
		isOneTimeKey: false,
	}
	
	p.addInstalledKey(keyId,&k)

	return nil
}

func (p *permissionlessProtocol) ShareGarbageCollector() {
	var toDeleteKeys []string

	appendToDelete := func(toDelete interface{}) {
		toDeleteKeys = append(toDeleteKeys, toDelete.(string))
	}

	p.installedKeys.Range(func(key, value interface{}) bool {
		keyInfo := value.(*keyInfo)

		if keyInfo.expired(){
			appendToDelete(key)
		}

		return true
	})

	for _,v := range toDeleteKeys {
		p.network.LeaveGroup(v)
		p.deleteInstalledKey(v)
	}

}

func NewPermissionlessProtocol(crypto crypto.ContextFactory, sc smartcontractengine.SCContextFactory, network network.Network) Protocol {

	p := permissionlessProtocol{
		requestsLock:      sync.Mutex{},
		requests:          make(map[string]*request),
		installedKeys:     sync.Map{},//make(map[string]*keyInfo),
		requestToKeyLock:  sync.RWMutex{},
		requestToKey:      make(map[string]string),
		crypto:            crypto,
		sc:                sc,
		network: network,
	}

	return &p
}
