package handel

func (blockChain *BlockChain) AddBlock(data string) {
	prevBlock := blockChain.Blocks[len(blockChain.Blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	blockChain.Blocks = append(blockChain.Blocks, newBlock)
}

func NewBlockChain() *BlockChain {
	return &BlockChain{[]*Block{NewGenesisBlock()}}
}



