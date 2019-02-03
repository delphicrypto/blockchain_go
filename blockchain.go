package main

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"math/big"
	"github.com/boltdb/bolt"
)

const dbFile = "blockchain_%s.db"
const blocksBucket = "blocks"
const problemsBucket = "problems"
const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"
const blocksPerTargetUpdate = 64
const initialTargetBits = 16
const targetBlocksPerMinute = 6
const secondsPerMinute = 60
const maxTargetChange = 4
// Blockchain implements interactions with a DB
type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

// CreateBlockchain creates a new blockchain DB
func CreateBlockchain(address, nodeID string) *Blockchain {
	dbFile := fmt.Sprintf(dbFile, nodeID)
	if dbExists(dbFile) {
		fmt.Println("Blockchain already exists.")
		bc := NewBlockchain(nodeID)
		return bc

		//os.Exit(1)
	}

	var tip []byte

	cbtx := NewCoinbaseTX(address, genesisCoinbaseData)
	genesis := NewGenesisBlock(cbtx)
	pow := NewProofOfWork(genesis)
	nonce, hash := pow.Run()
	genesis.Hash = hash[:]
	genesis.Nonce = nonce
	
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			log.Panic(err)
		}
		err = b.Put(genesis.Hash, genesis.Serialize())
		if err != nil {
			log.Panic(err)
		}
		err = b.Put([]byte("l"), genesis.Hash)
		if err != nil {
			log.Panic(err)
		}
		tip = genesis.Hash
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(problemsBucket))
		if err != nil {
			log.Panic(err)
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}

	return &bc
}

// NewBlockchain creates a new Blockchain with genesis Block
func NewBlockchain(nodeID string) *Blockchain {
	dbFile := fmt.Sprintf(dbFile, nodeID)
	if dbExists(dbFile) == false {
		fmt.Println("No existing blockchain found. Create one first.")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte("l"))

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}

	return &bc
}

// AddBlock saves the block into the blockchain
func (bc *Blockchain) AddBlock(block *Block) {
	err := bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		blockInDb := b.Get(block.Hash)

		if blockInDb != nil {
			return nil
		}

		blockData := block.Serialize()
		err := b.Put(block.Hash, blockData)
		if err != nil {
			log.Panic(err)
		}

		lastHash := b.Get([]byte("l"))
		lastBlockData := b.Get(lastHash)
		lastBlock := DeserializeBlock(lastBlockData)

		if block.Height > lastBlock.Height {
			err = b.Put([]byte("l"), block.Hash)
			if err != nil {
				log.Panic(err)
			}
			bc.tip = block.Hash
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

// FindTransaction finds a transaction by its ID
func (bc *Blockchain) FindTransaction(ID []byte) (Transaction, error) {
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return *tx, nil
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return Transaction{}, errors.New("Transaction is not found")
}

// FindUTXO finds all unspent transaction outputs and returns transactions with spent outputs removed
func (bc *Blockchain) FindUTXO() map[string]TXOutputs {
	UTXO := make(map[string]TXOutputs)
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Vout {
				// Was the output spent?
				if spentTXOs[txID] != nil {
					for _, spentOutIdx := range spentTXOs[txID] {
						if spentOutIdx == outIdx {
							continue Outputs
						}
					}
				}

				outs := UTXO[txID]
				outs.Outputs = append(outs.Outputs, out)
				UTXO[txID] = outs
			}

			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin {
					inTxID := hex.EncodeToString(in.Txid)
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return UTXO
}

// Iterator returns a BlockchainIterat
func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}

	return bci
}

// GetBestHeight returns the height of the latest block
func (bc *Blockchain) GetBestHeight() int {
	var lastBlock Block

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash := b.Get([]byte("l"))
		blockData := b.Get(lastHash)
		lastBlock = *DeserializeBlock(blockData)

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return lastBlock.Height
}

// GetBlockFromHash finds a block by its hash and returns it
func (bc *Blockchain) GetBlockFromHash(blockHash []byte) (Block, error) {
	var block Block

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		blockData := b.Get(blockHash)

		if blockData == nil {
			return errors.New("Block is not found.")
		}

		block = *DeserializeBlock(blockData)

		return nil
	})
	if err != nil {
		return block, err
	}

	return block, nil
}

// GetBlockFromHeight finds a block by its height and returns it
func (bc *Blockchain) GetBlockFromHeight(height int) (Block, error) {
	var block Block

	hashes := bc.GetBlockHashes()
	if height > len(hashes) - 1 {
		return block, errors.New("Block is not found.")
	}
	blockHash := hashes[len(hashes) - 1 - height]
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		blockData := b.Get(blockHash)

		if blockData == nil {
			return errors.New("Block is not found.")
		}

		block = *DeserializeBlock(blockData)

		return nil
	})
	if err != nil {
		return block, err
	}

	return block, nil
}

