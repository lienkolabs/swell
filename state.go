package swell

import "github.com/lienkolabs/swell/crypto"

type State interface {
	LastCheckPoint() Checkpoint
	ChecksumJob() chan crypto.Hash
}

type Checkpoint interface {
	Clock() uint64
}
