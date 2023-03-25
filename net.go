package swell

import (
	"time"

	"github.com/lienkolabs/swell/crypto"
)

type PeerRequest struct {
	Token    crypto.Hash
	Response chan bool
}

type Checksum struct {
	Token   crypto.Hash
	Check   []byte
	Confirm chan bool
}

type SyncRequest struct {
	Starting chan uint64
	Data     chan []byte
	Ok       chan bool
}

type ValidatedConnection struct {
	Token crypto.Hash
	Ok    chan bool
}

type Communication struct {
	PeerRequest     chan *PeerRequest // Node receives new peer requests from network
	NewBlock        chan *Block       // Node publishes to or receives new blocks from the network
	BlockSignature  chan *Signature   // Node publishes to or receives signatures from the network
	Checkpoint      chan *SignedBlock // Node publishes new checkpoint to observers network
	Checksum        chan *Checksum    // Node publishes to or receives checksums from the network
	Synchronization chan SyncRequest  // Node receives sync request
	ValidateConn    chan ValidatedConnection
	Events          chan Event
}

func NewCommunication() *Communication {
	return &Communication{
		PeerRequest:     make(chan *PeerRequest),
		NewBlock:        make(chan *Block),
		BlockSignature:  make(chan *Signature),
		Checkpoint:      make(chan *SignedBlock),
		Checksum:        make(chan *Checksum),
		Synchronization: make(chan SyncRequest),
		ValidateConn:    make(chan ValidatedConnection),
		Events:          make(chan Event),
	}
}

type ConsensusEngine func(BlockChain) *Communication

func IntervalToNewEpoch(epoch uint64, genesis time.Time) time.Duration {
	return time.Until(genesis.Add(time.Duration(int64(epoch) * 1000000000)))
}

func LauchNewGenesisConsensus(egine ConsensusEngine) {
	//pool := NewInstructionPool()

	//processInstruction := make(chan instructions.Instruction)
	go func() {
		for {
			select {
			//case newInstruction := <-processInstruction:
			//pool.Queue(newInstruction)
			}
		}
	}()
}