// GetBlockHashes returns a list of hashes of all the blocks in the chain
func (bc *Blockchain) GetBlockHashes() [][]byte {
	var blocks [][]byte
	bci := bc.Iterator()

	for {
		block := bci.Next()

		blocks = append(blocks, block.Hash)

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return blocks
}

// GetProlemGraphFromHash finds a Problemgraph by its hash and returns it
func (bc *Blockchain) GetProblemGraphFromHash(pgHash []byte) (ProblemGraph, error) {
	var pg ProblemGraph

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(problemsBucket))

		pgData := b.Get(pgHash)

		if pgData == nil {
			return errors.New("Problem is not found.")
		}

		pg = *DeserializeProblemGraph(pgData)

		return nil
	})
	if err != nil {
		return pg, err
	}

	return pg, nil
}

// GetProblemGraphHashes returns a list of hashes of all the problems in the chain
func (bc *Blockchain) GetProblemGraphHashes() [][]byte {
	var problems [][]byte
	bci := bc.Iterator()

	for {
		block := bci.Next()
		if len(block.ProblemGraphHash) > 0 {
			problems = append(problems, block.ProblemGraphHash)
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return problems
}

//Calculate the new target bits
func (bc *Blockchain) CalculateTarget(height int) *big.Int {
	var prevTarget *big.Int
	var newTarget *big.Int
	if height < blocksPerTargetUpdate {
		initialTarget := targetFromTargetBits(initialTargetBits)
		return initialTarget
	}

	hashes := bc.GetBlockHashes()
	total := len(hashes)
	index := ((height-1)/blocksPerTargetUpdate) * blocksPerTargetUpdate //this return only integer part of ratio since i'm divindg two integers
	baseBlock, _ := bc.GetBlockFromHash(hashes[total -1 - index])//this block is the first block in the batch of blocks we need to calculate difficulty

	if height%blocksPerTargetUpdate != 0 {
		return baseBlock.Target
	}
	
	lastBlock, _ := bc.GetBlockFromHash(hashes[total - 1 - (blocksPerTargetUpdate + index - 1)])//this block is the last block in the batch of blocks we need to calculate difficulty
	t := baseBlock.Timestamp - lastBlock.Timestamp
	prevTarget = baseBlock.Target
	timeTarget := secondsPerMinute * blocksPerTargetUpdate / targetBlocksPerMinute
	retarget := new(big.Float).SetFloat64(float64(t) / float64(timeTarget))
	floatTarget := new(big.Float).SetInt(prevTarget)
	floatNewTarget := new(big.Float).Mul(floatTarget, retarget)
	result := new(big.Int) 
	floatNewTarget.Int(result)
	newTarget = result
	maxChange := big.NewInt(maxTargetChange)
	maxTarget := new(big.Int).Mul(prevTarget, maxChange)
	minTarget := new(big.Int).Div(prevTarget, maxChange)
	if newTarget.Cmp(maxTarget) == 1 {
		newTarget = maxTarget
	} else if newTarget.Cmp(minTarget) == -1 {
		newTarget = minTarget
	}
	
	return newTarget

			
}

//Calculate the new target bits
func (bc *Blockchain) CurrentTarget() *big.Int {
	bci := bc.Iterator()
	block := bci.Next()
	height := block.Height + 1
	return bc.CalculateTarget(height)
}

// MineBlock mines a new block with the provided transactions
func (bc *Blockchain) MineBlock(transactions []*Transaction, solHash []byte, solution []int,  pgHash []byte) *Block {
	var lastHash []byte
	var lastHeight int

	for _, tx := range transactions {
		// TODO: ignore transaction if it's not valid
		if bc.VerifyTransaction(tx) != true {
			log.Panic("ERROR: Invalid transaction")
		}
	}
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		blockData := b.Get(lastHash)
		block := DeserializeBlock(blockData)

		lastHeight = block.Height

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	newTarget := bc.CalculateTarget(lastHeight+1)
	newBlock := NewBlock(transactions, lastHash, lastHeight+1, newTarget, solHash, solution, pgHash)
	pow := NewProofOfWork(newBlock)
	nonce, hash := pow.Run()

	newBlock.Hash = hash[:]
	newBlock.Nonce = nonce

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}

		bc.tip = newBlock.Hash

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return newBlock
}

// SignTransaction signs inputs of a Transaction
func (bc *Blockchain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	tx.Sign(privKey, prevTXs)
}

// VerifyTransaction verifies transaction input signatures
func (bc *Blockchain) VerifyTransaction(tx *Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	return tx.Verify(prevTXs)
}

// AddProblemGraph add a problem to the database
func (bc *Blockchain) AddProblemGraph(pg *ProblemGraph) {
	err := bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(problemsBucket))
		problemInDb := b.Get(pg.Hash)
		if problemInDb != nil {
			return nil
		}

		pgData := pg.Serialize()
		err := b.Put(pg.Hash, pgData)
		if err != nil {
			log.Panic(err)
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

func dbExists(dbFile string) bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}
