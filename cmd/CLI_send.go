package cmd

import (
	"NaiveChain/core"
	"log"
)

//send transaction
func (cli *CLI) send(from []string, to []string, amount []string) {
	if !core.DBExist() {
		log.Fatal("no blockchain")
	}
	blc := core.BlockchainObject()
	defer blc.DB.Close()

	blc.MineNewBlock(from, to, amount)

	utxoSet := &core.UTXOSet{blc}
	utxoSet.Update()
}
