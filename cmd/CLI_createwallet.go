package cmd

import (
	"NaiveChain/core"
	"fmt"
)

func (cli *CLI) createWallet() {
	wallets, _ := core.NewWallets()
	wallets.CreateNewWallet()

	fmt.Println(wallets.WalletsMap)
}
