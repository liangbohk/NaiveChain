package cmd

import "NaiveChain/core"

//create blockchain with genesis block
func (cli *CLI) createGenesisBlockChain(address string) {

	//create genesis block
	blc := core.CreateBlockchainWithAGenesisBlock(address)
	defer blc.DB.Close()
}
