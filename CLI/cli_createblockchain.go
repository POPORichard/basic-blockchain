package CLI

import (
	add "basic-blockchain/address"
	"basic-blockchain/handel"
	"fmt"
	"io/ioutil"
	"log"
)
//创建创世chain并将第一个block的奖励发给address
func (cli *CLI) createBlockchain(address string) {
	if !add.ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}
	bc := handel.CreateBlockchain(address)
	defer bc.Db.Close()

	UTXOSet := handel.UTXOSet{BlockChain:bc}
	UTXOSet.Reindex()

	fmt.Println("Done!")

	input, err := ioutil.ReadFile(handel.DbFile)
	err = ioutil.WriteFile(handel.DbFile+"_genesis", input, 0644)
	if err != nil{
		fmt.Println("can not create blockChain_genesis")
	}
}
