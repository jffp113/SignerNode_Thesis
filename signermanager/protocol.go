package signermanager

import (
	"context"
	"errors"
	"github.com/jffp113/CryptoProviderSDK/crypto"
	"github.com/jffp113/CryptoProviderSDK/keychain"
	ic "github.com/jffp113/SignerNode_Thesis/interconnect"
	"github.com/jffp113/SignerNode_Thesis/network"
	"github.com/jffp113/SignerNode_Thesis/smartcontractengine"
	"go.uber.org/atomic"
	"sync"
	"time"
)

const PERMISSIONED = "Permissioned"
const PERMISSIONLESS = "Permissionless"
const BYZANTINE = "Byzantine"

//TODO make this as default and give th possibility to configure other timeouts
const TimeoutRequestTime = 100 * time.Second

type Protocol interface {
	Register(ic ic.Interconnect) error
}

func GetProtocol(protocolName string, factory crypto.ContextFactory,
	keychain keychain.KeyChain, scFactory smartcontractengine.SCContextFactory, network network.Network, broadcastAnswer bool) (Protocol, error) {
	switch protocolName {
	case PERMISSIONED:
		return NewPermissionedProtocol(factory, keychain, scFactory, broadcastAnswer), nil
	case PERMISSIONLESS:
		return NewPermissionlessProtocol(factory, scFactory, network, broadcastAnswer), nil
	case BYZANTINE:
		return NewByzantineProtocol(), nil
	default:
		return nil, errors.New("protocol does not exist")
	}

}

type request struct {
	//Lock necessary to control insertion
	//in the signature shares slice
	lock sync.Mutex
	//Chan to respond to the client
	responseChan chan<- ic.HandlerResponse

	//Signature shares from every signernode
	shares                [][]byte
	sharesChan            chan []byte
	aggregatingInProgress atomic.Bool
	insertInSharesChan    bool

	//request information to build the signature
	t, n                  int
	scheme                string
	digest                []byte

	//uuid to identify unequivocally a request
	uuid                  string

	//fields to allow to cancel a request when to much
	//time has elapsed
	timer                 *time.Timer
	ctx						context.Context
}

//Will return true when has enough shares at the first time
func (r *request) AddSigAndCheckIfHaveEnoughShares(sig []byte) bool {
	//Lock request so no one changes the shares
	r.lock.Lock()
	defer r.lock.Unlock()

	//Check if shares slice was used
	if r.insertInSharesChan {
		//If used insert in the chan so that the goroutine
		//who is aggregating can get the new shares
		r.sharesChan <- sig
		return true //Only return true at the first time
	}

	//Add the share in the slice
	r.shares = append(r.shares, sig)
	r.insertInSharesChan = len(r.shares) >= r.t
	return r.insertInSharesChan
}

//DeleteNoneCompleteRequests deletes requests that weren't fulfilled in a predefined time.
func deleteNoneCompleteRequests(requests *sync.Map, deleteCh <-chan string, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case key := <-deleteCh:
			v, ok := requests.LoadAndDelete(key)
			if ok {
				req := v.(*request)
				logger.Error("Request timeout for: ", req.uuid)
				ic.SendErrorMessage(req.responseChan, errors.New("request timeout"))
			}
		}
	}
}

