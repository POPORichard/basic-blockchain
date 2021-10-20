package CLI

import (
	add "basic-blockchain/address"
	"basic-blockchain/handel"
	"fmt"
	"log"
)

//打印余额
func (cli *CLI) getBalance(address string) {
	if !add.ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}
	bc := handel.NewBlockchainLink()
	defer bc.Db.Close()
	UTXOSet := handel.UTXOSet{BlockChain:bc}

	balance := 0
	pubKeyHash := add.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := UTXOSet.FindUTXO(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}
