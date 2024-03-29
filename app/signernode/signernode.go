package main

import (
	_ "expvar"
	"fmt"
	"github.com/ipfs/go-log/v2"
	"github.com/jessevdk/go-flags"
	"github.com/jffp113/SignerNode_Thesis/api"
	"github.com/jffp113/SignerNode_Thesis/signermanager"
	_ "net/http/pprof"
	"os"
)

type Opts struct {
	Verbose         []bool `short:"v" long:"verbose" description:"Increase verbosity"`
	ApiPort         int    `short:"p" long:"port" description:"API Port" default:"8080"`
	SignerURI       string `short:"s" long:"signer" description:"Signer URI" default:"tcp://eth0:9000"`
	ScURI           string `short:"c" long:"smartcontract" description:"SmartContract URI" default:"tcp://eth0:4004"`
	BootstrapNode   string `short:"b" long:"bootstrap" description:"Boostrap Node to find other signer nodes"`
	KeyPath         string `short:"k" long:"keys" description:"Path for the private key and public key" default:"./resources/"`
	Protocol        string `short:"t" long:"protocol" description:"API Port" default:"Permissioned"`
	PeerPort        int    `long:"peerport" description:"P2P peer port" default:"0"`
	PeerAddress     string `long:"peeraddr" description:"P2P listening address" default:"/ip4/0.0.0.0/tcp/"`
	BroadcastAnswer bool   `long:"broadcastreply" description:"Defines if we a signer node should broadcast a signature share to all nodes"`
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
	_ = log.SetLogLevel("peerqueue", "warn")

	sm := signermanager.NewSignerManager(
		signermanager.SetBootStrapNode(opts.BootstrapNode),
		signermanager.SetKeyPath(opts.KeyPath),
		signermanager.SetProtocol(opts.Protocol),
		signermanager.SetSignerURI(opts.SignerURI),
		signermanager.SetScURI(opts.ScURI),
		signermanager.SetPeerPort(opts.PeerPort),
		signermanager.SetPeerAddress(opts.PeerAddress),
		signermanager.SetBroadcastAnswer(opts.BroadcastAnswer),
	)

	//Initiate signermanager
	err = sm.Init()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	api.Init(opts.ApiPort, sm.EmitEvent)
}

