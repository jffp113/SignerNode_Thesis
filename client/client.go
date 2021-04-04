package client

import (
	"github.com/jffp113/SignerNode_Thesis/signermanager/pb"
	"bytes"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/jffp113/CryptoProviderSDK/crypto"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

//A Client gives support to communicate with
//our signer nodes.
type Client interface {
	SendSignRequest(toSignBytes []byte, smartcontract string) (pb.ClientSignResponse, error)
	VerifySignature(digest []byte, sig []byte, scheme string, pubKey crypto.PublicKey) error
}

type signerNodeClient struct {
	protocol          string
	signerNodeAddress []string
}


//TODO installed share struct
//TODO add installed shares to a list
//

func (s *signerNodeClient) getSignerAddress() string {
	rand.Seed(time.Now().UnixNano())
	pos := rand.Intn(len(s.signerNodeAddress))
	return s.signerNodeAddress[pos]
}

func (s *signerNodeClient) VerifySignature(digest []byte, sig []byte, scheme string,
	pubKey crypto.PublicKey) error {

	keyBytes, err := pubKey.MarshalBinary()

	if err != nil {
		return err
	}

	req := pb.ClientVerifyMessage{
		Scheme:    scheme,
		PublicKey: keyBytes,
		Digest:    digest,
		Signature: sig,
	}

	reqBytes, _ := proto.Marshal(&req)

	completeAddress := fmt.Sprintf("http://%v/verify", s.getSignerAddress())
	resp, err := http.Post(completeAddress,
		"application/protobuf",
		bytes.NewReader(reqBytes))

	if err != nil {
		return err
	}

	body, _ := ioutil.ReadAll(resp.Body)

	msg := pb.ClientVerifyResponse{}

	proto.Unmarshal(body, &msg)

	if msg.Status != pb.ClientVerifyResponse_OK {
		return errors.New("incorrect signature")
	}

	return nil
}

func (s *signerNodeClient) SendSignRequest(toSignBytes []byte, smartcontract string) (pb.ClientSignResponse, error) {
	switch s.protocol {
	case PermissionedProtocol:
		return s.permissionedProtocol(toSignBytes, smartcontract)
	case PermissionlessProtocol:
		return s.permissionlessProtocol(toSignBytes,smartcontract)
	}

	return pb.ClientSignResponse{}, errors.New(fmt.Sprintf("protocol %s does not exist", s.protocol))
}

func (s *signerNodeClient) sign(toSignBytes []byte, smartcontract string) (pb.ClientSignResponse, error) {
	uuid := fmt.Sprint(uuid.NewV4())
	msg := pb.ClientSignMessage{
		UUID:                 fmt.Sprint(uuid),
		Content:              toSignBytes,
		SmartContractAddress: smartcontract,
	}

	b, err := proto.Marshal(&msg)

	reader := bytes.NewReader(b)

	completeAddress := fmt.Sprintf("http://%v/sign", s.getSignerAddress())
	resp, err := http.Post(completeAddress, "application/protobuf", reader)

	if err != nil {
		return pb.ClientSignResponse{}, err
	}


	if resp.StatusCode == 500{
		return pb.ClientSignResponse{}, errors.New("signer node did not sign")
	}

	body, _ := ioutil.ReadAll(resp.Body)

	respMsg := pb.ClientSignResponse{}

	err = proto.Unmarshal(body, &respMsg)
	return respMsg, err
}

func NewClient(configs ...Config) (Client, error) {
	c := signerNodeClient{}

	for _, conf := range configs {
		err := conf(&c)
		if err != nil {
			return nil, err
		}
	}

	return &c, nil
}
