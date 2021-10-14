package handel

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"
)

// 计算块hash
//func (block *Block) SetHash() {
//	timestamp := []byte(strconv.FormatInt(block.Timestamp, 10))
//	headers := bytes.Join([][]byte{block.PrevBlockHash, block.Data, timestamp}, []byte{})
//	hash := sha256.Sum256(headers)
//
//	block.Hash = hash[:]
//}

// 创建新块
//func NewBlock(data string, prevBlockHash []byte) *Block {
//	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}}
//	block.SetHash()
//	return block
//}

//带pow创建新block
func NewBlock(transaction []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
		Nonce:         0,
		Transaction:   transaction,
	}
	pow := NewProofOfWork(block)

	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// 创建创世块
func NewGenesisBlock(coinBase *Transaction) *Block {
	return NewBlock([]*Transaction{coinBase}, []byte{})
}

//序列化
func (block *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(block)
	if err != nil {
		fmt.Println("Error in Serialize err: ", err)
		return []byte{}
	}

	return result.Bytes()
}

// 反序列化
func Deserialize(b []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(b))
	err := decoder.Decode(&block)
	if err != nil {
		fmt.Println("Error in Deserialize err: ", err)
		return nil
	}
	return &block
}
