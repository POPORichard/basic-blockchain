package CLI

import (
	"basic-blockchain/address"
	"fmt"
)
//创建钱包并打印公钥给用户
func (cli *CLI) createWallet() {
	wallets, _ := address.NewWallets()
	address := wallets.CreateWallet()
	wallets.SaveToFile()

	fmt.Printf("Your new address: %s\n", address)
}
