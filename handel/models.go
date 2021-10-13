package handel

import "github.com/boltdb/bolt"

type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

type BlockChain struct {
	Tip []byte
	Db *bolt.DB
}
