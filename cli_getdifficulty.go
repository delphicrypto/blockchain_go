package main

import (
	"fmt"
)

func (cli *CLI) getDifficulty(nodeID string) {
	bc := NewBlockchain(nodeID)
	defer bc.db.Close()
	target := bc.CurrentTarget()
	diff := targetToDifficulty(target)
	fmt.Printf("Target: %d\r\n", target)
	fmt.Printf("Difficulty: %d\r\n", diff)
}
