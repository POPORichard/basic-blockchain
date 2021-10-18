package CLI

import (
	"basic-blockchain/handel"
	"bytes"
	"encoding/gob"
	"fmt"
	"strconv"
)

//打印区块连
func (cli *CLI) printChain(nodeID string) {
	bc := handel.NewBlockchainLink(nodeID)
	defer bc.Db.Close()

	bci := bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("============ Block %x ============\n", block.Hash)
		fmt.Printf("Height: %d\n", block.Height)
		fmt.Printf("Prev block: %x\n", block.PrevBlockHash)
		pow := handel.NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		//TODO : show Show transaction details
		//for _, tx := range block.Transaction {
		//	fmt.Println("Transaction:")
		//	for i,v := range tx.Vin{
		//		fmt.Println(i, "--- VOut: ",v.Vout, "===>", Deserialize(v.PubKey))
		//	}
		//	for i,v := range tx.VOut{
		//		fmt.Println(i, "--- VOut: ",v.Value, "===>", address.Base58Decode(v.PubKeyHash))
		//	}
		//
		//}
		fmt.Printf("\n\n")

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

func Deserialize(b []byte) *handel.Transaction {
	var tx handel.Transaction
	decoder := gob.NewDecoder(bytes.NewReader(b))
	err := decoder.Decode(&tx)
	if err != nil {
		fmt.Println("Error in Deserialize err: ", err)
		return nil
	}
	return &tx
}