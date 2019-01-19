package cmd

import (
	"NaiveChain/core"
	"fmt"
)

//look up balance
func (cli *CLI) getBalance(address string) {
	fmt.Printf("address: %s\n", address)

	blockchain := core.BlockchainObject()
	defer blockchain.DB.Close()

	//get balance from utxo
	utxoSet := &core.UTXOSet{blockchain}
	amount := utxoSet.GetBalance(address)

	//get txs with unspent output
	//balance := blockchain.GetBanlance(address)

	fmt.Printf("balance %d\n", amount)
}
