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
	//get txs with unspent output
	balance := blockchain.GetBanlance(address)

	fmt.Printf("balance %d\n", balance)
}
