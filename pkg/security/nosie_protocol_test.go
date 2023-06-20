package security

import (
	"crypto/rand"
	"fmt"
	"golang.org/x/crypto/curve25519"
	"io"
	"testing"
)

// DH test
func TestCurve(t *testing.T) {

	var privateKey [32]byte
	_, err := io.ReadFull(rand.Reader, privateKey[:])
	if err != nil {
		panic(err)
	}

	var pubKey [32]byte
	curve25519.ScalarBaseMult(&pubKey, &privateKey)

	var privateKey2 [32]byte
	_, err = io.ReadFull(rand.Reader, privateKey2[:])
	if err != nil {
		panic(err)
	}

	var pubKey2 [32]byte
	curve25519.ScalarBaseMult(&pubKey2, &privateKey2)

	//assert
	var shared [32]byte
	curve25519.ScalarMult(&shared, &privateKey, &pubKey2)

	fmt.Println(shared)
	var shared2 [32]byte
	curve25519.ScalarMult(&shared2, &privateKey2, &pubKey)

	fmt.Println(shared2)

}

func TestNewPrivateKey(t *testing.T) {

	privateKey, _ := NewPrivateKey()
	pubKey := privateKey.NewPubicKey()

	privateKey2, _ := NewPrivateKey()
	pubKey2 := privateKey2.NewPubicKey()

	shareKey := privateKey.NewSharedKey(pubKey2)
	shareKey2 := privateKey2.NewSharedKey(pubKey)

	fmt.Println(shareKey)
	fmt.Println(shareKey2)
}
