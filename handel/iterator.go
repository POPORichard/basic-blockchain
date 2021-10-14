package handel

import (
	"fmt"
	"github.com/boltdb/bolt"
)

type BlockChainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (bc *BlockChain) Iterator() *BlockChainIterator {
	bci := &BlockChainIterator{
		currentHash: bc.Tip,
		db:          bc.Db,
	}

	return bci
}

func (i *BlockChainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocksBucket"))
		encodedBlock := b.Get(i.currentHash)
		block = Deserialize(encodedBlock)

		return nil
	})

	if err != nil {
		fmt.Println("Iterator Error err:", err)
	}

	i.currentHash = block.PrevBlockHash

	return block
}