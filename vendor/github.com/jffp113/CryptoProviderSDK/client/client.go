package client

import (
	"github.com/ipfs/go-log"
	"github.com/jffp113/CryptoProviderSDK/crypto"
	"github.com/jffp113/go-util/messaging/routerdealerhandlers/handlerClient"
)

var logger = log.Logger("crypto_client")

type cryptoClient struct {
	client *handlerClient.HandlerClient
}

func NewCryptoFactory(uri string) (crypto.ContextFactory, error) {
	h,err := handlerClient.NewHandlerFactory(uri)

	if err != nil {
		return nil, err
	}

	return &cryptoClient{h}, nil
}

func (c *cryptoClient) Close() error {
	return c.client.Close()
}