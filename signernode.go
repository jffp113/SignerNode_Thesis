package main

import (
	"SignerNode/api"
	"SignerNode/signermanager"
	"fmt"
	"github.com/ipfs/go-log"
	"github.com/jessevdk/go-flags"
	"os"
)

type Opts struct {
	//Verbose []bool `short:"v" long:"verbose" description:"Increase verbosity"`
	ApiPort       int    `short:"p" long:"port" description:"API Port" default:"8080"`
	SignerURI     string `short:"s" long:"signer" description:"Signer URI" default:"tcp://127.0.0.1:9000"`
	BootstrapNode string `short:"b" long:"bootstrap" description:"Boostrap Node to find other signer nodes"`
	KeyPath       string `short:"k" long:"keys" description:"Path for the private key and public key" default:"./resources/"`
	Protocol      string `short:"t" long:"protocol" description:"API Port" default:"Permissioned"`
}

func main() {
	_ = log.SetLogLevel("network", "debug")
	_ = log.SetLogLevel("protocol", "debug")
	_ = log.SetLogLevel("pubsub", "debug")
	_ = log.SetLogLevel("connmgr", "debug")
	_ = log.SetLogLevel("crypto_client", "debug")
	_ = log.SetLogLevel("dht", "warn")

	var opts Opts

	parser := flags.NewParser(&opts, flags.Default)
	remaining, err := parser.Parse()
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			fmt.Printf("Failed to parse args: %v\n", err)
			os.Exit(2)
		}
	}

	if len(remaining) > 0 {
		fmt.Printf("Error: Unrecognized arguments passed: %v\n", remaining)
		os.Exit(2)
	}

	sm := signermanager.NewSignerManager()
	sm.SetBootStrapNode(opts.BootstrapNode)
	sm.SetKeyPath(opts.KeyPath)
	sm.SetProtocol(opts.Protocol)
	sm.SetSignerURI(opts.SignerURI)

	err = sm.Init()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	api.Init(opts.ApiPort, sm.Sign)
}
