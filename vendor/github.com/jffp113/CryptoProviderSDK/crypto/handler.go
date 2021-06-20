package crypto

import (
	"github.com/golang/protobuf/proto"
	"github.com/jffp113/CryptoProviderSDK/crypto/pb"
)


type handlerDecorator struct {
	THSignerHandler
}

func (h *handlerDecorator) Handle(msg []byte, msgType int32) ([]byte, int32) {
	var response []byte
	var responseType pb.Type
	switch pb.Type(msgType) {
		case pb.Type_SIGN_REQUEST:
			response,responseType =  h.sign(msg),pb.Type_SIGN_RESPONSE
		case pb.Type_VERIFY_REQUEST:
			response,responseType =  h.verify(msg),pb.Type_VERIFY_RESPONSE
		case pb.Type_AGGREGATE_REQUEST:
			response,responseType =  h.aggregate(msg),pb.Type_AGGREGATE_RESPONSE
		case pb.Type_GENERATE_THS_REQUEST:
			response,responseType =  h.generateTHS(msg),pb.Type_GENERATE_THS_RESPONSE
	}

	return response,int32(responseType)
}

func (h *handlerDecorator) Name() string {
	return h.SchemeName()
}

func (h *handlerDecorator) generateTHS(msg []byte) []byte {
	logger.Debugf("Generating THS keys")
	req := pb.GenerateTHSRequest{}
	err := proto.Unmarshal(msg, &req)

	if err != nil {
		logger.Warn("Ignoring message gen THS message")
		return createGenTHSErrorMsg()
	}

	pub, priv := h.Gen(int(req.N), int(req.T))

	pubBytes, err := pub.MarshalBinary()

	if err != nil {
		logger.Warn("Error marshalling pubkey")
		return createGenTHSErrorMsg()
	}

	privBytes, err := priv.MarshalBinary()

	resp := pb.GenerateTHSResponse{
		Status:      pb.GenerateTHSResponse_OK,
		PublicKey:   pubBytes,
		PrivateKeys: privBytes,
	}

	msgBytes, err := proto.Marshal(&resp)

	if err != nil {
		logger.Warn("Error marshalling answer")
		return createGenTHSErrorMsg()
	}

	logger.Debugf("Finished Generating THS keys")
	return msgBytes
}

func createGenTHSErrorMsg() []byte {
	resp := pb.GenerateTHSResponse{
		Status: pb.GenerateTHSResponse_ERROR,
	}

	msgBytes, _ := proto.Marshal(&resp)

	return msgBytes
}

func (h *handlerDecorator)  aggregate(msg []byte) []byte {
	req := pb.AggregateRequest{}
	err := proto.Unmarshal(msg, &req)
	logger.Debug("Start Aggregating")
	if err != nil {
		logger.Warn("Error unmarshalling request")
		return createAggregateTHSErrorMsg()
	}

	pubKey := h.UnmarshalPublic(req.PubKey)

	sig, err := h.Aggregate(req.Share, req.Digest, pubKey, int(req.T), int(req.N))

	if err != nil {
		logger.Warn("Error generating aggregated signature")
		return createAggregateTHSErrorMsg()
	}

	resp := pb.AggregateResponse{
		Status:    pb.AggregateResponse_OK,
		Signature: sig,
	}

	msgBytes, err := proto.Marshal(&resp)

	if err != nil {
		logger.Warn("Error marshalling answer")
		return createAggregateTHSErrorMsg()
	}

	logger.Debug("End Aggregating")

	return msgBytes
}

func createAggregateTHSErrorMsg() []byte {
	resp := pb.AggregateResponse{
		Status: pb.AggregateResponse_ERROR,
	}

	msgBytes, _ := proto.Marshal(&resp)

	return msgBytes
}

func (h *handlerDecorator) verify(msg []byte) []byte {
	req := pb.VerifyRequest{}
	err := proto.Unmarshal(msg, &req)

	if err != nil {
		logger.Warn("Error marshalling pubkey")
		return createsVerifyTHSErrorMsg()
	}

	pub := h.UnmarshalPublic(req.PubKey)

	err = h.Verify(req.Signature, req.Msg, pub)

	if err != nil {
		logger.Debug("Invalid Signature")
		return createsVerifyTHSErrorMsg()
	}

	resp := pb.VerifyResponse{
		Status: pb.VerifyResponse_OK,
	}

	msgBytes, err := proto.Marshal(&resp)

	if err != nil {
		logger.Warn("Error marshalling response")
		return createsVerifyTHSErrorMsg()
	}

	return msgBytes
}

func createsVerifyTHSErrorMsg() []byte {
	resp := pb.VerifyResponse{
		Status: pb.VerifyResponse_ERROR,
	}

	msgBytes, _ := proto.Marshal(&resp)

	return msgBytes
}

func (h *handlerDecorator)   sign(msg []byte) []byte {
	req := pb.SignRequest{}
	err := proto.Unmarshal(msg, &req)

	logger.Debug("Start Signing")

	if err != nil {
		logger.Warn("Ignoring sign message")
		return createsSignTHSErrorMsg()
	}

	priv := h.UnmarshalPrivate(req.PrivateKeys)

	data, err := h.Sign(req.Digest, priv)

	if err != nil {
		logger.Warn("Error marshalling pubkey")
		return createsSignTHSErrorMsg()
	}

	resp := pb.SignResponse{
		Status:    pb.SignResponse_OK,
		Signature: data,
	}

	msgBytes, err := proto.Marshal(&resp)

	if err != nil {
		logger.Warn("Error marshalling msgBytes")
		return createsSignTHSErrorMsg()
	}

	logger.Debug("End signing")

	return msgBytes
}

func createsSignTHSErrorMsg() []byte {
	resp := pb.SignResponse{
		Status: pb.SignResponse_ERROR,
	}

	msgBytes, _ := proto.Marshal(&resp)

	return msgBytes
}
