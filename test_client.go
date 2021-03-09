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
		Content:       []byte("Hello"),
		SmartContractAddress: "intkey",
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
