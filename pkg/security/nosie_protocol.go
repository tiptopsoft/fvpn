// Copyright 2023 Tiptopsoft, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package security

import (
	"crypto/rand"
	"github.com/tiptopsoft/fvpn/pkg/log"
	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/curve25519"
)

const (
	NoiseKeySize = 32
)

var (
	logger = log.Log()
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

	logger.Debugf("generate shared key: %v", shareKey)
	nonce := make([]byte, chacha20poly1305.NonceSize)

	//if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
	//	return nil
	//}
	copy(nonce, shareKey[:chacha20poly1305.NonceSize])
	return &cipher{
		key:   shareKey,
		nonce: nonce,
	}
}

type cipher struct {
	key   NoiseSharedKey
	nonce []byte
}

func (c *cipher) Encode(content []byte) ([]byte, error) {
	logger.Debugf(">>>>>>==before encode: %v, key: %v", content, c.key)
	cip, err := chacha20poly1305.New(c.key[:])
	if err != nil {
		return nil, err
	}

	encoded := cip.Seal(nil, c.nonce, content, nil)
	logger.Debugf(">>>>>>==after encode: %v", encoded)
	return encoded, nil
}

func (c *cipher) Decode(cipherBuff []byte) ([]byte, error) {
	logger.Debugf(">>>>>before decode: %v, key: %v", cipherBuff, c.key)
	cip, err := chacha20poly1305.New(c.key[:])
	if err != nil {
		return nil, err
	}

	decoded, err := cip.Open(nil, c.nonce, cipherBuff, nil)
	logger.Debugf(">>>>after decode: %v", decoded)
	return decoded, err
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
