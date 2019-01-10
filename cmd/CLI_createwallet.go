package cmd

import (
	"NaiveChain/core"
	"fmt"
)

func (cli *CLI) createWallet() {
	wallets := core.NewWallets()
	wallets.CreateNewWallet()
	fmt.Println(wallets.Wallets)
}
