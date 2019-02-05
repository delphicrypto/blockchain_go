package main

import (
	"fmt"
	"strconv"
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
		maxK := 6
		textsol := "var cliques = {"
		for k := 3; k <= maxK; k++ {
			kCliques := pg.FindAllKCliques(k)
			textsol += CliquesToString(strconv.Itoa(k), kCliques)
			if k < maxK {
				textsol += ",\n "
			}
		}
		textsol += "};\n"
		filenamesol := "jsgraph/data/sol.js"
		WriteToFile(filenamesol, textsol)
	} else {
		fmt.Println(err)
	}
}	