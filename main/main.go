package main

import (
	"basic-blockchain/database"
)

func main(){
	blockChain := database.NewBlockChain()
	defer blockChain.Db.Close()


	blockChain.AddBlock("message 1")

	//for no,block := range blockChain.Blocks{
	//	fmt.Println("No:",no)
	//	fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
	//	fmt.Printf("Data: %s\n", block.Data)
	//	fmt.Printf("Hash: %x\n", block.Hash)
	//	pow := handel.NewProofOfWork(block)
	//	fmt.Printf("nonce: %x\n", block.Nonce)
	//	fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
	//	fmt.Println()
	//}
}
