package client

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/jffp113/CryptoProviderSDK/crypto"
	"github.com/jffp113/SignerNode_Thesis/signermanager/pb"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
)



func signPermissioned(toSignBytes []byte, smartcontract string, signerAddress string) (pb.ClientSignResponse, error) {
	return signPermissionless(toSignBytes,smartcontract,signerAddress,"")
}

func signPermissionless(toSignBytes []byte, smartcontract string, signerAddress string,keyId string) (pb.ClientSignResponse, error) {
	uuid := fmt.Sprint(uuid.NewV4())
	msg := pb.ClientSignMessage{
		UUID:                 fmt.Sprint(uuid),
		Content:              toSignBytes,
		SmartContractAddress: smartcontract,
		KeyId: keyId,
	}

	b, err := proto.Marshal(&msg)

	reader := bytes.NewReader(b)

	completeAddress := fmt.Sprintf("http://%v/sign", signerAddress)
	resp, err := http.Post(completeAddress, "application/protobuf", reader)

	if err != nil {
		return pb.ClientSignResponse{}, err
	}

	if resp.StatusCode == 500{
		return pb.ClientSignResponse{}, errors.New("signer node did not sign")
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return pb.ClientSignResponse{}, err
	}

	respMsg := pb.ClientSignResponse{}

	err = proto.Unmarshal(body, &respMsg)

	fmt.Println(resp)

	return respMsg, err
}

func verifySignature(digest []byte, sig []byte, scheme string, pubKey crypto.PublicKey, signerAddress string) error {

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

	completeAddress := fmt.Sprintf("http://%v/verify", signerAddress)
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