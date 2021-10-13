package handel

// 向链上添加块
func (blockChain *BlockChain) AddBlock(data string) {
	prevBlock := blockChain.Blocks[len(blockChain.Blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	blockChain.Blocks = append(blockChain.Blocks, newBlock)
}

// 创建新创世链
func NewBlockChain() *BlockChain {
	return &BlockChain{[]*Block{NewGenesisBlock()}}
}



