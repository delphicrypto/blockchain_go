package crickchain


func (cli *CLI) createGraph(nodeID string, nodes int, edges int) {
	bc := NewBlockchain(nodeID)
	defer bc.db.Close()

	pg := NewProblemGraph(nodes, edges)
	pg.NicePrint(bc)
	bc.AddProblemGraph(pg)
}
