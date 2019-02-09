package crickchain


func (cli *CLI) printLast(nodeID string) {
	bc := NewBlockchain(nodeID)
	defer bc.db.Close()

	bci := bc.Iterator()
	block := bci.Next()
	block.NicePrint(bc)

}


