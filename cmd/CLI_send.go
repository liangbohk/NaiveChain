package cmd

import (
	"NaiveChain/core"
	"fmt"
)

//send transaction
func (cli *CLI) send(from []string, to []string, amount []string, nodeID string, mineNow bool) {

	blc := core.BlockchainObject(nodeID)
	defer blc.DB.Close()

	if mineNow {
		blc.MineNewBlock(from, to, amount, nodeID)

		utxoSet := &core.UTXOSet{blc}
		utxoSet.Update()
	} else {
		//send transaction to miner node
		fmt.Println("not mined now")
	}

}
