package signermanager

import (
	"github.com/jffp113/SignerNode_Thesis/signermanager/pb"
	"github.com/golang/protobuf/proto"
)

type byzantineProtocol struct {
}

func (p *byzantineProtocol) InstallShares(data []byte) error {
	//errors.New("operation not supported")
	return nil
}

func (p *byzantineProtocol) ProcessMessage(data []byte, ctx processContext) {
	logger.Debug("Received Sign Request, processing.")

	req := pb.ProtocolMessage{}
	proto.Unmarshal(data, &req)

	switch req.Type {
	case pb.ProtocolMessage_SIGN_REQUEST:
		p.processMessageSignRequest(&req, ctx)
	case pb.ProtocolMessage_SIGN_RESPONSE:
		p.processMessageSignResponse(&req, ctx)
	}

}

func (p *byzantineProtocol) processMessageSignResponse(req *pb.ProtocolMessage, ctx processContext) {
	logger.Debug("Byzantine Do nothing")

}

func (p *byzantineProtocol) processMessageSignRequest(req *pb.ProtocolMessage, ctx processContext) {
	logger.Debug("Received Sign(Byzantine) Request")
	reqSign := pb.ClientSignMessage{}
	err := proto.Unmarshal(req.Content, &reqSign)

	if err != nil {
		logger.Error(err)
		return
	}

	resp := pb.SignResponse{
		UUID:      reqSign.UUID,
		Signature: []byte("I am a byzantine share"),
	}

	data, err := proto.Marshal(&resp)

	if err != nil {
		logger.Error(err)
		return
	}

	data, err = createProtocolMessage(data, pb.ProtocolMessage_SIGN_RESPONSE)

	ctx.broadcast(data)
}

func (p *byzantineProtocol) Sign(data []byte, ctx signContext) {
	//Do nothing
}

func NewByzantineProtocol() Protocol {

	return &byzantineProtocol{}
}
