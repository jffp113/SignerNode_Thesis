package main

import (
	"SignerNode/client"
	"fmt"
	"github.com/jffp113/CryptoProviderSDK/keychain"
	"io/ioutil"
	"net/http"
	"time"
)

func sign(){

	c, _ := client.NewClient(client.SetProtocol("Permissioned"),
							client.SetSignerNodeAddresses("localhost:8080"))



	start := time.Now()
	resp, err := c.SendSignRequest([]byte("Hello"),"intkey")

	fmt.Println(resp)
	t := time.Now()
	elapsed := t.Sub(start)

	fmt.Printf("Elapsed: %v\n",elapsed)

	//Verify sig

	kc := keychain.NewKeyChain("./resources/keys/1/")

	pubKey,err := kc.LoadPublicKey("TBLS256_5_3")

	if err != nil {
		fmt.Print(err)
		return
	}

	err = c.VerifySignature([]byte("Hello"),resp.Signature,resp.Scheme,pubKey)

	if err != nil {
		fmt.Println("Invalid signature ", err)
	}
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
