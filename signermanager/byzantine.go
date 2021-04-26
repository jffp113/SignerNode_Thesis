package signermanager

import (
	"github.com/golang/protobuf/proto"
	ic "github.com/jffp113/SignerNode_Thesis/interconnect"
	"github.com/jffp113/SignerNode_Thesis/signermanager/pb"
)

type byzantineProtocol struct {
}

func (p *byzantineProtocol) Register(interconnect ic.Interconnect) error {
	interconnect.RegisterHandler(ic.SignClientRequest,p.Sign)
	interconnect.RegisterHandler(ic.InstallClientRequest,p.InstallShares)
	interconnect.RegisterHandler(ic.NetworkMessage,p.processMessage)
	return nil
}

func (p *byzantineProtocol) InstallShares(data []byte,ctx ic.P2pContext) ic.HandlerResponse {
	return ic.CreateOkMessage([]byte{})
}

func (p *byzantineProtocol) processMessage(data []byte, ctx ic.P2pContext) ic.HandlerResponse {
	logger.Debug("Received sign Request, processing.")

	req := pb.ProtocolMessage{}
	proto.Unmarshal(data, &req)

	switch req.Type {
	case pb.ProtocolMessage_SIGN_REQUEST:
		p.processMessageSignRequest(&req, ctx)
	case pb.ProtocolMessage_SIGN_RESPONSE:
		p.processMessageSignResponse(&req, ctx)
	}
	return ic.CreateOkMessage(data)
}

func (p *byzantineProtocol) processMessageSignResponse(req *pb.ProtocolMessage, ctx ic.P2pContext) {
	logger.Debug("Byzantine Do nothing")
}

func (p *byzantineProtocol) processMessageSignRequest(req *pb.ProtocolMessage, ctx ic.P2pContext) {
	logger.Debug("Received sign(Byzantine) Request")
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

	ctx.Broadcast(data)
}

func (p *byzantineProtocol) Sign(data []byte, ctx ic.P2pContext) ic.HandlerResponse {
	//Do nothing
	return ic.CreateOkMessage([]byte{})
}

func NewByzantineProtocol() Protocol {
	return &byzantineProtocol{}
}
