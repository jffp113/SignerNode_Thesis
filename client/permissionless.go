package client

import (
	"github.com/jffp113/SignerNode_Thesis/signermanager/pb"
	"bytes"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/jffp113/CryptoProviderSDK/crypto"
	"github.com/jffp113/CryptoProviderSDK/example/handlers/tbls"
	"github.com/jffp113/CryptoProviderSDK/example/handlers/trsa"
	"net/http"
	"time"
)


func (s *signerNodeClient) permissionlessProtocol(toSignBytes []byte, smartcontract string) (pb.ClientSignResponse, error) {
	panic("implement me")
}

type InstallMode int
const (
	GenerateAndInstallToKnow InstallMode = iota
	GenerateAndInstallToMembership
)

func (s *signerNodeClient) InstallShare(mode InstallMode, t,n int, scheme string, validUntil time.Time, isOneTimeKey bool) error{
	switch mode {
	case GenerateAndInstallToKnow:
		return s.generateAndInstallToKnow(t,n,scheme,validUntil,isOneTimeKey)
	}
	//TODO all known
	//TODO get membership and install
	return errors.New("mode does not exist")
}

func (s *signerNodeClient) generateAndInstallToKnow(t, n int, scheme string, validUntil time.Time, isOneTimeKey bool) error {
	if len(s.signerNodeAddress) < n {
		return errors.New("not enough signer nodes")
	}

	kg := getKeyGen(scheme)

	pub,privList := kg.Gen(n,t)

	membership := getLocalSubsetMembership(s.signerNodeAddress,n)

	for i,k := range privList {
		err := s.installShare(k,pub,validUntil,isOneTimeKey,membership[i])
		if err != nil {
			return err
		}
	}

	return nil
}


func (s *signerNodeClient) installShare(priv crypto.PrivateKey,pub crypto.PublicKey,validUntil time.Time,
	isOneTimeKey bool, address string) error{

	privBytes,_ := priv.MarshalBinary()
	pubBytes,_ := pub.MarshalBinary()

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


func getKeyGen(scheme string) crypto.KeyShareGenerator {
	switch scheme {
	case "TBLS256":
		return tbls.NewTBLS256KeyGenerator()
	case "TRSA1024":
		return trsa.NewTRSAKeyGenerator()
	default:
		return nil
	}
}
