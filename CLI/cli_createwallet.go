package CLI

import (
	"basic-blockchain/address"
	"fmt"
)

func (cli *CLI) createWallet() {
	wallets, _ := address.NewWallets()
	address := wallets.CreateWallet()
	wallets.SaveToFile()

	fmt.Printf("Your new address: %s\n", address)
}
