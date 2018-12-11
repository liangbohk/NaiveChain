package cmd

import (
	"NaiveChain/core"
	"flag"
	"fmt"
	"log"
	"os"
)

//define the
type CLI struct {
	//Blc *core.Blockchain
}

//print command usage
func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("\tcreateblockchain -data DATA -- transaction data")
	fmt.Println("\taddblock -data DATA -- transaction data")
	fmt.Println("\tprintchain -- print the block chain")
}

//check if the args are valid
func isValidArg() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}

//add a block to the blockchain
func (cli *CLI) addBlock(data string) {
	if !core.DBExist() {
		log.Fatal("no blockchain")
	}
	blc := core.BlockchainObject()
	defer blc.DB.Close()
	blc.AddBlockToBlockchain(data)
}

//print the blockchain
func (cli *CLI) printChain() {
	if !core.DBExist() {
		log.Fatal("no blockchain")
	}
	blc := core.BlockchainObject()
	defer blc.DB.Close()
	blc.PrintChain()
}

//create blockchain with genesis block
func (cli *CLI) createGenesisBlockChain(data string) {
	fmt.Printf("Genesis data: %s\n", data)
	core.CreateBlockchainWithAGenesisBlock(data)
}

func (cli *CLI) Run() {
	isValidArg()

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printBlockchainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createBlockchainWithGeneisCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)

	flagAddBlockData := addBlockCmd.String("data", "added block data", "transaction data")
	flagCreateBlockchainData := createBlockchainWithGeneisCmd.String("data", "Genesis block", "genesis block transaction data")

	switch os.Args[1] {
	case "addblock":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainWithGeneisCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		printUsage()
		os.Exit(1)
	}

	if createBlockchainWithGeneisCmd.Parsed() {
		if *flagCreateBlockchainData == "" {
			fmt.Println("empty data")
			printUsage()
			os.Exit(11)
		}
		cli.createGenesisBlockChain(*flagCreateBlockchainData)

	}

	if addBlockCmd.Parsed() {
		if *flagAddBlockData == "" {
			printUsage()
			os.Exit(1)
		}
		//fmt.Println(*flagAddBlockData)
		cli.addBlock(*flagAddBlockData)
	}

	if printBlockchainCmd.Parsed() {
		//fmt.Println("blockchain info")
		cli.printChain()
	}
}
