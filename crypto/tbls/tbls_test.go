package tbls

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTBLS(test *testing.T) {
	n := 10
	t := n/2 + 1

	for i := t; i <= n; i++ {
		tblsSuccessSignature(n, i, test)
	}
}

func tblsSuccessSignature(n, t int, test *testing.T) {
	var err error
	msg := []byte("Test TBLS")

	keygen := NewTBLS256KeyGenerator()
	tbls := NewTBLS256(n, t)
	pub, shares := keygen.Gen(n, t)

	sigShares := make([][]byte, 0)
	for _, x := range shares {
		s, err := tbls.Sign(msg, x)
		require.Nil(test, err)
		sigShares = append(sigShares, s)
	}

	sig, err := tbls.Aggregate(sigShares, msg, pub)

	require.Nil(test, err)

	err = tbls.Verify(sig, msg, pub)
	require.Nil(test, err)
}

func TestTBLSNotEnoughShares(test *testing.T) {
	var err error
	msg := []byte("Test TBLS")

	n := 10
	t := n/2 + 1

	keygen := NewTBLS256KeyGenerator()
	tbls := NewTBLS256(n, t)
	pub, shares := keygen.Gen(n, t)

	sigShares := make([][]byte, 0)
	for _, x := range shares[0 : t-1] {
		s, err := tbls.Sign(msg, x)
		require.Nil(test, err)
		sigShares = append(sigShares, s)
	}

	_, err = tbls.Aggregate(sigShares, msg, pub)

	require.NotNil(test, err)
}
