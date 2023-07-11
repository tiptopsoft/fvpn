package security

import (
	"crypto/rand"
	"fmt"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/util"
	"golang.org/x/crypto/curve25519"
	"io"
	"testing"
	"time"
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

	cip1 := NewCipher(privateKey, pubKey2)
	cip2 := NewCipher(privateKey2, pubKey)

	s := "hello, myworldhello, myworldhello, myworldhello, myworldhello, myworldhello, myworldhello, myworldhello, myworldhello, myworldhello, myworld"
	sBuff := []byte(s)
	frame := packet.NewFrame()
	h, _ := packet.NewHeader(util.MsgTypePacket, "123456")
	headerBuff, _ := packet.Encode(h)
	copy(frame.Packet, headerBuff)
	copy(frame.Packet[44:], sBuff)
	size := len(headerBuff) + len(sBuff)

	st := time.Now()
	fmt.Println("before encoded: ", frame.Packet[:size])

	encodedBuff, err := cip1.Encode(frame.Packet[:size])
	if err != nil {
		panic(err)
	}

	fmt.Println("After encoded: ", encodedBuff[:len(encodedBuff)])
	et := time.Since(st)
	fmt.Println("cost: ", et)

	decodedBuff, err := cip2.Decode(encodedBuff)
	if err != nil {
		panic(err)
	}

	fmt.Println("after decode: ", decodedBuff)

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
