package main

import (
	"NaiveChain/core"
	"log"
)

func main() {
	blc := core.CreateBlockchainWithAGenesisBlock()
	blc.AddBlockToBlockchain("first")
	blc.AddBlockToBlockchain("second")
	blc.AddBlockToBlockchain("third")
	blc.AddBlockToBlockchain("forth")

	err := blc.DB.Close()
	if err != nil {
		log.Panic(err)
	}

}
