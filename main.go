package main

import (
	"NaiveChain/core"
	"fmt"
)

func main() {
	blc := core.CreateBlockchainWithAGenesisBlock()
	blc.AddBlockToBlockchain("Second block", blc.Blocks[len(blc.Blocks)-1].Height+1, blc.Blocks[len(blc.Blocks)-1].Hash)
	fmt.Println(*blc.Blocks[1])

}
