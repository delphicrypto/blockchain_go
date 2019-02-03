package main

import (
	"fmt"
	// "log"
	"strconv"
)

func (cli *CLI) createGraph(nodes int, edges int) {
	pg := NewProblemGraph(nodes, edges)
	printBlue(fmt.Sprintf("Hash: %x\n",pg.Hash))
	for fr, to := range pg.Graph.AdjacencyList {
    	fmt.Println(fr, to)
	}
	clique := pg.FindClique()
	printGreen(fmt.Sprintf("Clique: %s\n", strconv.FormatBool(clique)))
	connected := pg.Graph.IsConnected()
	printGreen(fmt.Sprintf("Connected: %s\n", strconv.FormatBool(connected)))
}
