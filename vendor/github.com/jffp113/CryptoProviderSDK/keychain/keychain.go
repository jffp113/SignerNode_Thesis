package keychain

import (
	"encoding"
	"fmt"
	"github.com/jffp113/CryptoProviderSDK/crypto"
	"io/ioutil"
	"os"
)

type KeyChain interface {
	LoadPrivateKey(name string) (crypto.PrivateKey, error)
	LoadPublicKey(name string) (crypto.PublicKey, error)
	StorePublicKey(name string, pub crypto.PublicKey) error
	StorePrivateKey(name string, priv crypto.PrivateKey) error
}

const PrivateKeyPrefix = "priv_%v"
const PublicKeyPrefix = "pub_%v"

//This keychain can store and read public/private
//keys for the client only in binary
type keychain struct {
	directory string
	cache     map[string][]byte
}

type key []byte

func (k key) MarshalBinary() (data []byte, err error) {
	return k, nil
}

func NewKeyChain(directory string) KeyChain {
	return &keychain{directory: directory}
}

func (k *keychain) LoadPrivateKey(name string) (crypto.PrivateKey, error) {
	keyBytes, err := k.loadBytes(fmt.Sprintf(PrivateKeyPrefix, name))
	return key(keyBytes), err
}

func (k *keychain) LoadPublicKey(name string) (crypto.PublicKey, error) {
	keyBytes, err := k.loadBytes(fmt.Sprintf(PublicKeyPrefix, name))
	return key(keyBytes), err
}

func (k *keychain) loadBytes(name string) ([]byte, error) {
	v, ok := k.cache[k.directory+name]

	if ok {
		return v, nil
	}

	f, err := os.Open(k.directory + name)

	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(f)

	return b, err
}

func (k *keychain) StorePublicKey(name string, pub crypto.PublicKey) error {
	return k.storeKey(fmt.Sprintf(PublicKeyPrefix, name), pub)
}

func (k *keychain) StorePrivateKey(name string, priv crypto.PrivateKey) error {
	return k.storeKey(fmt.Sprintf(PrivateKeyPrefix, name), priv)
}

func (k *keychain) storeKey(name string, key encoding.BinaryMarshaler) error {
	f, err := os.Create(k.directory + name)

	if err != nil {
		return err
	}

	b, err := key.MarshalBinary()

	if err != nil {
		return err
	}

	_, err = f.Write(b)

	return err
}

func ConvertBytesToPubKey(bytes []byte) crypto.PublicKey {
	return key(bytes)
}

func ConvertBytesToPrivKey(bytes []byte) crypto.PrivateKey {
	return key(bytes)
}