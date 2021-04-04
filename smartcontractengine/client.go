package smartcontractengine

import "io"

type ScResponse struct {
	T, N   int
	Scheme string
	Valid  bool
	Error  bool
}

type SCContext interface {
	InvokeSmartContract(context []byte) ScResponse
}

type SCContextFactory interface {
	GetContext(scAddress string) (SCContext, io.Closer)
}
