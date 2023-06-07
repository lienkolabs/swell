package trusted

import (
	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/swell/util"
)

const (
	IBlockBroadcastInclude  byte = iota // request to be included in block broadcastimg pool
	IBlockBroadcastIncluded             // accept request to be included in block boradcasting pool
	IBlockBroadcastRemove               // accept request to be included in block boradcasting pool
	IBlockBroadcastRemoved              // accept request to be included in block boradcasting pool
	INewBlock                           // informs avalibility of new block
	IBlockRequest                       // request transmission of new block
	IBlockSend                          // send new block
	ISendEvent                          // send new event to be considered by the network

	Version = 0
)

type Serializer interface {
	Serialize() []byte
}

type BlockDescription struct {
	Age  uint64
	Hash crypto.Hash
}

func (n *BlockDescription) Serialize() []byte {
	ageBytes := make([]byte, 8)
	util.PutUint64(n.Age, &ageBytes)
	hashBytes := make([]byte, crypto.Size)
	util.PutByteArray(n.Hash[:], &hashBytes)
	return append([]byte{INewBlock}, append(ageBytes, hashBytes...)...)
}

func ParseBlockDescription(bytes []byte) *BlockDescription {
	if len(bytes) != 16+1 || bytes[0] != INewBlock {
		return nil
	}
	age, position := util.ParseUint64(bytes, 1)
	hash, _ := util.ParseHash(bytes, position)
	return &BlockDescription{Age: age, Hash: hash}
}
