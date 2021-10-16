package CLI

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type CLI struct{}

var printChainCmd *flag.FlagSet
var getBalanceCmd *flag.FlagSet
var createBlockchainCmd *flag.FlagSet
var sendCmd *flag.FlagSet
var createWalletCmd *flag.FlagSet
var listAddressesCmd *flag.FlagSet
var reindexUTXOCmd *flag.FlagSet

func init() {
	printChainCmd = flag.NewFlagSet("printchain", flag.ExitOnError)
	getBalanceCmd = flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd = flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd = flag.NewFlagSet("send", flag.ExitOnError)
	createWalletCmd = flag.NewFlagSet("createwallet", flag.ExitOnError)
	listAddressesCmd = flag.NewFlagSet("listaddresses", flag.ExitOnError)
	reindexUTXOCmd = flag.NewFlagSet("reindexutxo", flag.ExitOnError)

}

func (cli *CLI) Run() {
	fmt.Println("start run!")
	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")

	switch os.Args[1] {
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(err)
		}
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "listaddresses":
		err := listAddressesCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "reindexutxo":
		err := reindexUTXOCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		flag.Usage()
		os.Exit(1)
	}

	if printChainCmd.Parsed() {
		fmt.Println("start printChainCmd!")
		cli.printChain()
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceAddress)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createBlockchainAddress)
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}

		cli.send(*sendFrom, *sendTo, *sendAmount)
	}
	if createWalletCmd.Parsed() {
		cli.createWallet()
	}

	if listAddressesCmd.Parsed() {
		cli.listAddresses()
	}
	if reindexUTXOCmd.Parsed() {
		cli.reindexUTXO()
	}

}


func (cli *CLI)printUsage(){
	fmt.Println("Usage:")
	fmt.Println("  createblockchain -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println("  createwallet - Generates a new key-pair and saves it into the wallet file")
	fmt.Println("  getbalance -address ADDRESS - Get balance of ADDRESS")
	fmt.Println("  listaddresses - Lists all addresses from the wallet file")
	fmt.Println("  printchain - Print all the blocks of the blockchain")
	fmt.Println("  send -from FROM -to TO -amount AMOUNT - Send AMOUNT of coins from FROM address to TO")
	fmt.Println("  reindexutxo - Rebuilds the UTXO set")
}
