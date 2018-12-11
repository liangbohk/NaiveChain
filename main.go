package main

import (
	"NaiveChain/cmd"
)

func main() {

	//blc := core.CreateBlockchainWithAGenesisBlock("genesis block")

	cli := cmd.CLI{}
	cli.Run()

}
