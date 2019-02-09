package crickchain


func (cli *CLI) printHeight(nodeID string, height int) {
	bc := NewBlockchain(nodeID)
	defer bc.db.Close()

	block, _ := bc.GetBlockFromHeight(height)
	block.NicePrint(bc)

}


