package cmd

import (
	"NaiveChain/core"
)

func (cli *CLI) TestMethod() {
	blc := core.BlockchainObject()
	defer blc.DB.Close()

	utxoSet := &core.UTXOSet{blc}
	utxoSet.Reset()
	//utxoMap := blc.FindUTXOMap()
	//fmt.Println(utxoMap)
}
