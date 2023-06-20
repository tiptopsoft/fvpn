package security

import (
	"crypto/rand"
	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/curve25519"
	"io"
)

const (
	NoiseKeySize = 32
)

type (
	NoisePrivateKey [NoiseKeySize]byte
	NoisePublicKey  [NoiseKeySize]byte
	NoiseSharedKey  [NoiseKeySize]byte
)

type CipherFunc interface {
	Encode(content []byte) ([]byte, error)
	Decode(cipherBuff []byte) ([]byte, error)
}

func NewCipher(privateKey NoisePrivateKey, pubKey NoisePublicKey) CipherFunc {
	shareKey := privateKey.NewSharedKey(pubKey)
	return &cipher{
		key: shareKey,
	}
}

type cipher struct {
	key NoiseSharedKey
}

func (c *cipher) Encode(content []byte) ([]byte, error) {
	nonce := make([]byte, chacha20poly1305.NonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	cip, _ := chacha20poly1305.New(c.key[:])

	return cip.Seal(nil, nonce, content, nil), nil
}

func (c *cipher) Decode(cipherBuff []byte) ([]byte, error) {
	nonce := make([]byte, chacha20poly1305.NonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	cip, err := chacha20poly1305.New(c.key[:])
	if err != nil {
		return nil, err
	}

	return cip.Open(nil, nonce, cipherBuff, nil)
}

func NewPrivateKey() (npk NoisePrivateKey, err error) {
	_, err = rand.Read(npk[:])
	if err != nil {
		return
	}
	npk[0] &= 248
	npk[31] &= 127
	npk[31] |= 64
	return
}

func (npk NoisePrivateKey) NewPubicKey() (npc NoisePublicKey) {
	privateKey := (*[32]byte)(&npk)
	pubKey := (*[32]byte)(&npc)
	curve25519.ScalarBaseMult(pubKey, privateKey)
	return
}

func (npk NoisePrivateKey) NewSharedKey(npc NoisePublicKey) (shareKey NoiseSharedKey) {
	sk := (*[32]byte)(&shareKey)
	pubKey := (*[32]byte)(&npc)
	priKey := (*[32]byte)(&npk)
	curve25519.ScalarMult(sk, priKey, pubKey)
	return
}
