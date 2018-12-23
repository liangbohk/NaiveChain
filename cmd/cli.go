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
	fmt.Println("\tcreateblockchain -address ADDRESS -- address")
	fmt.Println("\tsend -from FROM -to TO -amount AMOUNT -- transaction")
	fmt.Println("\tgetbalance -address ADDRESS -- get balance of an address")
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
//func (cli *CLI) addBlock(txs []*core.Transaction) {
//	if !core.DBExist() {
//		log.Fatal("no blockchain")
//	}
//	blc := core.BlockchainObject()
//	defer blc.DB.Close()
//	blc.AddBlockToBlockchain([]*core.Transaction{})
//}

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
func (cli *CLI) createGenesisBlockChain(address string) {

	//create genesis block
	blc := core.CreateBlockchainWithAGenesisBlock(address)
	defer blc.DB.Close()
}

func (cli *CLI) send(from []string, to []string, amount []string) {
	if !core.DBExist() {
		log.Fatal("no blockchain")
	}
	blc := core.BlockchainObject()
	defer blc.DB.Close()

	blc.MineNewBlock(from, to, amount)
}

//look up balance
func (cli *CLI) getBalance(address string) {
	fmt.Printf("address: %s\n", address)

	blockchain := core.BlockchainObject()
	defer blockchain.DB.Close()
	//get txs with unspent output
	balance := blockchain.GetBanlance(address)

	fmt.Printf("balance %d\n", balance)
}

func (cli *CLI) Run() {
	isValidArg()

	sendBlockCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printBlockchainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createBlockchainWithGeneisCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)

	flagSendFrom := sendBlockCmd.String("from", "", "source address")
	flagSendTo := sendBlockCmd.String("to", "", "dist address")
	flagSendAmount := sendBlockCmd.String("amount", "", "transfer amount")

	flagCreateBlockchainAddress := createBlockchainWithGeneisCmd.String("address", "", "genesis block address")

	flagGetBalanceWithAddress := getBalanceCmd.String("address", "", "look up the balance of an address")

	switch os.Args[1] {
	case "send":
		err := sendBlockCmd.Parse(os.Args[2:])
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
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		printUsage()
		os.Exit(1)
	}

	if createBlockchainWithGeneisCmd.Parsed() {
		if *flagCreateBlockchainAddress == "" {
			fmt.Println("address cannot be empty ")
			printUsage()
			os.Exit(11)
		}
		cli.createGenesisBlockChain(*flagCreateBlockchainAddress)

	}

	if sendBlockCmd.Parsed() {
		if *flagSendFrom == "" || *flagSendTo == "" || *flagSendAmount == "" {
			printUsage()
			os.Exit(1)
		}
		//fmt.Println(*flagAddBlockData)
		//cli.addBlock([]*core.Transaction{})
		//fmt.Println(*flagSendFrom)
		//fmt.Println(*flagSendTo)
		//fmt.Println(*flagSendAmount)
		//
		//fmt.Println(core.Json2Array(*flagSendFrom))
		//fmt.Println(core.Json2Array(*flagSendTo))
		//fmt.Println(core.Json2Array(*flagSendAmount))
		cli.send(core.Json2Array(*flagSendFrom), core.Json2Array(*flagSendTo), core.Json2Array(*flagSendAmount))
	}

	if getBalanceCmd.Parsed() {
		if *flagGetBalanceWithAddress == "" {
			printUsage()
			os.Exit(1)
		}
		cli.getBalance(*flagGetBalanceWithAddress)
	}

	if printBlockchainCmd.Parsed() {
		//fmt.Println("blockchain info")
		cli.printChain()
	}
}
