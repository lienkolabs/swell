package swell

import (
	"time"

	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/swell/util"
)

const Version = 0

type Block struct {
	Clock       uint64
	Parent      crypto.Hash
	CheckPoint  uint64
	Publisher   crypto.Token
	PublishedAt time.Time
	Events      Events
	Hash        crypto.Hash
	Signature   crypto.Signature
}

func (b *Block) Sign(token crypto.PrivateKey) {
	b.Signature = token.Sign(b.serializeWithoutSignature())
}

func (b *Block) Serialize() []byte {
	bytes := b.serializeWithoutSignature()
	util.PutSignature(b.Signature, &bytes)
	return bytes
}

func (b *Block) serializeWithoutSignature() []byte {
	bytes := make([]byte, Version)
	util.PutUint64(b.Clock, &bytes)
	util.PutByteArray(b.Parent[:], &bytes)
	util.PutUint64(b.CheckPoint, &bytes)
	util.PutByteArray(b.Publisher[:], &bytes)
	util.PutTime(b.PublishedAt, &bytes)
	util.PutUint16(uint16(len(b.Events)), &bytes)
	for _, instruction := range b.Events {
		util.PutByteArray(instruction, &bytes)
	}
	util.PutByteArray(b.Hash[:], &bytes)
	return bytes
}

func ParseBlock(data []byte) *Block {
	position := 0
	block := Block{}
	block.Clock, position = util.ParseUint64(data, position)
	block.Parent, position = util.ParseHash(data, position)
	block.CheckPoint, position = util.ParseUint64(data, position)
	block.Publisher, position = util.ParseToken(data, position)
	block.PublishedAt, position = util.ParseTime(data, position)
	if position+1 >= len(data) {
		return nil
	}
	length := int(data[position+0]) | int(data[position+1])<<8
	position += 2
	block.Events = make(Events, length)
	for n := 0; n < length; n++ {
		var newEvent []byte
		newEvent, position = util.ParseByteArray(data, position)
		block.Events[n] = Event(newEvent)
	}
	block.Hash, position = util.ParseHash(data, position)
	msg := data[0:position]
	block.Signature, _ = util.ParseSignature(data, position)
	if !block.Publisher.Verify(msg, block.Signature) {
		return nil
	}
	return &block
}

func GetBlockEpoch(data []byte) uint64 {
	if len(data) < 8 {
		return 0
	}
	epoch, _ := util.ParseUint64(data, 0)
	return epoch
}

func (b *Block) JSONSimple() string {
	bulk := &util.JSONBuilder{}
	bulk.PutUint64("epoch", b.Clock)
	bulk.PutHex("parent", b.Parent[:])
	bulk.PutUint64("checkpoint", b.CheckPoint)
	bulk.PutHex("publisher", b.Publisher[:])
	bulk.PutTime("publishedAt", b.PublishedAt)
	bulk.PutUint64("instructionsCount", uint64(len(b.Events)))
	bulk.PutHex("hash", b.Parent[:])
	bulk.PutBase64("signature", b.Signature[:])
	return bulk.ToString()
}

type Signature struct {
	Hash      crypto.Hash
	Token     crypto.Token
	Signature crypto.Signature
}

type SignedBlock struct {
	Block      *Block
	Signatures []Signature
}

type SignedBlocks []*SignedBlock

func (blocks SignedBlocks) Less(i, j int) bool {
	return blocks[i].Block.Clock < blocks[j].Block.Clock
}

func (blocks SignedBlocks) Len() int {
	return len(blocks)
}

func (blocks SignedBlocks) Swap(i, j int) {
	blocks[i], blocks[j] = blocks[j], blocks[i]
}

type Blocks []*Block

func (blocks Blocks) Less(i, j int) bool {
	return blocks[i].Clock < blocks[j].Clock
}

func (blocks Blocks) Len() int {
	return len(blocks)
}

func (blocks Blocks) Swap(i, j int) {
	blocks[i], blocks[j] = blocks[j], blocks[i]
}
