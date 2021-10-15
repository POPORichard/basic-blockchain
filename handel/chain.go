package handel

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

//// 向链上添加块
//func (blockChain *BlockChain) AddBlock(data string) {
//	prevBlock := blockChain.Blocks[len(blockChain.Blocks)-1]
//	newBlock := NewBlock(data, prevBlock.Hash)
//	blockChain.Blocks = append(blockChain.Blocks, newBlock)
//}

//使用提供的交易挖掘新块
func (bc *BlockChain) MineBlock(transactions []*Transaction) {
	var prevBlockHash []byte

	for _,tx := range transactions{
		if bc.VerifyTransaction(tx) != true{
			panic("ERROR: Invalid transaction")
		}
	}

	//从数据库中读取最后一个blockChain，并由此挖出下一个
	err := bc.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocksBucket"))
		prevBlockHash = b.Get([]byte("l"))

		return nil
	})
	if err != nil {
		fmt.Println("Error when view prevBlockHash in db err :", err)
	}

	//开始计算下一个
	newBlock := NewBlock(transactions, prevBlockHash)

	//建立与数据库的读写链接(Update)
	err = bc.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocksBucket"))
		//将新挖到的blockChain序列化后放入数据库
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			fmt.Println("Error when put new block! err :", err)
			panic(err)
		}
		//并将key(l)更新为新的block的hash
		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			fmt.Println("Error when put new l ! err :", err)
			panic(err)
		}
		//同时将key(l)更新为新的block的hash
		bc.Tip = newBlock.Hash

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

//// 创建新创世链
//func NewBlockChain() *BlockChain {
//	return &BlockChain{[]*Block{NewGenesisBlock()}}
//}

//签署input交易
func (bc *BlockChain)SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey){
	prevTXs := make(map[string]Transaction)

	for _,vin := range tx.Vin{
		prevTx, err := bc.FindTransaction(vin.Txid)
		if err != nil{
			panic(err)
		}
		prevTXs[hex.EncodeToString(prevTx.ID)] = prevTx
	}
	tx.Sign(privKey, prevTXs)
}

//验证input交易签名
func (bc *BlockChain) VerifyTransaction(tx *Transaction) bool{
	prevTXs := make(map[string]Transaction)

	for _,vin := range tx.Vin{
		prevTX, err := bc.FindTransaction(vin.Txid)
		if err != nil{
			panic(err)
		}

		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	return tx.Verify(prevTXs)
}

