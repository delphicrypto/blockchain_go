package main

import (
	"fmt"
)

func (cli *CLI) mineblock(nodeID string, withProblemGraph bool, nodes int, density float64) {
	bc := NewBlockchain(nodeID)
	defer bc.db.Close()

	pgHash := []byte{}
	kclique := []int{}
	if withProblemGraph {
		edges := int(float64(nodes*(nodes-1)/2) * density)
		pg := NewProblemGraph(nodes, edges)
		bc.AddProblemGraph(pg)
		pgHash = pg.Hash
		//we mine the problem with an initial solution
		for k := 8; k >= 3; k-- {
			kclique = pg.FindKClique(k)
			if len(kclique) == k {
				break
			}
		}
	}
	var txs []*Transaction
	
	
	newBlock := bc.MineBlock(txs, pgHash, kclique, pgHash)
	fmt.Printf("New block hash: %x\r\n", newBlock.Hash)
}
