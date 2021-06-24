package protocol

import (
	"github.com/golang/protobuf/proto"
	ic "github.com/jffp113/SignerNode_Thesis/interconnect"
	"github.com/jffp113/SignerNode_Thesis/signermanager/pb"
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

type Status int

//Protocol

const (
	Ok Status = iota
	Error
	InvalidTransaction
)

type ManagerResponse struct {
	ResponseStatus Status
	ResponseData   []byte
	Err            error
}

//SignerManager

func CreateInvalidMessageVerifyResponse() ic.HandlerResponse {
	logger.Debug("Creating invalid verify response")
	msg := pb.ClientVerifyResponse{Status: pb.ClientVerifyResponse_INVALID}
	b, _ := proto.Marshal(&msg)
	return ic.CreateOkMessage(b)
}

func CreateValidMessageVerifyMessages() ic.HandlerResponse {
	logger.Debug("Creating valid verify response")
	msg := pb.ClientVerifyResponse{Status: pb.ClientVerifyResponse_OK}
	b, _ := proto.Marshal(&msg)
	return ic.CreateOkMessage(b)
}

func CreateValidMembershipResponse(addrs []peer.AddrInfo) ic.HandlerResponse {
	var peers []*pb.MembershipResponsePeer

	for _, v := range addrs {
		peers = append(peers, &pb.MembershipResponsePeer{
			Id:   v.ID.String(),
			Addr: peerAddressesToStringSlice(v.Addrs),
		})
	}

	resp := pb.MembershipResponse{Status: pb.MembershipResponse_OK, Peers: peers}
	b, _ := proto.Marshal(&resp)
	return ic.CreateOkMessage(b)
}

func peerAddressesToStringSlice(addr []ma.Multiaddr) []string {
	var result []string

	for _, v := range addr {
		result = append(result, v.String())
	}

	return result
}

func createSignResponse(UUID string, signature []byte) ([]byte, error) {
	resp := pb.SignResponse{
		UUID:      UUID,
		Signature: signature,
	}

	respData, err := proto.Marshal(&resp)

	if err != nil {
		logger.Error(err)
		return []byte{}, err
	}

	return createProtocolMessage(respData, pb.ProtocolMessage_SIGN_RESPONSE)
}
