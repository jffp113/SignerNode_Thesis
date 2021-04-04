package api

import (
	"github.com/jffp113/SignerNode_Thesis/signermanager"
	"fmt"
	"github.com/ipfs/go-log"
	"io/ioutil"
	"net/http"
)

var logger = log.Logger("api")

func httpGetHandler(w http.ResponseWriter, r *http.Request, f func(data []byte) <-chan signermanager.ManagerResponse) {
	httpFuncHandler(w, r, f, http.MethodGet)
}

func httpPostHandler(w http.ResponseWriter, r *http.Request, f func(data []byte) <-chan signermanager.ManagerResponse) {
	httpFuncHandler(w, r, f, http.MethodPost)
}

func httpFuncHandler(w http.ResponseWriter, r *http.Request,
	f func(data []byte) <-chan signermanager.ManagerResponse, method string) {

	switch r.Method {
	case method:
		b, _ := ioutil.ReadAll(r.Body)

		fmt.Println(b)
		respChan := f(b)
		resp := <-respChan

		switch resp.ResponseStatus {
		case signermanager.Ok:
			w.Write(resp.ResponseData)
			//			w.WriteHeader(200)
		case signermanager.Error:
			w.Write([]byte(resp.Err.Error()))
			w.WriteHeader(500)
		}
	default:
		w.WriteHeader(405)
	}
}

func Init(port int, singFunc SignFunc, verifyFunc VerifyFunc, membershipFunc MembershipFunc, installShareFunc InstallShareFunc) {
	http.HandleFunc("/sign", func(writer http.ResponseWriter, request *http.Request) {
		httpPostHandler(writer, request, singFunc)
	})
	http.HandleFunc("/verify", func(writer http.ResponseWriter, request *http.Request) {
		httpPostHandler(writer, request, verifyFunc)
	})
	http.HandleFunc("/install", func(writer http.ResponseWriter, request *http.Request) {
		httpPostHandler(writer, request, verifyFunc)
	})
	http.HandleFunc("/membership", func(writer http.ResponseWriter, request *http.Request) {
		httpGetHandler(writer, request, membershipFunc)
	})
	logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}
