package interconnect

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetNumberOfHandlerWorkers(t *testing.T) {
	ic := interconnect{
		handlers:  nil,
		eventChan: nil,
		nWorkers:  0,
	}
	assert.Equal(t, 0, ic.nWorkers)
	err := SetNumberOfHandlerWorkers(10)(&ic)
	assert.Nil(t, err)
	assert.Equal(t, 10, ic.nWorkers)
}

func TestSetContext(t *testing.T) {
	ic := interconnect{
		handlers:  nil,
		eventChan: nil,
		nWorkers:  0,
	}
	assert.Nil(t, ic.p2pCtx.Broadcast)
	err := SetContext(P2pContext{
		Broadcast: func(msg []byte) error { return nil },
	})(&ic)
	assert.Nil(t, err)
	assert.NotNil(t, ic.p2pCtx.Broadcast)
}
