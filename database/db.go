package database

import (
	"basic-blockchain/handel"
	"fmt"
	"github.com/boltdb/bolt"
)

const dbFile = "my.db"

func NewBlockChain() *handel.BlockChain {
	var tip []byte

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		fmt.Println("Error when open database err :", err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocksBucket"))
		//若db中不存在
		if b == nil {
			genesis := handel.NewGenesisBlock()

			b, err := tx.CreateBucket([]byte("blockBucket"))
			if err != nil {
				fmt.Println("Error when create blockBucket err :", err)
			}

			err = b.Put(genesis.Hash, genesis.Serialize())
			if err != nil {
				fmt.Println("Error when put genesis err :", err)
			}

			err = b.Put([]byte("l"), genesis.Hash)
			if err != nil {
				fmt.Println("", err)
			}

			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}

		return nil
	})

	return &handel.BlockChain{Tip: tip, Db: db}
}
