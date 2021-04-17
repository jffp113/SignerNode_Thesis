package api

import (
	ic "github.com/jffp113/SignerNode_Thesis/interconnect"
	"fmt"
	"github.com/ipfs/go-log"
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

		fmt.Println(b)
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
			w.Write([]byte(resp.Err.Error()))
		}
	default:
		w.WriteHeader(405)
	}
}

func Init(port int, emitEvent func(t ic.HandlerType, content []byte) ic.HandlerResponse) {
	http.HandleFunc("/sign", func(writer http.ResponseWriter, request *http.Request) {
		httpPostHandler(writer, request, func(data []byte) ic.HandlerResponse {
			return emitEvent(ic.SignClientRequest,data)
		})
	})
	http.HandleFunc("/verify", func(writer http.ResponseWriter, request *http.Request) {
		httpPostHandler(writer, request, func(data []byte) ic.HandlerResponse {
			return emitEvent(ic.VerifyClientRequest,data)
		})
	})
	http.HandleFunc("/install", func(writer http.ResponseWriter, request *http.Request) {
		httpPostHandler(writer, request, func(data []byte) ic.HandlerResponse {
			return emitEvent(ic.InstallClientRequest,data)
		})
	})
	http.HandleFunc("/membership", func(writer http.ResponseWriter, request *http.Request) {
		httpGetHandler(writer, request, func(data []byte) ic.HandlerResponse {
			return emitEvent(ic.MembershipClientRequest,data)
		})
	})
	logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}
