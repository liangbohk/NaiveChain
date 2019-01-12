package cmd

import (
	"NaiveChain/core"
	"fmt"
)

//print all addresses
func (cli *CLI) getAddressList() {
	wallets, _ := core.NewWallets()
	for address, _ := range wallets.WalletsMap {
		fmt.Println(address)
	}
}
