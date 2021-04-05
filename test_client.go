package main

import (
	"fmt"
	"github.com/jffp113/CryptoProviderSDK/example/handlers/tbls"
	"github.com/jffp113/CryptoProviderSDK/keychain"
	"github.com/jffp113/SignerNode_Thesis/client"
	"io/ioutil"
	"net/http"
	"time"
)

func permissionless(){
	t := 1
	n := 5
	membership := []string{"localhost:8080","localhost:8081","localhost:8082","localhost:8083","localhost:8084"}
	gen := tbls.NewTBLS256KeyGenerator()
	pub, priv := gen.Gen(n,t)

	k := client.Key{
		T:               t,
		N:               n,
		Scheme:          "TBLS256",
		ValidUntil:      time.Now().Add(-24 * time.Hour),
		IsOneTimeKey:    false,
		PubKey:          pub,
		PrivKeys:        priv,
		GroupMembership: membership,
	}

	c := client.NewPermissionlessClient()
	err := c.InstallShare(&k)
	if err != nil {
		fmt.Println("failed ",err)
		return
	}
	fmt.Println(k.GetKeyId())

	//time.Sleep(10*time.Second)
	v,err := c.SendSignRequest([]byte("Hello"), "intkey",&k)

	fmt.Println(v,err)

	fmt.Println(c.VerifySignature([]byte("Hello"),v.Signature,v.Scheme,&k))
}

func sign() {

	c, _ := client.NewPermissionedClient(client.SetSignerNodeAddresses("localhost:8080"))

	start := time.Now()
	resp, err := c.SendSignRequest([]byte("Hello"), "intkey")

	t := time.Now()
	elapsed := t.Sub(start)

	fmt.Printf("Elapsed: %v\n", elapsed)


	if err != nil {
		fmt.Println(err)
		return
	}

	//Verify sig

	kc := keychain.NewKeyChain("./resources/keys/1/")

	pubKey, err := kc.LoadPublicKey("TBLS256_5_3")

	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Println("ok")

	err = c.VerifySignature([]byte("Hello"), resp.Signature, resp.Scheme, pubKey)

	if err != nil {
		fmt.Println("Invalid signature ", err)
	}
}

func membership() {
	resp, _ := http.Get("http://localhost:8080/membership")

	fmt.Println(resp)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func main() {
	permissionless()
}
