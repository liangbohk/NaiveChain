package main

import (
	"NaiveChain/cmd"
	"NaiveChain/core"
	"fmt"
)

func test() {
	//wallet := core.NewWallet()
	//address := wallet.GetAddress()
	//
	//fmt.Printf("%s\n", address)
	//fmt.Println(string(address))
	//fmt.Println(address)
	//fmt.Printf("%s",hex.EncodeToString(address))
	//fmt.Println(wallet.IsValidAddress(address))

	wallets := core.NewWallets()
	wallets.CreateNewWallet()
	wallets.CreateNewWallet()
	fmt.Println(wallets.Wallets)
}

func main() {
	//test()

	cli := cmd.CLI{}
	cli.Run()

}
