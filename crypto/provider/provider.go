package provider

import (
	"SignerNode/crypto"
	"SignerNode/crypto/tbls"
	"errors"
)

const (
	TBLS256 = "TBLS256"
)

func GetThresholdSignatureScheme(schema string, n int, t int) (crypto.SignerVerifierAggregator, error) {
	//TODO Verifications

	switch schema {
	case TBLS256:
		return tbls.NewTBLS256(n, t), nil
	}

	return nil, errors.New("scheme not found")
}

func GetSignatureScheme(schema string) crypto.SignerVerifier {
	return nil
}

func GetThresholdKeyGenerator(schema string, n int, t int) (crypto.KeyShareGenerator, error) {
	//TODO Verifications

	switch schema {
	case TBLS256:
		return tbls.NewTBLS256KeyGenerator(), nil
	}

	return nil, errors.New("generator not found")
}
