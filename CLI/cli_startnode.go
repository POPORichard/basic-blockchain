package CLI

import (
	"basic-blockchain/address"
	"basic-blockchain/server"
	"fmt"
)

func (cli *CLI)startNode(minerAddress string){
	fmt.Printf("Starting...\n")
	if len(minerAddress) > 0{
		if address.ValidateAddress(minerAddress){
			fmt.Println("Mining is on.Address to receive rewards:", minerAddress)
		} else{
			panic("wrong miner address!")
		}
	}
	server.StartServer(minerAddress)
}
