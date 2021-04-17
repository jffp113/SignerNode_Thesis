package interconnect

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateOkMessage(t *testing.T) {
	msg := CreateOkMessage([]byte("hello"))
	assert.Equal(t,Ok,msg.ResponseStatus)
	assert.Equal(t,[]byte("hello"),msg.ResponseData)
}

func TestCreateErrorMessage(t *testing.T) {
	err := errors.New("severe error")
	msg := CreateErrorMessage(err)
	assert.Equal(t,Error,msg.ResponseStatus)
	assert.Equal(t,err,msg.Err)
}

func TestCreateInvalidTransactionMessage(t *testing.T) {
	msg := CreateInvalidTransactionMessage()
	assert.Equal(t,InvalidTransaction,msg.ResponseStatus)
}

func TestSendOkMessage(t *testing.T) {
	ch := make(chan HandlerResponse,1)
	SendOkMessage(ch,[]byte("hello"))
	resp := <-ch

	assert.Equal(t,Ok,resp.ResponseStatus)
	assert.Equal(t,[]byte("hello"),resp.ResponseData)
}

func TestSendErrorMessage(t *testing.T) {
	ch := make(chan HandlerResponse,1)
	err := errors.New("severe error")
	SendErrorMessage(ch,err)
	resp := <-ch
	assert.Equal(t,Error,resp.ResponseStatus)
	assert.Equal(t,err,resp.Err)
}

func TestSendInvalidTransactionMessage(t *testing.T) {
	ch := make(chan HandlerResponse,1)
	SendInvalidTransactionMessage(ch)
	resp := <-ch
	assert.Equal(t,InvalidTransaction,resp.ResponseStatus)
}