package smartcontractengine

import (
	"github.com/golang/protobuf/proto"
	"github.com/ipfs/go-log/v2"
	"github.com/jffp113/SignerNode_Thesis/smartcontractengine/pb"
	"github.com/jffp113/go-util/messaging/routerdealerhandlers/handlerClient"
	"io"
)

var logger = log.Logger("smartcontract_engine")

type scClient struct {
	client *handlerClient.HandlerClient
}

type context struct {
	scID   string
	client *scClient
	invoker handlerClient.Invoker
}

func (c *context) InvokeSmartContract(payload []byte) ScResponse {
	req := pb.SmartContractValidationRequest{
		Payload:              payload,
		SmartContractAddress: c.scID,
	}

	msg, err := proto.Marshal(&req)

	if err != nil {
		logger.Error("Error invoking smartcontract :",err)
	}

	content,respType,err := c.invoker.Invoke(msg,int32(pb.MessageType_SMART_CONTRACT_VALIDATE_REQUEST))

	if err != nil {
		logger.Error("Error invoking smartcontract :",err)
	}

	if pb.MessageType(respType) != pb.MessageType_SMART_CONTRACT_VALIDATE_RESPONSE {
		panic("Wrong message received")
	}

	replyTHS := pb.SmartContractValidationResponse{}

	err = proto.Unmarshal(content, &replyTHS)

	if err != nil {
		panic(err)
	}

	return ScResponse{
		T:      int(replyTHS.T),
		N:      int(replyTHS.N),
		Scheme: replyTHS.SignatureScheme,
		Valid:  replyTHS.Status == pb.SmartContractValidationResponse_OK,
		Error:  replyTHS.Status == pb.SmartContractValidationResponse_INTERNAL_ERROR,
	}
}

func (c *scClient) GetContext(scAddress string) (SCContext, io.Closer) {
	invoker, closer := c.client.GetContext(scAddress)

	r := context{scAddress, c, invoker}


	return &r,closer
}

func NewSmartContractClientFactory(uri string) (SCContextFactory, error) {
	h,err := handlerClient.NewHandlerFactory(uri)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return &scClient{h}, nil
}

func (c *scClient) Close() error {
	return c.client.Close()
}
