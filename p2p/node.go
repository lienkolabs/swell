package p2p

import (
	"time"

	"github.com/lienkolabs/swell"
	"github.com/lienkolabs/swell/crypto"
)

const (
	validationNodePort             = 7080
	blockBroadcastPort             = 7801
	messageReceiveConnectionPort   = 7802
	messageBroadcastConnectionPort = 7803
	syncPort                       = 7804
)

var BlockWindow, _ = time.ParseDuration("1s")
var GenesisTime = time.Date(2021, time.November, 18, 0, 0, 0, 0, time.UTC)

type MsgValidator struct {
	msg []byte
	ok  chan bool
}

type MsgValidatorChan chan *MsgValidator

func NewNode(prvKey crypto.PrivateKey,
	trusted map[crypto.Token]string,
	comm *swell.Communication,
	epoch uint64,
) {
	//
	newBlockSignal := make(chan uint64)
	peers := ValidatorNetwork(ConnectTCPPool(trusted, prvKey))
	instructionBroker := NewInstructionBroker(prvKey, &peers, comm, newBlockSignal, epoch)
	NewInstructionNetwork(messageReceiveConnectionPort, prvKey, instructionBroker, comm.ValidateConn)
	attendees := NewGatewayNetwork(
		blockBroadcastPort,
		prvKey,
		comm,
	)
	go func() {
		for {
			signedBlock := <-comm.Checkpoint
			newBlockSignal <- signedBlock.Block.Clock + 1
			attendees.comm <- signedBlock
		}
	}()
}
