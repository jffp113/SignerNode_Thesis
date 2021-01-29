package main

import (
	"SignerNode/signermanager/pb"
	"bytes"
	"fmt"
	"github.com/golang/protobuf/proto"
	"net/http"
	uuid "github.com/satori/go.uuid"
	"time"
)

func main() {

	uuid := fmt.Sprint(uuid.NewV4())
	msg:=pb.ClientMessage{
		UUID:          fmt.Sprint(uuid),
		SmartContract: "none",
		T:             3,
		N:             5,
		Scheme:        "TBLS256",
		Content:       []byte("Hello"),
	}

	b,err := proto.Marshal(&msg)

	reader := bytes.NewReader(b)

	fmt.Println(uuid)
	start := time.Now()
	resp, err := http.Post("http://localhost:8080/sign","application/protobuf",reader)
	t := time.Now()
	elapsed := t.Sub(start)

	fmt.Printf("Elapsed: %v\n",elapsed)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp)


}
