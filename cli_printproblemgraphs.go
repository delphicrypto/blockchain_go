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
			text := ProblemToString(pg)
			filename := "jsgraph/data/graph.js"
			WriteToFile(filename, text)
			if len(pg.Graph.AdjacencyList) < 50 {
				maxCliques := pg.FindMaxClique()
				textsol := CliquesToString(maxCliques)
				filenamesol := "jsgraph/data/sol.js"
				WriteToFile(filenamesol, textsol)	
			} 
		} else {
			fmt.Println(err)
		}
	}

	
}

