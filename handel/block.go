package handel

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)
// 计算块hash
func (block *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(block.Timestamp, 10))
	headers := bytes.Join([][]byte{block.PrevBlockHash, block.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)

	block.Hash = hash[:]
}

// 创建新块
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}}
	block.SetHash()
	return block
}

// 创建创世块
func NewGenesisBlock() *Block{
	return NewBlock("Genesis Block", []byte{})
}
