package cmd

import (
	"NaiveChain/core"
)

//print the blockchain
func (cli *CLI) printChain(nodeID string) {

	blc := core.BlockchainObject(nodeID)
	defer blc.DB.Close()
	blc.PrintChain()
}
