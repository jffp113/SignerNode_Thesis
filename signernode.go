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
	Verbose 	  []bool `short:"v" long:"verbose" description:"Increase verbosity"`
	ApiPort       int    `short:"p" long:"port" description:"API Port" default:"8080"`
	SignerURI     string `short:"s" long:"signer" description:"Signer URI" default:"tcp://eth0:9000"`
	ScURI     	  string `short:"c" long:"smartcontract" description:"SmartContract URI" default:"tcp://eth0:4004"`
	BootstrapNode string `short:"b" long:"bootstrap" description:"Boostrap Node to find other signer nodes"`
	KeyPath       string `short:"k" long:"keys" description:"Path for the private key and public key" default:"./resources/"`
	Protocol      string `short:"t" long:"protocol" description:"API Port" default:"Permissioned"`
}

func main() {
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

	
	//Set log level 

	switch len(opts.Verbose) {
		case 2:
			log.SetAllLoggers(log.LevelDebug)
		case 1:
			log.SetAllLoggers(log.LevelInfo)
		default:
			log.SetAllLoggers(log.LevelWarn)
	}
	_ = log.SetLogLevel("dht", "warn")


	//Set configs in signermanager
	sm := signermanager.NewSignerManager()
	sm.SetBootStrapNode(opts.BootstrapNode)
	sm.SetKeyPath(opts.KeyPath)
	sm.SetProtocol(opts.Protocol)
	sm.SetSignerURI(opts.SignerURI)
	sm.SetScURI(opts.ScURI)

	//Initiate signermanager
	err = sm.Init()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	api.Init(opts.ApiPort, sm.Sign)
}
