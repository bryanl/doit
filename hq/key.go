package hq

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/ssh"
)

// KeyRing stores a RSA key.
type KeyRing struct {
	PrivateKey *rsa.PrivateKey
}

// NewKeyRing creates a new KeyRing.
func NewKeyRing() (*KeyRing, error) {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return nil, err
	}

	return &KeyRing{key}, nil
}

// EncodedPublicKey converts the public key to rfc 4716 format.
func (k *KeyRing) EncodedPublicKey() (string, error) {
	publicKey, err := ssh.NewPublicKey(&k.PrivateKey.PublicKey)
	if err != nil {
		return "", err
	}

	b := publicKey.Marshal()

	encoded := base64.StdEncoding.EncodeToString(b)

	return fmt.Sprintf("%s %s\n", publicKey.Type(), encoded), nil
}

func (k *KeyRing) EncodedPrivateKey()
