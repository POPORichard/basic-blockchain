package CLI

import (
	"basic-blockchain/database"
	"basic-blockchain/handel"
	"fmt"
	"strconv"
)

func (cli *CLI) printChain() {
	bc := database.NewBlockchainLink()
	defer bc.Db.Close()

	bci := bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("============ Block %x ============\n", block.Hash)
		fmt.Printf("Prev. block: %x\n", block.PrevBlockHash)
		pow := handel.NewProofOfWork(block)
		fmt.Printf("PoW: %s\n\n", strconv.FormatBool(pow.Validate()))
		for _, tx := range block.Transaction {
			fmt.Println(tx)
		}
		fmt.Printf("\n\n")

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}