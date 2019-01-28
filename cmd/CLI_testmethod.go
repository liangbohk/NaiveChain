package cmd

import (
	"NaiveChain/core"
)

func (cli *CLI) TestMethod(nodeID string) {
	blc := core.BlockchainObject(nodeID)
	defer blc.DB.Close()

	utxoSet := &core.UTXOSet{blc}
	utxoSet.Reset()
	//utxoMap := blc.FindUTXOMap()
	//fmt.Println(utxoMap)
}
