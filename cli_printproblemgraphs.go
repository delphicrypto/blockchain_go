package crickchain

import (
	"fmt"
	"encoding/hex"
)

func (cli *CLI) printProblemGraphs(nodeID string) {
	bc := NewBlockchain(nodeID)
	defer bc.db.Close()

	hashes := bc.GetProblemGraphHashes()

	for i, h := range hashes {
		pg, err := bc.GetProblemGraphFromHash(h)
		if err == nil {
			fmt.Println("Problem ", i)
			pg.NicePrint(bc)			
		} else {
			fmt.Println(err)
		}
	}	
}

func (cli *CLI) printProblemGraph(nodeID string, hash string) {
	bc := NewBlockchain(nodeID)
	defer bc.db.Close()

	h, err := hex.DecodeString(hash)
	if err != nil {
	    panic(err)
	}
	pg, err := bc.GetProblemGraphFromHash(h)
	if err == nil {
		pg.NicePrint(bc)
		text := ProblemToString(pg)
		filename := "jsgraph/data/graph.js"
		WriteToFile(filename, text)
		textsol := "var cliques = ["

		allSolutions := bc.GetAllSolutions(&pg)
		for i, s := range allSolutions {
			textsol += CliqueToString(s)
			if i < len(allSolutions) - 1{
				textsol += ",\n "
			}
		}

		textsol += "];\n"
		filenamesol := "jsgraph/data/sol.js"
		WriteToFile(filenamesol, textsol)
	} else {
		fmt.Println(err)
	}
}	