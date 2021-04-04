package client

const (
	PermissionedProtocol   = "Permissioned"
	PermissionlessProtocol = "Permissionless"
)

type Config func(m *signerNodeClient) error

func SetProtocol(protocol string) Config {
	return func(m *signerNodeClient) error {
		m.protocol = protocol
		return nil
	}
}

func SetSignerNodeAddresses(addresses ...string) Config {
	return func(m *signerNodeClient) error {
		m.signerNodeAddress = addresses
		return nil
	}
}
