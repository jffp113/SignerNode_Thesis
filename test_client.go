package main

import (
	"SignerNode/signermanager/pb"
	"bytes"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/jffp113/CryptoProviderSDK/keychain"
	"io/ioutil"
	"net/http"
	uuid "github.com/satori/go.uuid"
	"time"
)

func sign(){

	uuid := fmt.Sprint(uuid.NewV4())
	msg:=pb.ClientSignMessage{
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
	body,_ :=ioutil.ReadAll(resp.Body)

	respMsg := pb.ClientSignResponse{}

	proto.Unmarshal(body,&respMsg)

	fmt.Println(respMsg)

	verifySig(respMsg.Signature,[]byte("Hello"),respMsg.Scheme)
}

func verifySig(sig []byte,digest []byte,scheme string){
	kc := keychain.NewKeyChain("./resources/keys/1/")

	pubKey,err := kc.LoadPublicKey("TBLS256_5_3")

	if err != nil {
		fmt.Print(err)
		return
	}

	keyBytes,_ := pubKey.MarshalBinary()

	req := pb.ClientVerifyMessage{
		Scheme:   scheme,
		PublicKey: keyBytes,
		Digest:    digest,
		Signature: sig,
	}

	reqBytes,_ := proto.Marshal(&req)

	resp,err := http.Post("http://localhost:8080/verify",
		"application/protobuf",
		bytes.NewReader(reqBytes))

	body,_ :=ioutil.ReadAll(resp.Body)

	fmt.Println(body)
}

func membership(){
	resp,_ := http.Get("http://localhost:8080/membership")

	fmt.Println(resp)
	body,_ :=ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func main() {
	sign()
}
