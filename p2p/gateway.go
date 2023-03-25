package p2p

import (
	"fmt"

	"github.com/lienkolabs/swell"
	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/swell/util"
)

const maxEpochReceiveMessage = 100
const Version = 0

// Events are received from trusted gateways. 

type Gateways struct {
	authorized []crypto.Token
}

func (g *Gateways) ValidateConnection(token crypto.Token) chan bool {
	check := make(chan bool)
	go func() {
		for _, gateway := range g.authorized {
			if token.Equal(gateway) {
				ok <- true
				return
			}
		}
		ok <- false
	}()
	return check
}

type HashedEventBytes struct {
	msg     []byte
	hash    crypto.Hash
	clock   int
}

func getEventClock(event []byte) int {
	if event[0] != Version || len(event) < 9 {
		return -1
	}
	clock, _ := util.ParseUint64(event,1)
	if clock == 0 {
		return - 1
	}
	return int(clock)
}

type EventBroker chan *HashedEventBytes

func (e EventBroker) Queue(event []byte) {
	e <- &HashedEventBytes{
		msg: event,
		hash: crypto.Hasher(msg),
		clock: getEventClock(msg),
	}
}


func NewEventBroker(
	token crypto.PrivateKey,
	peers *ValidatorNetwork,
	comm *swell.Communication,
	newBlockSignal chan uint64,
	epoch uint64,
) EventBroker {
	broker := make(EventBroker)
	recentHashes := make([]map[crypto.Hash]struct{}, maxEpochReceiveMessage)
	for n := 0; n < maxEpochReceiveMessage; n++ {
		recentHashes[n] = make(map[crypto.Hash]struct{})
	}
	currentEpoch := int(epoch)
	go func() {
		for {
			select {
			case hashInst := <-broker:
				if deltaEpoch := currentEpoch - int(hashInst.epoch); deltaEpoch < 100 && deltaEpoch >= 0 {
					if _, exists := recentHashes[deltaEpoch][hashInst.hash]; !exists {
						recentHashes[deltaEpoch][hashInst.hash] = struct{}{}
						comm.Events <- &HashTransaction{
								Transaction: transaction,
								Hash:        hashInst.hash,
							}
							if hashInst.nonpeer {
								message := NewNetworkMessage(BroadcastInstruction(hashInst.msg), token, false)
								peers.Broadcast(message)
							}
						}
						//broker <- hashInst
						// if instruction was not received from peer it should be broadcasted
					}
				}
			case newEpoch := <-newBlockSignal:
				deltaEpoch := int(newEpoch) - currentEpoch
				if deltaEpoch != 1 {
					panic(fmt.Sprintf("TODO: decide what to do... %v, %v", newEpoch, currentEpoch))
				}
				recentHashes = append(recentHashes[1:], make(map[crypto.Hash]struct{}))
				currentEpoch = int(newEpoch)
				fmt.Printf("current epoch: %v\n", currentEpoch)
			}
		}
	}()
	return broker
}
