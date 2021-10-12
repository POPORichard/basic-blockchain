package main

import (
	"basic-blockchain/handel"
	"fmt"
)

func main(){
	blockChain := handel.NewBlockChain()

	blockChain.AddBlock("message 1")

	for no,block := range blockChain.Blocks{
		fmt.Println("No:",no)
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()
	}
}
