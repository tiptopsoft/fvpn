package security

import (
	"crypto/rand"
	"golang.org/x/crypto/chacha20poly1305"
	"io"
)

type CipherInterface interface {
	Encode(key, content []byte) ([]byte, error)
	Decode(key, cipherBuff []byte) ([]byte, error)
}

func NewCipher() CipherInterface {
	return &cipher{}
}

type cipher struct {
	key string
}

func (c *cipher) Encode(key, content []byte) ([]byte, error) {
	nonce := make([]byte, chacha20poly1305.NonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	cip, _ := chacha20poly1305.New(key)

	return cip.Seal(nil, nonce, content, nil), nil
}

func (c *cipher) Decode(key, cipherBuff []byte) ([]byte, error) {
	nonce := make([]byte, chacha20poly1305.NonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	cip, err := chacha20poly1305.New(key)
	if err != nil {
		return nil, err
	}

	return cip.Open(nil, nonce, cipherBuff, nil)
}
