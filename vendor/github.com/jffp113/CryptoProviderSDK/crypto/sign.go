package crypto

import (
	"encoding"
	"io"
)

// PublicKey represents a public key using an unspecified algorithm.
type PublicKey interface {
	encoding.BinaryMarshaler
}

// PrivateKey represents a private key using an unspecified algorithm.
type PrivateKey interface {
	encoding.BinaryMarshaler
}

type PrivateKeyList []PrivateKey

func (priv PrivateKeyList) MarshalBinary() (data [][]byte, err error) {
	privSlice := make([][]byte, len(priv))

	for i, v := range priv {
		bytes, err := v.MarshalBinary()

		if err != nil {
			return nil, err
		}
		privSlice[i] = bytes
	}

	return privSlice, nil
}

type KeyGenerator interface {
	Gen() (PublicKey, PrivateKey)
}

type KeyShareGenerator interface {
	Gen(n int, t int) (PublicKey, PrivateKeyList)
}

type THSignerHandler interface {
	KeyShareGenerator
	SignerVerifierAggregator
	SchemeName() string
	UnmarshalPublic(data []byte) PublicKey
	UnmarshalPrivate(data []byte) PrivateKey
}

type SignerVerifier interface {
	Signer
	Verifier
}

type SignerVerifierAggregator interface {
	Signer
	Verifier
	Aggregator
}

type Signer interface {
	Sign(digest []byte, key PrivateKey) (signature []byte, err error)
}

type Verifier interface {
	Verify(signature []byte, msg []byte, key PublicKey) error
}

type Aggregator interface {
	Aggregate(share [][]byte, digest []byte, key PublicKey, t, n int) (signature []byte, err error)
}

type ContextFactory interface {
	GetSignerVerifierAggregator(cryptoId string) (SignerVerifierAggregator, io.Closer)
	GetKeyGenerator(cryptoId string) (KeyShareGenerator, io.Closer)
}
