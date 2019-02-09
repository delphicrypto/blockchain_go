package crickchain


func (cli *CLI) printChain(nodeID string) {
	bc := NewBlockchain(nodeID)
	defer bc.db.Close()

	bci := bc.Iterator()
	for {
		block := bci.Next()
		block.NicePrint(bc)

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
