package CLI

import (
	"basic-blockchain/address"
	"fmt"
)
//创建钱包并打印公钥给用户
func (cli *CLI) createWallet(nodeID string) {
	wallets, _ := address.NewWallets(nodeID)
	address := wallets.CreateWallet()
	wallets.SaveToFile(nodeID)

	fmt.Printf("Your new address: %s\n", address)
}
