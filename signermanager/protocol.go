package signermanager

import (
	"SignerNode/network"
	"SignerNode/smartcontractengine"
	"errors"
	"github.com/jffp113/CryptoProviderSDK/crypto"
	"github.com/jffp113/CryptoProviderSDK/keychain"
	"go.uber.org/atomic"
	"sync"
)

const PERMISSIONED = "Permissioned"
const PERMISSIONLESS = "Permissionless"
const BYZANTINE = "Byzantine"

type Protocol interface {
	ProcessMessage(data []byte, ctx processContext)
	Sign(data []byte, ctx signContext)
	InstallShares(data []byte) error
}

func GetProtocol(protocolName string, factory crypto.ContextFactory,
	keychain keychain.KeyChain, scFactory smartcontractengine.SCContextFactory,network network.Network) (Protocol, error) {
	switch protocolName {
	case PERMISSIONED:
		return NewPermissionedProtocol(factory, keychain, scFactory), nil
	case PERMISSIONLESS:
		return NewPermissionlessProtocol(factory,scFactory,network),nil
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
	responseChan chan<- ManagerResponse
	//Signature shares from every signernode
	shares                [][]byte
	sharesChan            chan []byte
	aggregatingInProgress atomic.Bool
	insertInSharesChan    bool
	t, n                  int
	scheme                string
	uuid                  string
	digest                []byte
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
