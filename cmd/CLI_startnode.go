package cmd

import (
	"NaiveChain/core"
	"fmt"
	"os"
)

func (cli *CLI) startNode(nodeID string, mineAddress string) {
	//start server
	if mineAddress == "" || core.IsValidAddress([]byte(mineAddress)) {
		fmt.Printf("Start server localhost:%s\n", nodeID)
		core.StartServer(nodeID, mineAddress)
	} else {
		fmt.Println("Invalid address")
		os.Exit(1)
	}
}
