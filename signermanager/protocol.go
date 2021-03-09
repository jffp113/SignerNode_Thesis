package signermanager

import (
	"SignerNode/smartcontractengine"
	"errors"
	"github.com/jffp113/CryptoProviderSDK/crypto"
	"github.com/jffp113/CryptoProviderSDK/keychain"
)

const PERMISSIONED = "Permissioned"

type Protocol interface {
	ProcessMessage(data []byte, ctx processContext)
	Sign(data []byte, ctx signContext)
}

func GetProtocol(protocolName string, factory crypto.ContextFactory,
	keychain keychain.KeyChain, scFactory smartcontractengine.SCContextFactory) (Protocol, error) {
	switch protocolName {
	case PERMISSIONED:
		return NewPermissionedProtocol(factory, keychain,scFactory), nil
	default:
		return nil, errors.New("protocol does not exist")
	}

}
