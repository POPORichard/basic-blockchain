package handel

import "github.com/boltdb/bolt"

type Block struct {
	Timestamp     int64
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
	Transaction   []*Transaction
}

type BlockChain struct {
	Tip []byte
	Db  *bolt.DB
}
