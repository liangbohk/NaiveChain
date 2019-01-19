package cmd

import (
	"NaiveChain/core"
	"fmt"
)

func (cli *CLI) TestMethod() {
	blc := core.BlockchainObject()
	defer blc.DB.Close()

	utxoMap := blc.FindUTXOMap()
	fmt.Println(utxoMap)
}
