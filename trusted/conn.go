package trusted

import (
	"errors"
	"net"

	"github.com/lienkolabs/swell/crypto"
)

var ErrMessageTooLarge = errors.New("message size cannot be larger than 65.536 bytes")
var ErrInvalidSignature = errors.New("signature is invalid")

type SignedConnection struct {
	token         crypto.Token
	key           crypto.PrivateKey
	conn          net.Conn
	terminate     chan struct{}
	blockListener bool
}

func (s *SignedConnection) WriteMessage(msg []byte) error {
	lengthWithSignature := len(msg) + crypto.SignatureSize
	if lengthWithSignature > 1<<40-1 {
		return ErrMessageTooLarge
	}
	msgToSend := []byte{byte(lengthWithSignature), byte(lengthWithSignature >> 8),
		byte(lengthWithSignature >> 16), byte(lengthWithSignature >> 24), byte(lengthWithSignature >> 32)}
	signature := s.key.Sign(msg)
	msgToSend = append(append(msgToSend, msg...), signature[:]...)
	if n, err := s.conn.Write(msgToSend); n != lengthWithSignature+4 {
		return err
	}
	return nil
}

func (s *SignedConnection) readMessageWithoutCheck() ([]byte, error) {
	lengthBytes := make([]byte, 5)
	if n, err := s.conn.Read(lengthBytes); n != 5 {
		return nil, err
	}
	lenght := int(lengthBytes[0]) + (int(lengthBytes[1]) << 8) + (int(lengthBytes[2]) << 16) + (int(lengthBytes[3]) << 24) + (int(lengthBytes[4]) << 32)
	msg := make([]byte, lenght)
	if n, err := s.conn.Read(msg); n != int(lenght) {
		return nil, err
	}
	return msg, nil
}

func (s *SignedConnection) read() ([]byte, error) {
	bytes, err := s.readMessageWithoutCheck()
	if err != nil {
		return nil, err
	}
	msg := bytes[0 : len(bytes)-crypto.SignatureSize]
	var signature crypto.Signature
	copy(signature[:], bytes[len(bytes)-crypto.SignatureSize:])
	if !s.token.Verify(msg, signature) {
		return nil, ErrInvalidSignature
	}
	return msg, nil
}

func (s *SignedConnection) Listen() chan Message {
	newMessages := make(chan Message)
	go func() {
		for {
			data, err := s.read()
			if err != nil {
				s.terminate <- struct{}{}
				return
			}
			newMessages <- Message{token: s.token, msg: data}
		}
	}()
	return newMessages
}

func ConnectGateway(address string, prvKey crypto.PrivateKey, pubKey crypto.Token, messages chan Message) *SignedConnection {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil
	}
	secureConnection, err := PerformClientHandShake(conn, prvKey, pubKey)
	if err != nil {
		conn.Close()
		return nil
	}
	secureConnection.terminate = make(chan struct{})
	go func() {
		for {
			data, err := secureConnection.read()
			if err != nil {
				return
			}
			messages <- Message{token: pubKey, msg: data}
		}
	}()
	return secureConnection
}

type connResult struct {
	hash crypto.Hash
	conn *SignedConnection
}
