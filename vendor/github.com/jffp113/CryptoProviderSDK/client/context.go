package client

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/jffp113/CryptoProviderSDK/crypto"
	"github.com/jffp113/CryptoProviderSDK/crypto/pb"
	"github.com/jffp113/go-util/messaging/routerdealerhandlers/handlerClient"
	"io"
)

type context struct {
	client *cryptoClient
	scheme string
	context handlerClient.Invoker
}

type key []byte

func (key key) MarshalBinary() (data []byte, err error) {
	return key, nil
}

func (c *cryptoClient) GetSignerVerifierAggregator(cryptoId string) (crypto.SignerVerifierAggregator, io.Closer) {
	invoker, closer := c.client.GetContext(cryptoId)

	return &context{c,cryptoId,invoker}, closer
}

func (c *cryptoClient) GetKeyGenerator(cryptoId string) (crypto.KeyShareGenerator, io.Closer) {
	invoker, closer := c.client.GetContext(cryptoId)

	return &context{c,cryptoId,invoker}, closer
}

func (c *context) Sign(digest []byte, key crypto.PrivateKey) (signature []byte, err error) {
	logger.Debugf("Sign Key for %v", c.scheme)

	d, _ := key.MarshalBinary()

	req := pb.SignRequest{
		Scheme:      c.scheme,
		Digest:      digest,
		PrivateKeys: d,
	}

	b,err := proto.Marshal(&req)
	if err != nil {
		return nil, err
	}

	content,_,err := c.context.Invoke(b,int32(pb.Type_SIGN_REQUEST))

	if err != nil {
		return nil, err
	}

	replySign := pb.SignResponse{}
	err = proto.Unmarshal(content, &replySign)

	if err != nil {
		return nil, err
	}

	if replySign.Status != pb.SignResponse_OK {
		return nil, errors.New("error signing")
	}

	return replySign.Signature, nil

}

func (c *context) Verify(signature []byte, msg []byte, key crypto.PublicKey) error {
	logger.Debugf("Verify Request for %v", c.scheme)

	keyBytes, _ := key.MarshalBinary()

	req := pb.VerifyRequest{
		Scheme:    c.scheme,
		Signature: signature,
		Msg:       msg,
		PubKey:    keyBytes,
	}

	b,err := proto.Marshal(&req)
	if err != nil {
		return nil
	}

	content,_,err := c.context.Invoke(b,int32(pb.Type_VERIFY_REQUEST))

	if err != nil {
		return nil
	}

	replySign := pb.VerifyResponse{}
	_ = proto.Unmarshal(content, &replySign)

	if replySign.Status == pb.VerifyResponse_ERROR {
		return errors.New("invalid signature")
	}

	return nil
}

func (c *context) Aggregate(share [][]byte, digest []byte, key crypto.PublicKey, t, n int) (signature []byte, err error) {
	logger.Debugf("Aggregating Request for %v", c.scheme)

	keyBytes, _ := key.MarshalBinary()

	req := pb.AggregateRequest{
		Scheme: c.scheme,
		Share:  share,
		Digest: digest,
		PubKey: keyBytes,
		T:      int32(t),
		N:      int32(n),
	}

	b,err := proto.Marshal(&req)
	if err != nil {
		return nil, nil
	}

	content,_,err := c.context.Invoke(b,int32(pb.Type_AGGREGATE_REQUEST))

	if err != nil {
		return nil, err
	}

	replySign := pb.AggregateResponse{}
	err = proto.Unmarshal(content, &replySign)

	if replySign.Status == pb.AggregateResponse_ERROR {
		return nil, errors.New("error aggregating")
	}

	return replySign.Signature, nil
}

func (c *context) Gen(n int, t int) (crypto.PublicKey, crypto.PrivateKeyList) {
	logger.Debugf("Requesting Key Gen for %v", c.scheme)

	req := pb.GenerateTHSRequest{
		Scheme: c.scheme,
		T:      uint32(t),
		N:      uint32(n),
	}

	b,err := proto.Marshal(&req)
	if err != nil {
		return nil, nil
	}

	content,respType,err := c.context.Invoke(b,int32(pb.Type_GENERATE_THS_REQUEST))

	if err != nil {
		return nil, nil
	}

	if respType != int32(pb.Type_GENERATE_THS_RESPONSE) {
		panic("Wrong message received")
	}

	replyTHS := pb.GenerateTHSResponse{}

	err = proto.Unmarshal(content, &replyTHS)

	if err != nil {
		panic("error unmarshalling msg")
	}

	pubKey := key(replyTHS.PublicKey)
	privKeySlice := make([]crypto.PrivateKey, len(replyTHS.PrivateKeys))

	for i, v := range replyTHS.PrivateKeys {
		privKeySlice[i] = key(v)
	}

	return &pubKey, privKeySlice
}

func (c *context) Close() error {
	//TODO
	return nil
}