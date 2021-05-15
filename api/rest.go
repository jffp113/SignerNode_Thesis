package api

import (
	"fmt"
	"github.com/ipfs/go-log"
	ic "github.com/jffp113/SignerNode_Thesis/interconnect"
	"io/ioutil"
	"net/http"
)

var logger = log.Logger("api")

func httpGetHandler(w http.ResponseWriter, r *http.Request, f func(data []byte) ic.HandlerResponse) {
	httpFuncHandler(w, r, f, http.MethodGet)
}

func httpPostHandler(w http.ResponseWriter, r *http.Request, f func(data []byte) ic.HandlerResponse) {
	httpFuncHandler(w, r, f, http.MethodPost)
}

func httpFuncHandler(w http.ResponseWriter, r *http.Request,
	f func(data []byte) ic.HandlerResponse, method string) {

	switch r.Method {
	case method:
		b, _ := ioutil.ReadAll(r.Body)

		resp := f(b)
		switch resp.ResponseStatus {
		case ic.Ok:
			logger.Debug("Execution Ok")
			w.Write(resp.ResponseData)
		case ic.InvalidTransaction:
			fallthrough
		case ic.Error:
			logger.Debug("Execution Failed, ")
			w.WriteHeader(500)
			if resp.Err != nil {
				w.Write([]byte(resp.Err.Error()))
			}
		}
	default:
		w.WriteHeader(405)
	}
}

func initHttp(port int, emitEvent EmitFunc, handleFuncRegister HandleFunc) {
	handleFuncRegister("/sign", func(writer http.ResponseWriter, request *http.Request) {
		httpPostHandler(writer, request, func(data []byte) ic.HandlerResponse {
			return emitEvent(ic.SignClientRequest, data)
		})
	})

	handleFuncRegister("/verify", func(writer http.ResponseWriter, request *http.Request) {
		httpPostHandler(writer, request, func(data []byte) ic.HandlerResponse {
			return emitEvent(ic.VerifyClientRequest, data)
		})
	})

	handleFuncRegister("/install", func(writer http.ResponseWriter, request *http.Request) {
		httpPostHandler(writer, request, func(data []byte) ic.HandlerResponse {
			return emitEvent(ic.InstallClientRequest, data)
		})
	})

	handleFuncRegister("/membership", func(writer http.ResponseWriter, request *http.Request) {
		httpGetHandler(writer, request, func(data []byte) ic.HandlerResponse {
			return emitEvent(ic.MembershipClientRequest, data)
		})
	})

}

type EmitFunc func(t ic.HandlerType, content []byte) ic.HandlerResponse
type HandleFunc func(pattern string, handler func(http.ResponseWriter, *http.Request))

func Init(port int, emitEvent EmitFunc) {
	initHttp(port, emitEvent, http.HandleFunc)
	logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}
