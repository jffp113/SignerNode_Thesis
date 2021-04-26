package api

import (
	"bytes"
	"errors"
	ic "github.com/jffp113/SignerNode_Thesis/interconnect"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

type mockResponseWriter struct {
	status int
	buff bytes.Buffer
}

func (m *mockResponseWriter) Header() http.Header {
	return nil
}

func (m *mockResponseWriter) Write(bytes []byte) (int, error) {
	return m.buff.Write(bytes)
}

func (m *mockResponseWriter) WriteHeader(statusCode int) {
	m.status = statusCode
}

func TestFuncHandlerOk(t *testing.T) {
	data := []byte{123}
	m := mockResponseWriter{status: 200} //Like http status starts at 200
	httpFuncHandler(
		&m,
		&http.Request{
			Method: "POST",
			Body: ioutil.NopCloser(bytes.NewReader(data)),
		},
		func(data []byte) ic.HandlerResponse {
			return ic.CreateOkMessage(data)
		},
		"POST")

	assert.Equal(t,200,m.status)
	assert.Equal(t,data,m.buff.Bytes())
}

func TestFuncHandlerWrongMethod(t *testing.T) {
	data := []byte{123}
	m := mockResponseWriter{status: 200} //Like http status starts at 200
	httpFuncHandler(
		&m,
		&http.Request{
			Method: "GET",
			Body: ioutil.NopCloser(bytes.NewReader(data)),
		},
		func(data []byte) ic.HandlerResponse {
			return ic.CreateOkMessage(data)
		},
		"POST")

	assert.Equal(t,405,m.status)
	assert.Equal(t,[]byte(nil),m.buff.Bytes())
}

func TestFuncHandlerError(t *testing.T) {
	data := []byte{123}
	err := errors.New("mock error")
	m := mockResponseWriter{status: 200} //Like http status starts at 200
	httpFuncHandler(
		&m,
		&http.Request{
			Method: "POST",
			Body: ioutil.NopCloser(bytes.NewReader(data)),
		},
		func(data []byte) ic.HandlerResponse {
			return ic.CreateErrorMessage(err)
		},
		"POST")

	assert.Equal(t,500,m.status)
	assert.Equal(t,err.Error(),string(m.buff.Bytes()))
}

func TestFuncHandlerInvalid(t *testing.T) {
	m := mockResponseWriter{status: 200} //Like http status starts at 200
	data := []byte{123}
	httpFuncHandler(
		&m,
		&http.Request{
			Method: "POST",
			Body: ioutil.NopCloser(bytes.NewReader(data)),
		},
		func(data []byte) ic.HandlerResponse {
			return ic.CreateInvalidTransactionMessage()
		},
		"POST")

	assert.Equal(t,500,m.status)
	assert.Equal(t,[]byte(nil),m.buff.Bytes())
}
