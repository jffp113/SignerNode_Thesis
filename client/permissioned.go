package client

import (
	"github.com/jffp113/CryptoProviderSDK/crypto"
	"github.com/jffp113/SignerNode_Thesis/signermanager/pb"
	"math/rand"
	"time"
)

//A PermissionedClient gives support to communicate with
//our signer nodes.
type PermissionedClient interface {
	SendSignRequest(toSignBytes []byte, smartcontract string) (pb.ClientSignResponse, error)
	VerifySignature(digest []byte, sig []byte, scheme string, pubKey crypto.PublicKey) error
}

type permisisonedClient struct {
	signerNodeAddress []string
}

func (s *permisisonedClient) getSignerAddress() string {
	rand.Seed(time.Now().UnixNano())
	pos := rand.Intn(len(s.signerNodeAddress))
	return s.signerNodeAddress[pos]
}

func (s *permisisonedClient) VerifySignature(digest []byte, sig []byte, scheme string,
	pubKey crypto.PublicKey) error {
	return verifySignature(digest,sig,scheme,pubKey,s.getSignerAddress())
}

func (s *permisisonedClient) SendSignRequest(toSignBytes []byte, smartcontract string) (pb.ClientSignResponse, error) {
	return signPermissioned(toSignBytes,smartcontract,s.getSignerAddress())
}

func NewPermissionedClient(configs ...PermissionedConfig) (PermissionedClient, error) {
	c := permisisonedClient{}

	for _, conf := range configs {
		err := conf(&c)
		if err != nil {
			return nil, err
		}
	}

	return &c, nil
}
