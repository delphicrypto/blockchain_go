package crickchain

import (
	"fmt"
	"log"
)

func (cli *CLI) createBlockchain(address, dbFile string) {
	if !ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}
	bc := CreateBlockchain(address, dbFile)
	defer bc.db.Close()

	UTXOSet := UTXOSet{bc}
	UTXOSet.Reindex()

	fmt.Println("Done!")
}
