package network

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/assert"
	ic "github.com/libp2p/go-libp2p-core/crypto"
	"testing"
)

func TestMarshallUnmarshallMessage(t *testing.T){
	_,pub1, err := ic.GenerateKeyPair(0,2048)
	assert.Nil(t,err)

	_,pub2, err := ic.GenerateKeyPair(0,2048)
	assert.Nil(t,err)

	id1, err := peer.IDFromPublicKey(pub1)
	assert.Nil(t,err)

	id2, err := peer.IDFromPublicKey(pub2)
	assert.Nil(t,err)

	msg := networkMessage{
		To:      id1,
		From:    id2,
		Content: []byte("hello"),
	}

	bytes,err := msg.MarshalBinary()

	assert.Nil(t,err)

	msgUnmarshalled := networkMessage{}
	err = msgUnmarshalled.UnmarshalBinary(bytes)

	assert.Nil(t,err)

	assert.Equal(t,msg.GetTo(),msgUnmarshalled.GetTo())
	assert.Equal(t,msg.GetFrom(),msgUnmarshalled.GetFrom())
	assert.Equal(t,msg.GetData(),msgUnmarshalled.GetData())
}
