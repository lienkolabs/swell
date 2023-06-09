package consensus

import (
	"time"

	"github.com/lienkolabs/swell/block"
	"github.com/lienkolabs/swell/crypto"
)

/*
	Interface:

		Receives New Instructions
		Receives New Blocks
		Receives Block signatures
		Peer validation request
*/

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
	NewBlock        chan *block.Block // Node publishes to or receives new blocks from the network
	BlockSignature  chan *Signature   // Node publishes to or receives signatures from the network
	Checkpoint      chan *SignedBlock // Node publishes new checkpoint to observers network
	Checksum        chan *Checksum    // Node publishes to or receives checksums from the network
	Synchronization chan SyncRequest  // Node receives sync request
	ValidateConn    chan ValidatedConnection
	Instructions    chan []byte
}

func NewCommunication() *Communication {
	return &Communication{
		PeerRequest:     make(chan *PeerRequest),
		NewBlock:        make(chan *block.Block),
		BlockSignature:  make(chan *Signature),
		Checkpoint:      make(chan *SignedBlock),
		Checksum:        make(chan *Checksum),
		Synchronization: make(chan SyncRequest),
		ValidateConn:    make(chan ValidatedConnection),
		Instructions:    make(chan []byte]),
	}
}

type ConsensusEngine func(BlockChain) *Communication

func IntervalToNewEpoch(epoch uint64, genesis time.Time) time.Duration {
	return time.Until(genesis.Add(time.Duration(int64(epoch) * 1000000000)))
}

/*
	func (c *Consensus) PushNewBlock(block *instructions.Block) {
		epoch := block.Epoch
		newSignedBlock := SignedBlock{
			Block:      block,
			Signatures: []Signature{},
		}
		if blocks, ok := c.blockChain.CandidateBlocks[epoch]; ok {
			blocks = append(blocks, &newSignedBlock)
		} else {
			c.blockChain.CandidateBlocks[epoch] = SignedBlocks{&newSignedBlock}
		}
		consensus := c.blockChain.Engine.GetConsensus(c.blockChain.CandidateBlocks[epoch])
		if consensus != nil {
			c.blockChain.RecentBlocks = append(c.blockChain.RecentBlocks, &newSignedBlock)
			delete(c.blockChain.CandidateBlocks, epoch)
			sort.Sort(c.blockChain.RecentBlocks)
			c.blockChain.IncorporateBlocks()
		}
	}

func (b *BlockChain) IncoporateBlocks() {

}

	func IntervalToNewEpoch(epoch uint64, genesis time.Time) time.Duration {
		return time.Until(genesis.Add(time.Duration(int64(epoch) * 1000000000)))
	}

	func (b *BlockChain) AppendSignature(signature Signature) {
		if len(b.CandidateBlocks) != 0 {
			for _, block := range b.CandidateBlocks {
				if block.Block.Hash.Equals(signature.Hash[:]) {
					for _, signBlock := range block.Signatures {
						if bytes.Equal(signBlock.Token, signature.Token) {
							return
						}
					}
					newSignature := Signature{
						Hash:      signature.Hash,
						Token:     signature.Token,
						Signature: signature.Signature,
					}
					block.Signatures = append(block.Signatures, newSignature)
					//if b.Engine.IsConsensus()
					return
				}
			}
		}
		// TODO append to recent blocks also
	}

	func (b *BlockChain) GetLastValidator() (*instructions.Validator, uint64) {
		starting := b.CurrentState.Epoch
		if len(b.RecentBlocks) == 0 || b.RecentBlocks[0].Epoch != starting+1 {
			return &instructions.Validator{
				State:     b.CurrentState,
				Mutations: instructions.NewMutation(),
			}
		}
		sequential := make([]*instructions.Block, 0)
		for _, block := range b.RecentBlocks {
			if block.Epoch != starting+1 {
				break
			}
			starting += 1
			sequential = append(sequential, block)
		}
		return &instructions.Validator{
			State:     b.CurrentState,
			Mutations: instructions.GroupBlockMutations(sequential),
		}
	}
*/
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

/*
func NewGenesisConsensus(engine Consensual) {
	pool := NewInstructionPool()
	state, token := instructions.NewGenesisState()
	timeOfGenesis := time.Now()
	epoch := int64(0)
	blockTick := time.NewTimer(time.Until(timeOfGenesis.Add(time.Duration((epoch + 1) * 1000000000))))
	go func() {
		for {
			starting := <-blockTick.C
			if engine.IsLeader(starting) {
				BlockBuilder()
			}
		}
	}()
}
*/
