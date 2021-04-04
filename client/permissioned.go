package client

import "github.com/jffp113/SignerNode_Thesis/signermanager/pb"

func (s *signerNodeClient) permissionedProtocol(toSignBytes []byte, smartcontract string) (pb.ClientSignResponse, error) {
	return s.sign(toSignBytes,smartcontract)
}
