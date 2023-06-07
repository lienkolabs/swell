package trusted

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/swell/util"
)

type Message struct {
	token crypto.Token
	msg   []byte
}

// for whom signed blocks should be forwarded
type Gateway struct {
	mu        sync.Mutex
	key       crypto.PrivateKey
	outbound  map[crypto.Token]*SignedConnection
	inbound   map[crypto.Token]*SignedConnection
	terminate chan struct{}
	messages  chan Message
}

func (g *Gateway) NewMessage(kind byte, data []byte) []byte {
	output := []byte{kind}
	now := time.Now()
	util.PutTime(now, &output)
	util.PutByteArray(data, &output)
	signature := g.key.Sign(output)
	util.PutSignature(signature, &output)
	return output
}

func NewGateway(port int, prvKey crypto.PrivateKey, validator ValidateConnection) *Gateway {

	router := Gateway{
		outbound:  make(map[crypto.Token]*SignedConnection),
		inbound:   make(map[crypto.Token]*SignedConnection),
		terminate: make(chan struct{}),
		messages:  make(chan Message),
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		panic(err)
	}

	// listener loop
	go func() {
		for {
			conn, err := listener.Accept()
			if err == nil {
				secureConnection, err := PerformServerHandShake(conn, prvKey, validator)
				if err != nil {
					conn.Close()
				} else {
					router.mu.Lock()
					router.outbound[secureConnection.token] = secureConnection
					router.mu.Unlock()
				}
			}
		}
	}()

	// write/termination loop
	go func() {
		for {
			select {
			case <-router.terminate:
				router.mu.Lock()
				for hash, conn := range router.outbound {
					conn.conn.Close()
					delete(router.outbound, hash)
				}
				close(router.terminate)
				router.mu.Unlock()
				return
				//case block := <-router.inbound:
				//	for hash, conn := range router.outbound {
				//		if err := conn.WriteMessage(block); err != nil {
				//			conn.conn.Close()
				//				router.mu.Lock()
				//				delete(router.outbound, hash)
				//				router.mu.Unlock()
				//		}
				//				}
			}
		}
	}()

	return &router
}
