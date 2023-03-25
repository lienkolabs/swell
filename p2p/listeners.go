package p2p

import (
	"fmt"
	"net"

	"github.com/lienkolabs/swell"
	"github.com/lienkolabs/swell/crypto"
)

// for whom signed blocks should be forwarded
type BlockBroadcastNewtWork struct {
	attendees map[crypto.Hash]*SecureConnection
	comm      chan *swell.SignedBlock
}

func NewGatewayClient(address string, prv crypto.PrivateKey, rmt crypto.Token) (*SecureConnection, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	secure, err := PerformClientHandShake(conn, prv, rmt)
	if err != nil {
		return nil, err
	}
	return secure, nil
}

func NewGatewayNetwork(port int,
	prvKey crypto.PrivateKey, comm *swell.Communication) BlockBroadcastNewtWork {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	network := BlockBroadcastNewtWork{
		attendees: make(map[crypto.Hash]*SecureConnection),
		comm:      make(chan *swell.SignedBlock),
	}

	if err != nil {
		panic(err)
	}
	// listener loop
	go func() {
		for {
			conn, err := listener.Accept()
			if err == nil {
				secureConnection, err := PerformServerHandShake(conn, prvKey, comm.ValidateConn)
				if err != nil {
					conn.Close()
				} else {
					network.attendees[secureConnection.hash] = secureConnection
				}
			}
		}
	}()
	// write loop
	go func() {
		for {
			block := <-network.comm
			blockBytes := block.Block.Serialize()
			for hash, conn := range network.attendees {
				if err := conn.WriteMessage(blockBytes); err != nil {
					conn.conn.Close()
					delete(network.attendees, hash)
				}
			}
		}
	}()

	return network
}
