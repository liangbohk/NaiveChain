package main

import (
	"NaiveChain/core"
	"fmt"
)

func test() {
	wallet := core.NewWallet()
	address := wallet.GetAddress()
	fmt.Printf("%s\n", address)
	fmt.Println(len(address))
}

func main() {
	test()

	//cli := cmd.CLI{}
	//cli.Run()

}
