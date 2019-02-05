package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"fmt"
	"strconv"
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
	SolutionHash []byte
	Solution	  []int
	ProblemGraphHash []byte
}

// NewBlock creates and returns Block
func NewBlock(transactions []*Transaction, prevBlockHash []byte, height int, target *big.Int, solHash []byte, solution []int, pgHash []byte) *Block {
	block := &Block{time.Now().UnixNano(), transactions, prevBlockHash, []byte{}, 0, height, target, solHash, solution, pgHash}
	
	return block
}

// NewGenesisBlock creates and returns genesis Block
func NewGenesisBlock(coinbase *Transaction, pgHash []byte) *Block {
	target := targetFromTargetBits(initialTargetBits)
	
	return NewBlock([]*Transaction{coinbase}, []byte{}, 0, target, []byte{}, []int{}, pgHash)
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

//NicePrint print nicely the block properties
func (b *Block) NicePrint(bc *Blockchain) {
	fmt.Printf("\n")
	printGreen(fmt.Sprintf("============ Block %d ============\n", b.Height))
	printBlue(fmt.Sprintf("Hash: %x\n", b.Hash))
	fmt.Printf("Prev: %x\n", b.PrevBlockHash)
	fmt.Printf("Difficulty: %d\n", targetToDifficulty(b.Target))
	prevBlock, _ := bc.GetBlockFromHash(b.PrevBlockHash)
	time := (b.Timestamp - prevBlock.Timestamp) / 1e9
	fmt.Printf("Time: %d seconds\n", time)
	blockchainTarget := bc.CalculateTarget(b.Height)
	validBlock := b.Validate(blockchainTarget)
	if validBlock {
		printGreen(fmt.Sprintf("PoW: %s\n", strconv.FormatBool(validBlock)))
	} else {
		printRed(fmt.Sprintf("PoW: %s\n", strconv.FormatBool(validBlock)))
	}

	if len(b.SolutionHash) > 0 {
		printGreen(fmt.Sprintf("Solution to %x: ", b.SolutionHash))
		fmt.Println(b.Solution)
	} else {
		printRed("No solution\n")
	}

	if len(b.ProblemGraphHash) > 0 {
		printGreen(fmt.Sprintf("New Problem %x \n", b.ProblemGraphHash))
		// pg, err := bc.GetProblemGraphFromHash(b.ProblemGraphHash)
		// if err == nil {
		// 	pg.NicePrint()
		// }
		
	} else {
		printRed("No problem ;)\n")
	}

	for _, tx := range b.Transactions {
		printYellow(fmt.Sprintln(tx))
	}
	fmt.Printf("\n")
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

