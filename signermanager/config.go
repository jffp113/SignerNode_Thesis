package signermanager

type Config func(m *signermanager) error

func SetBootStrapNode(bootstrap string) Config {
	return func(m *signermanager) error {
		m.bootstrapNode = bootstrap
		return nil
	}
}

func SetKeyPath(keyPath string) Config {
	return func(m *signermanager) error {
		m.keyPath = keyPath
		return nil
	}
}

func SetProtocol(protocol string) Config {
	return func(m *signermanager) error {
		m.protocolName = protocol
		return nil
	}
}

func SetSignerURI(uri string) Config {
	return func(m *signermanager) error {
		m.signerURI = uri
		return nil
	}
}

func SetScURI(uri string) Config {
	return func(m *signermanager) error {
		m.scURI = uri
		return nil
	}
}

func SetPeerPort(port int) Config {
	return func(m *signermanager) error {
		m.peerPort = port
		return nil
	}
}

func SetPeerAddress(addr string) Config {
	return func(m *signermanager) error {
		m.peerAddress = addr
		return nil
	}
}

func SetBroadcastAnswer(broadcast bool) Config {
	return func(m *signermanager) error {
		m.broadcastAnswer = broadcast
		return nil
	}
}
