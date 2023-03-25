package p2p

import (
	"fmt"
	"net"

	"github.com/lienkolabs/swell/crypto"
)

const (
	validatorNetworkMsg byte = iota
	broadcastValidatedMsg
	broadcastBlockMsg
	pingMsg
	pongMsg
)

type ValidateConnection interface {
	ValidateConnection(token crypto.Token) chan bool
}

type ValidatorNetwork map[crypto.Hash]*SecureConnection

func (v ValidatorNetwork) Broadcast(msg *NetworkMessageTemplate) {
	msgToSend := msg.Serialize()
	for _, peer := range v {
		peer.WriteMessage(msgToSend)
	}
}
func NewValidatorNetwork(port int, prvKey crypto.PrivateKey, comm chan *HashedEventBytes,
	validator ValidateConnection, dial map[crypto.Token]string) ValidatorNetwork {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	network := make(ValidatorNetwork)
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			conn, err := listener.Accept()
			if err == nil {
				secureConnection, err := PerformServerHandShake(conn, prvKey, validator)
				if err != nil {
					conn.Close()
				} else {
					network[secureConnection.hash] = secureConnection
					handleValidatorConnection(secureConnection, comm)
				}
			}
		}
	}()
	for publicKey, address := range dial {
		go func(address string, publicKey crypto.Token) {
			net, err := net.Dial("tcp", address)
			if err != nil {
				return
			}
			conn, err := PerformClientHandShake(net, prvKey, publicKey)
			if err != nil {
				return
			}
			network[crypto.HashToken(publicKey)] = conn
			handleValidatorConnection(conn, comm)
		}(address, publicKey)
	}
	return network
}

func handleValidatorConnection(conn *SecureConnection, comm chan *HashedEventBytes) {
	for {
		data, err := conn.ReadMessage()
		if err != nil {
			conn.conn.Close()
			return
		}
		hashed := HashedEventBytes{msg: data}
		hashed.hash = crypto.Hasher(data)
		hashed.epoch = int(event.Clock())
		comm <- &hashed
	}
}
