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
	fmt.Println("\tgetaddresslist -- output all address")
	fmt.Println("\tcreatewallet -- create a wallet")
	fmt.Println("\tcreateblockchain -address ADDRESS -- create a blockchain")
	fmt.Println("\tsend -from FROM -to TO -amount AMOUNT -mine -- send value by transaction")
	fmt.Println("\tgetbalance -address ADDRESS -- get balance of an address")
	fmt.Println("\tprintchain -- print the block chain")
	fmt.Println("\tstartnode -mineaddress ADDRESS -- start sever node and specify mining address")
	fmt.Println("\ttest -- test")
}

//check if the args are valid
func isValidArg() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) Run() {
	isValidArg()

	//get NODE_ID
	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		fmt.Println("NODE_ID not set")
		os.Exit(1)
	}
	fmt.Printf("NODE_ID:%s\n", nodeID)

	getAddressListCmd := flag.NewFlagSet("getaddresslist", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	sendBlockCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printBlockchainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createBlockchainWithGeneisCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	startNodeCmd := flag.NewFlagSet("startnode", flag.ExitOnError)
	testCmd := flag.NewFlagSet("test", flag.ExitOnError)

	flagSendFrom := sendBlockCmd.String("from", "", "source address")
	flagSendTo := sendBlockCmd.String("to", "", "dist address")
	flagSendAmount := sendBlockCmd.String("amount", "", "transfer amount")
	flagSendMine := sendBlockCmd.Bool("mine", false, "if the transaction(s) are mine now")

	flagCreateBlockchainAddress := createBlockchainWithGeneisCmd.String("address", "", "genesis block address")

	flagGetBalanceWithAddress := getBalanceCmd.String("address", "", "look up the balance of an address")

	flagMineAddress := startNodeCmd.String("mineaddress", "", "coinbase address")

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
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "startnode":
		err := startNodeCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getaddresslist":
		err := getAddressListCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "test":
		err := testCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		printUsage()
		os.Exit(1)
	}

	if createBlockchainWithGeneisCmd.Parsed() {
		if core.IsValidAddress([]byte(*flagCreateBlockchainAddress)) == false {
			fmt.Println("invalid address")
			printUsage()
			os.Exit(11)
		}
		cli.createGenesisBlockChain(*flagCreateBlockchainAddress, nodeID)

	}

	if sendBlockCmd.Parsed() {
		if *flagSendFrom == "" || *flagSendTo == "" || *flagSendAmount == "" {
			printUsage()
			os.Exit(1)
		}
		from := core.Json2Array(*flagSendFrom)
		to := core.Json2Array(*flagSendTo)
		for index, _ := range from {
			if core.IsValidAddress([]byte(from[index])) == false || core.IsValidAddress([]byte(to[index])) == false {
				fmt.Println("invalid address")
				os.Exit(1)
			}
		}
		cli.send(from, to, core.Json2Array(*flagSendAmount), nodeID, *flagSendMine)
	}

	if getBalanceCmd.Parsed() {
		if *flagGetBalanceWithAddress == "" {
			printUsage()
			os.Exit(1)
		}
		cli.getBalance(*flagGetBalanceWithAddress, nodeID)
	}

	if createWalletCmd.Parsed() {
		cli.createWallet(nodeID)
	}
	if getAddressListCmd.Parsed() {
		cli.getAddressList(nodeID)
	}
	if startNodeCmd.Parsed() {
		//if core.IsValidAddress([]byte(*flagMineAddress)) == false {
		//	fmt.Println("invalid address")
		//	printUsage()
		//	os.Exit(11)
		//}
		cli.startNode(nodeID, *flagMineAddress)
	}

	if printBlockchainCmd.Parsed() {
		//fmt.Println("blockchain info")
		cli.printChain(nodeID)
	}

	if testCmd.Parsed() {
		cli.TestMethod(nodeID)
	}
}
