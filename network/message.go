package network

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"google.golang.org/protobuf/proto"
)
import "SignerNode/network/pb"

type networkMessage struct {
	To      peer.ID
	From    peer.ID
	Content []byte
}

func (msg *networkMessage) MarshalBinary() (data []byte, err error) {
	p := pb.NetworkMessage{}
	p.From, err = msg.From.Marshal()

	if err != nil {
		return nil, err
	}

	p.To, err = msg.To.Marshal()

	if err != nil {
		return nil, err
	}

	p.Payload = msg.Content
	return proto.Marshal(&p)
}

func (msg *networkMessage) UnmarshalBinary(data []byte) error {
	p := pb.NetworkMessage{}
	err := proto.Unmarshal(data, &p)

	if err != nil {
		return err
	}

	msg.Content = p.Payload
	err = msg.From.UnmarshalBinary(p.From)

	if err != nil {
		return err
	}

	_ = msg.To.UnmarshalBinary(p.To)

	return err
}
