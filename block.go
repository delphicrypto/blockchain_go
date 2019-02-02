package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
	"math/big"
)



// Block represents a block in the blockchain
type Block struct {
	Timestamp     int64
	Transactions  []*Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
	Height        int
	Target	  	  *big.Int
}

// NewBlock creates and returns Block
func NewBlock(transactions []*Transaction, prevBlockHash []byte, height int, target *big.Int) *Block {
	block := &Block{time.Now().UnixNano(), transactions, prevBlockHash, []byte{}, 0, height, target}
	
	return block
}

// NewGenesisBlock creates and returns genesis Block
func NewGenesisBlock(coinbase *Transaction) *Block {
	target := targetFromTargetBits(initialTargetBits)
	return NewBlock([]*Transaction{coinbase}, []byte{}, 0, target)
}

// HashTransactions returns a hash of the transactions in the block
func (b *Block) HashTransactions() []byte {
	var transactions [][]byte

	for _, tx := range b.Transactions {
		transactions = append(transactions, tx.Serialize())
	}
	mTree := NewMerkleTree(transactions)

	return mTree.RootNode.Data
}

// Serialize serializes the block
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

func (b *Block) Validate(chainTarget *big.Int) bool {

//check that the targetBits is correct
	if b.Target.Cmp(chainTarget) != 0 {
		return false
	}
	pow := NewProofOfWork(b)
	return pow.Validate()
}

// DeserializeBlock deserializes a block
func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}

