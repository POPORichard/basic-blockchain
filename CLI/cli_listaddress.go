package CLI


import (
	"basic-blockchain/address"
	"fmt"
	"log"
)

//打印钱包的所有地址
func (cli *CLI) listAddresses() {
	wallets, err := address.NewWallets()
	if err != nil {
		log.Panic(err)
	}
	addresses := wallets.GetAddresses()

	for _, address := range addresses {
		fmt.Println(address)
	}
}