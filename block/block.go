package block

import (
	"fmt"
	"time"

	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/swell/util"
)

type MutatingState interface {
	Validate(data []byte) bool
}

type Block struct {
	epoch         uint64
	Parent        crypto.Hash
	CheckPoint    uint64
	Publisher     crypto.Token
	PublishedAt   time.Time
	Transactions  [][]byte
	Hash          crypto.Hash
	FeesCollected uint64
	Signature     crypto.Signature
	validator     MutatingState
}

func NewBlock(parent crypto.Hash, checkpoint, epoch uint64, publisher crypto.Token, validator MutatingState) *Block {
	return &Block{
		Parent:       parent,
		epoch:        epoch,
		CheckPoint:   checkpoint,
		Publisher:    publisher,
		Transactions: make([][]byte, 0),
		validator:    validator,
	}
}

func (b *Block) Incorporate(transaction []byte) bool {
	if !b.validator.Validate(transaction) {
		return false
	}
	b.Transactions = append(b.Transactions, transaction)
	return true
}

func (b *Block) Epoch() uint64 {
	return b.epoch
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
	bytes := make([]byte, 0)
	util.PutUint64(b.epoch, &bytes)
	util.PutByteArray(b.Parent[:], &bytes)
	util.PutUint64(b.CheckPoint, &bytes)
	util.PutByteArray(b.Publisher[:], &bytes)
	util.PutTime(b.PublishedAt, &bytes)
	util.PutUint16(uint16(len(b.Transactions)), &bytes)
	for _, instruction := range b.Transactions {
		util.PutByteArray(instruction, &bytes)
	}
	util.PutByteArray(b.Hash[:], &bytes)
	util.PutUint64(b.FeesCollected, &bytes)
	return bytes
}

func ParseBlock(data []byte) *Block {
	position := 0
	block := Block{}
	block.epoch, position = util.ParseUint64(data, position)
	block.Parent, position = util.ParseHash(data, position)
	block.CheckPoint, position = util.ParseUint64(data, position)
	block.Publisher, position = util.ParseToken(data, position)
	block.PublishedAt, position = util.ParseTime(data, position)
	block.Transactions, position = util.ParseByteArrayArray(data, position)
	block.Hash, position = util.ParseHash(data, position)
	block.FeesCollected, position = util.ParseUint64(data, position)
	msg := data[0:position]
	block.Signature, _ = util.ParseSignature(data, position)
	if !block.Publisher.Verify(msg, block.Signature) {
		fmt.Println("wrong signature")
		return nil
	}
	return &block
}

func (b *Block) SetValidator(validator MutatingState) {
	b.validator = validator
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
	bulk.PutUint64("epoch", b.epoch)
	bulk.PutHex("parent", b.Parent[:])
	bulk.PutUint64("checkpoint", b.CheckPoint)
	bulk.PutHex("publisher", b.Publisher[:])
	bulk.PutTime("publishedAt", b.PublishedAt)
	bulk.PutUint64("instructionsCount", uint64(len(b.Transactions)))
	bulk.PutHex("hash", b.Parent[:])
	bulk.PutUint64("feesCollectes", b.FeesCollected)
	bulk.PutBase64("signature", b.Signature[:])
	return bulk.ToString()
}
