package handel

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

type BlockChain struct {
	Tip []byte
	Db  *bolt.DB
}

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

//查找并返回所有未花费的transaction outputs
func (bc *BlockChain) FindUTXO(pubKeyHash []byte) []TXOutput {
	var UTXOs []TXOutput
	unspentTransactions := bc.FindUnspentTransactions(pubKeyHash)

	for _, tx := range unspentTransactions {
		for _, out := range tx.VOut {
			if out.IsLockedWithKey(pubKeyHash) {
				UTXOs = append(UTXOs, out)
			}
		}
	}
	return UTXOs
}

//找到有未花费的输出交易
func (bc *BlockChain) FindUnspentTransactions(pubKeyHash []byte) []Transaction {
	var unspentTXs []Transaction
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transaction {
			txID := hex.EncodeToString(tx.ID)

		Output:
			for outIdx, out := range tx.VOut {
				if spentTXOs[txID] != nil {
					for _, spentOutIdx := range spentTXOs[txID] {
						if spentOutIdx == outIdx {
							continue Output
						}
					}
				}

				if out.IsLockedWithKey(pubKeyHash) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}

			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin {
					if in.UsesKey(pubKeyHash) {
						inTxId := hex.EncodeToString(in.Txid)
						spentTXOs[inTxId] = append(spentTXOs[inTxId], in.Vout)
					}
				}
			}
		}
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return unspentTXs
}

//通过ID查找交易
func (bc *BlockChain)FindTransaction(ID []byte) (Transaction, error){
	bci := bc.Iterator()

	for{
		block := bci.Next()

		for _,tx := range block.Transaction{
			if bytes.Compare(tx.ID, ID) == 0{
				return *tx, nil
			}
		}

		if len(block.PrevBlockHash) == 0{
			break
		}
	}

	return Transaction{}, errors.New("Transaction is not found")
}

//查找并返回未花费的outputs
func (bc *BlockChain) FindSpendableOutPuts(pubKeyHash []byte, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTXs := bc.FindUnspentTransactions(pubKeyHash)
	accumulated := 0

Work:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.VOut {
			if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOutputs
}

