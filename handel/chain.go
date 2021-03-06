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
func (bc *BlockChain) MineBlock(transactions []*Transaction) *Block{
	var prevBlockHash []byte
	var lastHeight int

	for _,tx := range transactions{
		// TODO: 错误处理，不要panic
		if bc.VerifyTransaction(tx) != true{
			panic("ERROR: Invalid transaction")
		}
	}

	//从数据库中读取最后一个blockChain，并由此挖出下一个
	err := bc.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucketName))
		prevBlockHash = b.Get([]byte("l"))

		blockData := b.Get(prevBlockHash)
		block := DeserializeBlock(blockData)

		lastHeight = block.Height

		return nil
	})
	if err != nil {
		fmt.Println("Error when view prevBlockHash in db err :", err)
	}

	//开始计算下一个
	newBlock := NewBlock(transactions, prevBlockHash, lastHeight+1)

	//建立与数据库的读写链接(Update)
	err = bc.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucketName))
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

	return newBlock
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
	if tx.IsCoinbase(){
		return true
	}
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

//查找所有未花费的交易输出并返回删除了花费输出的交易
func (bc *BlockChain) FindUTXO() map[string]TXOutputs{
	UTXO := make(map[string]TXOutputs)
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
				outs := UTXO[txID]
				outs.Outputs = append(outs.Outputs, out)
				UTXO[txID] = outs
			}

			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin {
					inTxID := hex.EncodeToString(in.Txid)
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
				}
			}
		}
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return UTXO
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

func(bc *BlockChain)AddBlock (block *Block){
	err := bc.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucketName))
		blockInDb := b.Get(block.Hash)

		if blockInDb != nil{
			return nil
		}

		blockData := block.Serialize()
		err := b.Put(block.Hash, blockData)
		if err != nil{
			panic(err)
		}

		lastHash := b.Get([]byte("l"))
		lastBlockData := b.Get(lastHash)
		lastBlock := DeserializeBlock(lastBlockData)

		if block.Height > lastBlock.Height{
			err = b.Put([]byte("l"), block.Hash)
			if err != nil{
				panic(err)
			}
			bc.Tip = block.Hash
		}
		return nil
	})
	if err != nil{
		panic(err)
	}
}

func (bc *BlockChain) GetBestHeight() int{
	var lastBlock Block
	err := bc.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucketName))
		lastHash := b.Get([]byte("l"))
		blockData := b.Get(lastHash)
		lastBlock = *DeserializeBlock(blockData)

		return nil
	})
	if err != nil{
		panic(err)
	}

	return lastBlock.Height
}

func (bc *BlockChain) GetBlock(blockHash []byte) (Block, error){
	var block Block

	err := bc.Db.View(func(tx *bolt.Tx)error{
		b := tx.Bucket([]byte(BlocksBucketName))

		blockData := b.Get(blockHash)

		if blockData == nil{
			return  errors.New("block is not found")
		}

		block = *DeserializeBlock(blockData)

		return nil
	})

	return block, err
}

func (bc *BlockChain)GetBlockHashes() [][]byte{
	var blocks [][]byte
	bci := bc.Iterator()

	for{
		block := bci.Next()

		blocks = append(blocks, block.Hash)

		if len(block.PrevBlockHash) == 0{
			break
		}
	}

	return blocks
}