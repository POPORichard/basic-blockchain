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

	tx := handel.NewUTXOTransaction(from, to, amount, bc)
	bc.MineBlock([]*handel.Transaction{tx})
	fmt.Println("Success!")
}