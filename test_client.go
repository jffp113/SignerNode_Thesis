package main

import (
	"SignerNode/signermanager/pb"
	"bytes"
	"fmt"
	"github.com/golang/protobuf/proto"
	"net/http"
	uuid "github.com/satori/go.uuid"
)

func main() {

	msg:=pb.ClientMessage{
		UUID:          fmt.Sprint(uuid.NewV4()),
		SmartContract: "none",
		T:             3,
		N:             5,
		Scheme:        "TBLS256",
		Content:       []byte("Hello"),
	}

	b,err := proto.Marshal(&msg)

	reader := bytes.NewReader(b)

	resp, err := http.Post("http://localhost:8081/sign","application/protobuf",reader)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp)


}
