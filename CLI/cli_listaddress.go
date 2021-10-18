package CLI


import (
	"basic-blockchain/address"
	"fmt"
	"log"
)

//打印钱包的所有地址
func (cli *CLI) listAddresses(nodeID string) {
	wallets, err := address.NewWallets(nodeID)
	if err != nil {
		log.Panic(err)
	}
	addresses := wallets.GetAddresses()

	for _, address := range addresses {
		fmt.Println(address)
	}
}