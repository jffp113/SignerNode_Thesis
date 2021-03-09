package api

import (
	"SignerNode/signermanager"
	"fmt"
	"github.com/ipfs/go-log"
	"io/ioutil"
	"net/http"
)

var logger = log.Logger("api")

func signHandler(w http.ResponseWriter, r *http.Request, f func(data []byte) <-chan signermanager.ManagerResponse) {

	switch r.Method {
	case http.MethodPost:
		b, _ := ioutil.ReadAll(r.Body)

		fmt.Println(b)
		respChan := f(b)
		resp := <-respChan

		switch resp.ResponseStatus {
			case signermanager.Ok:
				w.WriteHeader(202)
			case signermanager.Error:
				w.WriteHeader(500)
		}
	default:
		w.WriteHeader(405)
	}
}

func Init(port int, f func(data []byte) <-chan signermanager.ManagerResponse) {
	http.HandleFunc("/sign", func(writer http.ResponseWriter, request *http.Request) {
		signHandler(writer, request, f)
	})
	logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}
