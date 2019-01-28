package cmd

import (
	"NaiveChain/core"
)

//create blockchain with genesis block
func (cli *CLI) createGenesisBlockChain(address string, nodeID string) {

	//create genesis block
	blc := core.CreateBlockchainWithAGenesisBlock(address, nodeID)
	defer blc.DB.Close()

	utxoSet := &core.UTXOSet{blc}
	utxoSet.Reset()
}
