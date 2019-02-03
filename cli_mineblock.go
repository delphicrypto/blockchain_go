package main

import (
	"fmt"
)

func (cli *CLI) mineblock(nodeID string, withProblemGraph bool) {
	bc := NewBlockchain(nodeID)
	defer bc.db.Close()

	pgHash := []byte{}
	if withProblemGraph {
		pg := NewProblemGraph(500, 800)
		bc.AddProblemGraph(pg)
		pgHash = pg.Hash
	}
	var txs []*Transaction
	newBlock := bc.MineBlock(txs, []byte{}, []int{}, pgHash)
	fmt.Printf("New block hash: %x\r\n", newBlock.Hash)
}
