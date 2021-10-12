package handel

type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
}

type BlockChain struct {
	Blocks []*Block
}
