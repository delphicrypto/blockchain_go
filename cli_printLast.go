package main

import (
	"fmt"
	"strconv"
)

func (cli *CLI) printLast(nodeID string) {
	bc := NewBlockchain(nodeID)
	defer bc.db.Close()

	bci := bc.Iterator()
	block := bci.Next()
	printGreen(fmt.Sprintf("============ Block %d ============\n", block.Height))
	printBlue(fmt.Sprintf("Hash: %x\n", block.Hash))
	fmt.Printf("Prev: %x\n", block.PrevBlockHash)
	fmt.Printf("Block target: %d\n", block.Target)
	blockchainTarget := bc.CalculateTarget(block.Height)
	fmt.Printf("Chain target: %d\n", blockchainTarget)
	fmt.Printf("Difficulty: %d\n", targetToDifficulty(block.Target))
	validBlock := block.Validate(blockchainTarget)
	if validBlock {
		printGreen(fmt.Sprintf("PoW: %s\n\n", strconv.FormatBool(validBlock)))
	} else {
		printRed(fmt.Sprintf("PoW: %s\n\n", strconv.FormatBool(validBlock)))
	}
	for _, tx := range block.Transactions {
		printYellow(fmt.Sprintln(tx))
	}
	fmt.Printf("\n\n")
}
