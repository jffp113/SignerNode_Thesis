package api

import (
	"fmt"
	"github.com/ipfs/go-log"
	"io/ioutil"
	"net/http"
)

var logger = log.Logger("api")

func signHandler(w http.ResponseWriter, r *http.Request, f func(data []byte) <-chan []byte) {

	switch r.Method {
	case http.MethodPost:
		b, _ := ioutil.ReadAll(r.Body)

		fmt.Println(b)
		respChan := f(b)
		respBytes := <-respChan

		w.Write(respBytes)
	default:
		w.WriteHeader(405)
	}
}

func Init(port int, f func(data []byte) <-chan []byte) {
	http.HandleFunc("/sign", func(writer http.ResponseWriter, request *http.Request) {
		signHandler(writer, request, f)
	})
	logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}
