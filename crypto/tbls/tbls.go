package tbls

import (
	"SignerNode/crypto"
	"errors"
	"go.dedis.ch/kyber/v3/pairing"
	"go.dedis.ch/kyber/v3/pairing/bn256"
	"go.dedis.ch/kyber/v3/share"
	"go.dedis.ch/kyber/v3/sign/bls"
	ths "go.dedis.ch/kyber/v3/sign/tbls"
)

var (
	privateKeyError = errors.New("invalid private key")
)

type tbls struct {
	t, n  int
	suite pairing.Suite
}

func (t *tbls) Sign(digest []byte, key crypto.PrivateKey) ([]byte, error) {
	priv, ok := key.(*share.PriShare)

	if !ok {
		return nil, privateKeyError
	}

	return ths.Sign(t.suite, priv, digest)
}

func (t *tbls) Verify(signature, msg []byte, key crypto.PublicKey) error {
	pub, ok := key.(*share.PubPoly)

	if !ok {
		return privateKeyError
	}

	return bls.Verify(t.suite, pub.Commit(), msg, signature)
}

func (t *tbls) Aggregate(shares [][]byte, digest []byte, key crypto.PublicKey) ([]byte, error) {
	pub, ok := key.(*share.PubPoly)

	if !ok {
		return nil, privateKeyError
	}

	return ths.Recover(t.suite, pub, digest, shares, t.t, t.n)
}

func NewTBLS256(n int, t int) crypto.SignerVerifierAggregator {
	return &tbls{
		t,
		n,
		bn256.NewSuite(),
	}
}

type tblsKeyGenerator struct {
	suite pairing.Suite
}

func (g *tblsKeyGenerator) Gen(n int, t int) (crypto.PublicKey, []crypto.PrivateKey) {
	suite := g.suite
	secret := suite.G1().Scalar().Pick(suite.RandomStream())
	priPoly := share.NewPriPoly(suite.G2(), t, secret, suite.RandomStream())
	pubPoly := priPoly.Commit(suite.G2().Point().Base())

	shares := make([]crypto.PrivateKey, n)
	for i, v := range priPoly.Shares(n) {
		shares[i] = v
	}

	return pubPoly, shares
}

func NewTBLS256KeyGenerator() crypto.KeyShareGenerator {
	return &tblsKeyGenerator{
		bn256.NewSuite(),
	}
}
