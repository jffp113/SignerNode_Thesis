package client

import "SignerNode/signermanager/pb"

func (s *signerNodeClient) permissionedProtocol(toSignBytes []byte, smartcontract string) (pb.ClientSignResponse, error) {
	return s.sign(toSignBytes,smartcontract)
}
