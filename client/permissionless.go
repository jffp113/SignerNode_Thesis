package client

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/jffp113/CryptoProviderSDK/crypto"
	"github.com/jffp113/SignerNode_Thesis/signermanager/pb"
	"go.uber.org/atomic"
	"net/http"
	"time"
)

//Key is a struct that aggregates everything
//related to a installed key
//The user of this API should keep it to ask
//signer nodes to sign a transaction
type Key struct {
	T, N            int
	Scheme          string
	ValidUntil      time.Time
	IsOneTimeKey    bool
	used            atomic.Bool
	PubKey          crypto.PublicKey
	PrivKeys        crypto.PrivateKeyList
	GroupMembership []string
}

func (k *Key) Validate() error {
	if k.T > k.N {
		return errors.New("T can not be higher than N")
	} else if len(k.GroupMembership) != k.N {
		return errors.New(fmt.Sprintf("group membership should have size %v", k.N))
	} else if len(k.PrivKeys) != k.N {
		return errors.New(fmt.Sprintf("priv keys should contain %v keys", k.N))
	}
	return nil
}

func (k *Key) GetKeyId() (string, error) {
	b, err := k.PubKey.MarshalBinary()

	if err != nil {
		return "", err
	}

	h := sha256.New()
	h.Write(b)

	return hex.EncodeToString(h.Sum(nil)), nil
}

type PermissionlessClient interface {
	SendSignRequest(toSignBytes []byte, smartcontract string, key *Key) (pb.ClientSignResponse, error)
	VerifySignature(digest []byte, sig []byte, scheme string, key *Key) error
	InstallShare(key *Key) error
}

type permissionlessClient struct{}

func (s permissionlessClient) SendSignRequest(toSignBytes []byte, smartcontract string, key *Key) (pb.ClientSignResponse, error) {
	keyId, err := key.GetKeyId()
	if err != nil {
		return pb.ClientSignResponse{}, err
	}

	return signPermissionless(toSignBytes, smartcontract, GetRandomGroupMember(key.GroupMembership), keyId)
}

func (s permissionlessClient) VerifySignature(digest []byte, sig []byte, scheme string, key *Key) error {
	return verifySignature(digest, sig, scheme, key.PubKey, GetRandomGroupMember(key.GroupMembership))
}

func NewPermissionlessClient() PermissionlessClient {
	return permissionlessClient{}
}

func (s permissionlessClient) InstallShare(key *Key) error {
	if err := key.Validate(); err != nil {
		return err
	}

	//membership := GetSubsetMembership(s.signerNodeAddress,N)
	for i, privKeyShare := range key.PrivKeys {
		err := s.installShare(privKeyShare, key.PubKey, key.ValidUntil,
			key.IsOneTimeKey, key.GroupMembership[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (s permissionlessClient) installShare(priv crypto.PrivateKey, pub crypto.PublicKey, validUntil time.Time,
	isOneTimeKey bool, address string) error {

	privBytes, _ := priv.MarshalBinary()
	pubBytes, _ := pub.MarshalBinary()

	msg := pb.ClientInstallShareRequest{
		PublicKey:    pubBytes,
		PrivateKey:   privBytes,
		ValidUntil:   validUntil.Unix(),
		IsOneTimeKey: isOneTimeKey,
	}

	b, err := proto.Marshal(&msg)

	if err != nil {
		return err
	}

	reader := bytes.NewReader(b)

	completeAddress := fmt.Sprintf("http://%v/install", address)
	_, err = http.Post(completeAddress, "application/protobuf", reader)

	return err
}
