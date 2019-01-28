package cmd

import (
	"NaiveChain/core"
)

func (cli *CLI) createWallet(nodeID string) {
	wallets, _ := core.NewWallets(nodeID)
	wallets.CreateNewWallet(nodeID)
}
