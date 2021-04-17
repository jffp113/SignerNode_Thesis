package interconnect

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterHandler(t *testing.T) {
	ic,err := NewInterconnect()
	defer ic.Done()

	assert.Nil(t,err)

	assert.Equal(t,0,len(ic.handlers[SignClientRequest]))
	assert.Equal(t,0,len(ic.handlers[NetworkMessage]))

	ic.RegisterHandler(SignClientRequest, func(content []byte, ctx P2pContext) HandlerResponse {
		return HandlerResponse{}
	})

	ic.RegisterHandler(SignClientRequest, func(content []byte, ctx P2pContext) HandlerResponse {
		return HandlerResponse{}
	})

	ic.RegisterHandler(SignClientRequest, func(content []byte, ctx P2pContext) HandlerResponse {
		return HandlerResponse{}
	})

	assert.Equal(t,3,len(ic.handlers[SignClientRequest]))
	assert.Equal(t,0,len(ic.handlers[NetworkMessage]))
}

func TestEmitEvent(t *testing.T) {
	ic,err := NewInterconnect()
	defer ic.Done()

	assert.Nil(t,err)

	ic.RegisterHandler(SignClientRequest, func(content []byte, ctx P2pContext) HandlerResponse {
		assert.Equal(t,[]byte("Jorge: "),content)
		return CreateOkMessage(append(content, []byte("Hello")...))
	})

	ic.RegisterHandler(SignClientRequest, func(content []byte, ctx P2pContext) HandlerResponse {
		return CreateOkMessage(append(content, []byte(" World")...))
	})

	resp := ic.EmitEvent(SignClientRequest,[]byte("Jorge: "))

	assert.Equal(t,resp.ResponseData,[]byte("Jorge: Hello World"))
}

func TestEmitEventNoHandlers(t *testing.T) {
	ic,err := NewInterconnect()
	defer ic.Done()

	assert.Nil(t,err)

	resp := ic.EmitEvent(SignClientRequest,[]byte("Jorge: "))

	assert.NotNil(t,resp.Err)
}

func TestConfigInterconnect(t *testing.T){
	ic,err := NewInterconnect(SetNumberOfHandlerWorkers(1))
	defer ic.Done()

	assert.Nil(t,err)
	assert.Equal(t,ic.nWorkers,1)

}

func TestConfigErrorInterconnect(t *testing.T){
	ic,err := NewInterconnect(func(m *interconnect) error {
		return errors.New("random error")
	})
	defer ic.Done()
	assert.NotNil(t,err)
}
