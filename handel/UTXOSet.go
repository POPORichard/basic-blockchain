package database

import (
	"basic-blockchain/handel"
	"encoding/hex"
	"github.com/boltdb/bolt"
	"log"
)

//unspent transaction out
type UTXOSet struct {
	BlockChain *handel.BlockChain
}

const UTXOBucket = "UTXOBucket"

//重构UTXOSet
func (u UTXOSet)Reindex(){
	db := u.BlockChain.Db
	bucketName := []byte(UTXOBucket)

	err := db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(bucketName)
		if err != nil && err != bolt.ErrBucketNotFound{
			panic(err)
		}

		_,err = tx.CreateBucket(bucketName)
		if err != nil{
			panic(err)
		}

		return nil
	})
	if err != nil {
		log.Panic("ERROR : can not create UTXOBucket")
	}

	UTXO := u.BlockChain.FindUTXO()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)

		for txID, outs := range UTXO{
			key, err := hex.DecodeString(txID)
			if err != nil {
				panic(err)
			}

			err = b.Put(key, outs.Serialize())
			if err != nil{
				panic(err)
			}
		}
		return nil
	})
}

//查找并返回unspent outputs to reference in inputs
func (u UTXOSet)FindSpendableOutputs(pubKeyHash []byte, amount int)(int, map[string][]int){
	unspentOUtputs := make(map[string][]int)
	accumulated := 0
	db := u.BlockChain.Db

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(UTXOBucket))
		c := b.Cursor()

		for k,v := c.First(); k!=nil;k,v = c.Next(){
			txID := hex.EncodeToString(k)
			outs := handel.DeserializeOutputs(v)

			for outIdx ,out := range outs.Outputs {
				if out.IsLockedWithKey(pubKeyHash) && accumulated < amount{
					accumulated += out.Value
					unspentOUtputs[txID] = append(unspentOUtputs[txID], outIdx)
				}

			}
		}
		return nil
	})

	if err != nil{
		panic(err)
	}

	return accumulated, unspentOUtputs
}

//查找公钥对应的UTXOs
func (u UTXOSet)FindUTXO(pubKeyHash []byte) []handel.TXOutput {
	var UTXOs []handel.TXOutput
	db := u.BlockChain.Db

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(UTXOBucket))
		c := b.Cursor()

		for k,v := c.First(); k!= nil; k,v = c.Next(){
			outs := handel.DeserializeOutputs(v)

			for _,out :=range outs.Outputs{
				if out.IsLockedWithKey(pubKeyHash){
					UTXOs = append(UTXOs, out)
				}
			}
		}
		return nil
	})
	if err != nil{
		panic(err)
	}
	return UTXOs
}

// Update 使用来自区块的交易更新 UTXO 集
// 该块被认为是区块链的尖端
func (u UTXOSet) Update(block *handel.Block){
	db := u.BlockChain.Db

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(UTXOBucket))

		for _,tx := range block.Transaction{
			if tx.IsCoinbase() == false{
				for _, vin := range tx.Vin{
					updateOuts := handel.TXOutputs{}
					outsBytes := b.Get(vin.Txid)
					outs := handel.DeserializeOutputs(outsBytes)

					for outIdx, out := range outs.Outputs{
						if outIdx != vin.Vout{
							updateOuts.Outputs = append(updateOuts.Outputs, out)
						}
					}
					if len(updateOuts.Outputs) == 0{
						err := b.Delete(vin.Txid)
						if err != nil{
							panic(err)
						}
					}else {
						err := b.Put(vin.Txid, updateOuts.Serialize())
						if err != nil{
							panic(err)
						}
					}
				}
			}
			newOutpust := handel.TXOutputs{}
			for _,out := range tx.VOut{
				newOutpust.Outputs = append(newOutpust.Outputs, out)
			}

			err := b.Put(tx.ID, newOutpust.Serialize())
			if err != nil{
				panic(err)
			}
		}
		return nil
	})
	if err != nil{
		panic(err)
	}
}

//返回UTXOSet中交易的数量
func (u UTXOSet) CountTransactions()int{
	db := u.BlockChain.Db
	counter := 0

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(UTXOBucket))
		c := b.Cursor()

		for k,_ := c.First();k != nil; k,_ =c.Next(){
			counter++
		}

		return nil
	})
	if err != nil{
		panic(err)
	}
	return counter
}