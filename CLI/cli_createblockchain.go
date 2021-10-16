package CLI

import (

	add "basic-blockchain/address"
	"basic-blockchain/database"
	"basic-blockchain/handel"
	"fmt"
	"log"
)
//创建创世chain并将第一个block的奖励发给address
func (cli *CLI) createBlockchain(address string) {
	if !add.ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}
	bc := database.CreateBlockchain(address)
	defer bc.Db.Close()

	UTXOSet := handel.UTXOSet{BlockChain:bc}
	UTXOSet.Reindex()

	fmt.Println("Done!")
}
