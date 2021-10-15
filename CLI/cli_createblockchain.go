package CLI

import (

	add "basic-blockchain/address"
	"basic-blockchain/database"
	"fmt"
	"log"
)

func (cli *CLI) createBlockchain(address string) {
	if !add.ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}
	bc := database.CreateBlockchain(address)
	bc.Db.Close()
	fmt.Println("Done!")
}
