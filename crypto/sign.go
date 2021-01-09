package crypto

// PublicKey represents a public key using an unspecified algorithm.
type PublicKey interface{}

// PrivateKey represents a private key using an unspecified algorithm.
type PrivateKey interface{}

type KeyGenerator interface {
	Gen() (PublicKey, PrivateKey)
}

type KeyShareGenerator interface {
	Gen(n int, t int) (PublicKey, []PrivateKey)
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
	Aggregate(share [][]byte, digest []byte, key PublicKey) (signature []byte, err error)
}
