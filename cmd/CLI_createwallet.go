package cmd

import (
	"NaiveChain/core"
)

func (cli *CLI) createWallet() {
	wallets, _ := core.NewWallets()
	wallets.CreateNewWallet()
}
