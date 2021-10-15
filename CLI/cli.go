package CLI

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type CLI struct{}

//var addBlockCmd *flag.FlagSet
var printChainCmd *flag.FlagSet
var getBalanceCmd *flag.FlagSet
var createBlockchainCmd *flag.FlagSet
var sendCmd *flag.FlagSet
var createWalletCmd *flag.FlagSet
var listAddressesCmd *flag.FlagSet

func init() {
	//	addBlockCmd = flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd = flag.NewFlagSet("printchain", flag.ExitOnError)
	getBalanceCmd = flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd = flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd = flag.NewFlagSet("send", flag.ExitOnError)
	createWalletCmd = flag.NewFlagSet("createwallet", flag.ExitOnError)
	listAddressesCmd = flag.NewFlagSet("listaddresses", flag.ExitOnError)

}

func (cli *CLI) Run() {
	fmt.Println("start run!")
	//flag.Parse()
	//addBlockData := addBlockCmd.String("data", "", "Block data")
//	printNumber := printChainCmd.Int("number", 1, "print length. -1 means all")
	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")

	switch os.Args[1] {
	//case "addblock":
	//	err := addBlockCmd.Parse(os.Args[2:])
	//	if err != nil{
	//		fmt.Println(err)
	//	}
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
	default:
		flag.Usage()
		os.Exit(1)
	}

	//if addBlockCmd.Parsed(){
	//	fmt.Println("start addBlockCmd!")
	//	if *addBlockData == ""{
	//		addBlockCmd.Usage()
	//		os.Exit(1)
	//	}
	//	cli.AddBlock(*addBlockData)
	//}

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

}

//func (cli *CLI) AddBlock(data string) {
//	cli.Bc.AddBlock(data)
//	fmt.Println("Add message success!")
//}

//func (cli *CLI) PrintChain(num int) {
//	bc := database.NewBlockchainLink()
//	defer bc.Db.Close()
//	bci := bc.Iterator()
//	block := bci.Next()
//	if len(block.PrevBlockHash) == 0 {
//		fmt.Println("only one block in chain")
//		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
//		//fmt.Printf("Data: %s\n", block.Transaction)
//		fmt.Printf("Hash: %x\n", block.Hash)
//		pow := handel.NewProofOfWork(block)
//		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
//		fmt.Printf("Nonce: %x\n", block.Nonce)
//		fmt.Println("Transaction :")
//		for _, t := range block.Transaction {
//			fmt.Println("    Vout:", t.VOut, "<-----", "Vin:", t.Vin)
//		}
//		fmt.Println()
//	} else {
//		all := false
//		if num == -1 {
//			all = true
//		}
//		for num > 0 || all {
//
//			fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
//			//fmt.Printf("Data: %s\n", block.Transaction)
//			fmt.Printf("Hash: %x\n", block.Hash)
//			pow := handel.NewProofOfWork(block)
//			fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
//			fmt.Printf("Nonce: %x\n", block.Nonce)
//			for _, t := range block.Transaction {
//				fmt.Println("    Vout:", t.VOut, "<-----", "Vin:", t.Vin)
//			}
//			fmt.Println()
//
//			if len(block.PrevBlockHash) == 0 {
//				break
//			}
//
//			block = bci.Next()
//			num--
//		}
//	}
//
//}

//func (cli *CLI) send(from, to string, amount int) {
//	bc := database.NewBlockchainLink()
//	defer bc.Db.Close()
//
//	tx := handel.NewUTXOTransaction(from, to, amount, bc)
//	bc.MineBlock([]*handel.Transaction{tx})
//	fmt.Println("Send Success!")
//
//}

//func (cli *CLI) getBalance(address string) {
//	bc := database.NewBlockchainLink()
//	defer bc.Db.Close()
//
//	balance := 0
//	UTXOs := bc.FindUTXO(address)
//
//	for _, out := range UTXOs {
//		balance += out.Value
//	}
//
//	fmt.Printf("Balance of '%s': %d\n", address, balance)
//}

//func (cli *CLI) createBlockchain(address string) {
//	fmt.Println("The first coin will send to", address)
//	bc := database.CreateBlockchain(address)
//	bc.Db.Close()
//	fmt.Println("Done!")
//}

func (cli *CLI)printUsage(){
	fmt.Println("Usage:")
	fmt.Println("  createblockchain -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println("  createwallet - Generates a new key-pair and saves it into the wallet file")
	fmt.Println("  getbalance -address ADDRESS - Get balance of ADDRESS")
	fmt.Println("  listaddresses - Lists all addresses from the wallet file")
	fmt.Println("  printchain - Print all the blocks of the blockchain")
	fmt.Println("  send -from FROM -to TO -amount AMOUNT - Send AMOUNT of coins from FROM address to TO")
}
