package client

const (
	PermissionedProtocol   = "Permissioned"
	PermissionlessProtocol = "Permissionless"
)

type PermissionedConfig func(m *permisisonedClient) error

func SetSignerNodeAddresses(addresses ...string) PermissionedConfig {
	return func(m *permisisonedClient) error {
		m.signerNodeAddress = addresses
		return nil
	}
}