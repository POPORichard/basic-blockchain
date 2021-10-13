package handel

import (
	"fmt"
	"github.com/boltdb/bolt"
)

//// 向链上添加块
//func (blockChain *BlockChain) AddBlock(data string) {
//	prevBlock := blockChain.Blocks[len(blockChain.Blocks)-1]
//	newBlock := NewBlock(data, prevBlock.Hash)
//	blockChain.Blocks = append(blockChain.Blocks, newBlock)
//}

//向链上添加块
func (bc *BlockChain) AddBlock (data string) {
	var prevBlockHash []byte

	err := bc.Db.View(func (tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocksBucket"))
		prevBlockHash = b.Get([]byte("l"))

		return nil
	})
	if err != nil{
		fmt.Println("Error when view prevBlockHash in db err :", err)
	}

	newBlock := NewBlock(data, prevBlockHash)

	err = bc.Db.Update(func (tx *bolt.Tx) error{
		b := tx.Bucket([]byte("blocksBucket"))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil{
			fmt.Println("Error when put new block! err :", err)
		}

		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil{
			fmt.Println("Error when put new l ! err :",err)
		}
		bc.Tip = newBlock.Hash

		return nil
	})
}

//// 创建新创世链
//func NewBlockChain() *BlockChain {
//	return &BlockChain{[]*Block{NewGenesisBlock()}}
//}





