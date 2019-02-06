package main

import (
	"fmt"
	"encoding/hex"
	"log"
)

func (cli *CLI) mineblock(nodeID string) {
	bc := NewBlockchain(nodeID)
	defer bc.db.Close()
	var txs []*Transaction
	
	newBlock := bc.MineBlock(txs, []byte{}, []int{}, []byte{})
	fmt.Printf("New block hash: %x\r\n", newBlock.Hash)
}

func (cli *CLI) mineblockWithNewProblem(nodeID string, nodes int, density float64) {
	bc := NewBlockchain(nodeID)
	defer bc.db.Close()

	kclique := []int{}
	edges := int(float64(nodes*(nodes-1)/2) * density)
	pg := NewProblemGraph(nodes, edges)
	bc.AddProblemGraph(pg)
	//we mine the problem with an initial solution
	for k := 8; k >= 3; k-- {
		kclique = pg.FindKClique(k)
		if len(kclique) == k {
			break
		}
	}
	
	var txs []*Transaction
	
	
	newBlock := bc.MineBlock(txs, pg.Hash, kclique, pg.Hash)
	fmt.Printf("New block hash: %x\r\n", newBlock.Hash)
}

func (cli *CLI) mineblockWithSolution(nodeID string, pgHash string) {
	bc := NewBlockchain(nodeID)
	defer bc.db.Close()
	
	hash, err := hex.DecodeString(pgHash)
	if err != nil {
		fmt.Println("Invalid problem hash")
		log.Panic(err)
	}
	pg, err := bc.GetProblemGraphFromHash(hash)
	if err != nil {
		fmt.Println("Invalid problem hash")
		log.Panic(err)
	}

	height := bc.GetBestHeight()
	bestSolution := bc.GetBestSolution(&pg, height)
	kclique := pg.FindKClique(len(bestSolution) + 1)

	var txs []*Transaction

	newBlock := bc.MineBlock(txs, hash, kclique, []byte{})
	fmt.Printf("New block hash: %x\r\n", newBlock.Hash)
}
