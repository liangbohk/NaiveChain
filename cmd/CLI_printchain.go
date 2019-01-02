package cmd

import (
	"NaiveChain/core"
	"log"
)

//print the blockchain
func (cli *CLI) printChain() {
	if !core.DBExist() {
		log.Fatal("no blockchain")
	}
	blc := core.BlockchainObject()
	defer blc.DB.Close()
	blc.PrintChain()
}
