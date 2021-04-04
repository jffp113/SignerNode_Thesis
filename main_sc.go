package main

import (
	"SignerNode/smartcontractengine"
	"fmt"
	"github.com/ipfs/go-log"
	"time"
)

func main() {
	_ = log.SetLogLevel("smartcontract_engine", "debug")
	s, err := smartcontractengine.NewSmartContractClientFactory("tcp://127.0.0.1:9000")

	if err != nil {
		fmt.Println(err)
	}

	time.Sleep(2 * time.Second)

	c, _ := s.GetContext("intkey")

	r := c.InvokeSmartContract([]byte("hello"))

	fmt.Println(r)

	r2 := c.InvokeSmartContract([]byte("hello"))

	fmt.Println(r2)

	r3 := c.InvokeSmartContract([]byte("hello"))

	fmt.Println(r3)

	select {}
}
