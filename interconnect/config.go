package interconnect


type Config func(m *interconnect) error

func SetNumberOfHandlerWorkers(nWorkers int) Config {
	return func(m *interconnect) error {
		m.nWorkers = nWorkers
		return nil
	}
}

func SetContext(ctx P2pContext) Config {
	return func(m *interconnect) error {
		m.p2pCtx = ctx
		return nil
	}
}