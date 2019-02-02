package main

import (
	"fmt"
)

func (cli *CLI) mineblock(nodeID string) {
	bc := NewBlockchain(nodeID)
	defer bc.db.Close()
	var txs []*Transaction
	newBlock := bc.MineBlock(txs)
	fmt.Printf("New block hash: %x\r\n", newBlock.Hash)
}
