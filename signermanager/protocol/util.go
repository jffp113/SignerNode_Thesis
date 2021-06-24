package protocol

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/golang/protobuf/proto"
	"github.com/jffp113/CryptoProviderSDK/crypto"
	"github.com/jffp113/SignerNode_Thesis/signermanager/pb"
)

func createProtocolMessage(msg []byte, messageType pb.ProtocolMessage_Type) ([]byte, error) {
	req := pb.ProtocolMessage{
		Type:    messageType,
		Content: msg,
	}

	b, err := proto.Marshal(&req)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return b, err
}

func signWithShare(content []byte, privShare crypto.PrivateKey, crypto crypto.ContextFactory,
	scheme string, n, t int) ([]byte, error) {

	context, closer := crypto.GetSignerVerifierAggregator(scheme)
	defer closer.Close()
	b, err := context.Sign(content, privShare)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return b, nil
}

func aggregateShares(req *request, pubKey crypto.PublicKey, crypto crypto.ContextFactory) ([]byte, error) {
	context, closer := crypto.GetSignerVerifierAggregator(req.scheme)
	defer closer.Close()

	fullSig, err := context.Aggregate(req.shares, req.digest, pubKey, req.t, req.n)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return fullSig, nil
}

func hash(hash []byte) string {
	h := sha256.New()
	h.Write(hash)

	return hex.EncodeToString(h.Sum(nil))
}
