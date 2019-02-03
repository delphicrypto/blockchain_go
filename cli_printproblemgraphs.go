package main

import (
	"fmt"
)

func (cli *CLI) printProblemGraphs(nodeID string) {
	bc := NewBlockchain(nodeID)
	defer bc.db.Close()

	hashes := bc.GetProblemGraphHashes()

	for i, h := range hashes {
		pg, err := bc.GetProblemGraphFromHash(h)
		if err == nil {
			fmt.Println("Problem ", i)
			pg.NicePrint()
		} else {
			fmt.Println(err)
		}
	}
}
