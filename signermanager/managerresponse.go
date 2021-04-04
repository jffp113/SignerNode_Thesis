package signermanager

import (
	"SignerNode/signermanager/pb"
	"github.com/golang/protobuf/proto"
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

func sendErrorMessage(c chan<- ManagerResponse, err error) {
	c <- ManagerResponse{
		ResponseStatus: Error,
		Err:            err,
	}
}

func sendInvalidTransactionMessage(c chan<- ManagerResponse) {
	c <- ManagerResponse{
		ResponseStatus: InvalidTransaction,
		Err:            nil,
	}
}

func sendOkMessage(c chan<- ManagerResponse, data []byte) {
	c <- ManagerResponse{
		ResponseStatus: Ok,
		ResponseData:   data,
		Err:            nil,
	}
}

//SignerManager

func createValidMembershipResponse(addrs []peer.AddrInfo, ch chan ManagerResponse) {

	var peers []*pb.MembershipResponsePeer

	for _, v := range addrs {
		peers = append(peers, &pb.MembershipResponsePeer{
			Id:   v.ID.String(),
			Addr: peerAddressesToStringSlice(v.Addrs),
		})
	}
	resp := pb.MembershipResponse{Status: pb.MembershipResponse_OK, Peers: peers}

	b, _ := proto.Marshal(&resp)
	ch <- ManagerResponse{Ok, b, nil}
}

func peerAddressesToStringSlice(addr []ma.Multiaddr) []string {
	var result []string

	for _, v := range addr {
		result = append(result, v.String())
	}

	return result
}

func createInvalidMessageVerifyResponse(ch chan ManagerResponse) {
	logger.Debug("Creating invalid verify response")
	msg := pb.ClientVerifyResponse{Status: pb.ClientVerifyResponse_INVALID}
	b, _ := proto.Marshal(&msg)
	ch <- ManagerResponse{Ok,
		b,
		nil}
}
func createValidMessageVerifyMessages(ch chan ManagerResponse) {
	logger.Debug("Creating valid verify response")
	msg := pb.ClientVerifyResponse{Status: pb.ClientVerifyResponse_OK}
	b, _ := proto.Marshal(&msg)
	ch <- ManagerResponse{Ok,
		b,
		nil}
}
