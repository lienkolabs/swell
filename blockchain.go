package swell

import (
	"time"
)

type BlockChain struct {
	GenesisTime     time.Time
	TotalStake      uint64
	Epoch           uint64
	CurrentState    State
	RecentBlocks    SignedBlocks
	CandidateBlocks map[uint64]SignedBlocks
}
