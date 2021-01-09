package provider

import (
	"SignerNode/crypto"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetSignatureSchemeAndSign(test *testing.T) {
	n := 10
	t := 6

	signer, err := GetThresholdSignatureScheme(TBLS256, n, t)
	require.Nil(test, err)

	keygen, err := GetThresholdKeyGenerator(TBLS256, n, t)
	require.Nil(test, err)

	successThresholdSignature(n, t, signer, keygen, test)
}

func successThresholdSignature(n, t int, sva crypto.SignerVerifierAggregator, keygen crypto.KeyShareGenerator, test *testing.T) {
	var err error
	msg := []byte("Test TBLS")

	pub, shares := keygen.Gen(n, t)

	sigShares := make([][]byte, 0)
	for _, x := range shares {
		s, err := sva.Sign(msg, x)
		require.Nil(test, err)
		sigShares = append(sigShares, s)
	}

	sig, err := sva.Aggregate(sigShares, msg, pub)

	require.Nil(test, err)

	err = sva.Verify(sig, msg, pub)
	require.Nil(test, err)
}
