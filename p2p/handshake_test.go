package p2p

import (
	"fmt"
	"net"
	"testing"

	"github.com/lienkolabs/swell"
	"github.com/lienkolabs/swell/crypto"
)

type ciphernonce struct {
	msg   []byte
	nonce []byte
}

var validator chan swell.ValidatedConnection = func() chan swell.ValidatedConnection {
	validator := make(chan swell.ValidatedConnection)
	go func() {
		validate := <-validator
		validate.Ok <- true
	}()
	return validator
}()

func TestSecureConnection(t *testing.T) {

	pubSv, prvSv := crypto.RandomAsymetricKey()
	listener, _ := net.Listen("tcp", ":7780")
	cipher := make(chan ciphernonce)
	go func() {
		conn, _ := listener.Accept()
		sec, err := PerformServerHandShake(conn, prvSv, validator)
		if err != nil {
			fmt.Println("---------", err)
			t.Error(err)
			return
		}
		var msg ciphernonce
		msg.msg, msg.nonce = sec.cipher.SealWithNewNonce([]byte("thats correct"))
		cipher <- msg
		msg.msg, msg.nonce = sec.cipherRemote.SealWithNewNonce([]byte("thats also correct"))
		cipher <- msg
	}()

	_, prvCl := crypto.RandomAsymetricKey()
	client, _ := net.Dial("tcp", ":7780")
	sec, err := PerformClientHandShake(client, prvCl, pubSv)
	if err != nil {
		t.Error(err)
	}
	msg := <-cipher
	msgData, err := sec.cipherRemote.OpenNewNonce(msg.msg, msg.nonce)
	if err != nil {
		t.Fatal(err)
	}
	if string(msgData) != "thats correct" {
		t.Fatalf("wrong message:%v", string(msgData))
	}
	msg = <-cipher
	msgData, err = sec.cipher.OpenNewNonce(msg.msg, msg.nonce)
	if err != nil {
		t.Fatal(err)
	}
	if string(msgData) != "thats also correct" {
		t.Fatalf("wrong message:%v", string(msgData))
	}

}
