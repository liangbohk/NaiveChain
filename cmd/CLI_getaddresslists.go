package cmd

import (
	"NaiveChain/core"
	"fmt"
)

//print all addresses
func (cli *CLI) getAddressList(nodeID string) {
	wallets, _ := core.NewWallets(nodeID)
	for address, _ := range wallets.WalletsMap {
		fmt.Println(address)
	}
}
