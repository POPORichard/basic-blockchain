package CLI

import (
	"basic-blockchain/address"
	"basic-blockchain/database"
	"basic-blockchain/handel"
	"fmt"
	"log"
)

//打钱
func (cli *CLI) send(from, to string, amount int) {
	if !address.ValidateAddress(from) {
		log.Panic("ERROR: Sender address is not valid")
	}
	if !address.ValidateAddress(to) {
		log.Panic("ERROR: Recipient address is not valid")
	}

	bc := database.NewBlockchainLink()
	defer bc.Db.Close()
	UTXOSet := handel.UTXOSet{BlockChain: bc}

	//tx := handel.NewUTXOTransaction(from, to, amount, bc)
	//bc.MineBlock([]*handel.Transaction{tx})

	tx := handel.NewUTXOTransaction(from, to, amount, &UTXOSet)
	cbTx := handel.NewCoinbaseTX(from, "coinBase")
	txs := []*handel.Transaction{cbTx, tx}

	newBlock := bc.MineBlock(txs)
	UTXOSet.Update(newBlock)
	fmt.Println("Success!")
}