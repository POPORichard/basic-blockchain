package handel

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"
)

type Block struct {
	Timestamp     int64
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
	Transaction   []*Transaction
	Height 		  int
}

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
func NewBlock(transaction []*Transaction, prevBlockHash []byte, height int) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
		Nonce:         0,
		Transaction:   transaction,
		Height:        height,
	}
	pow := NewProofOfWork(block)

	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// 创建创世块
func NewGenesisBlock(coinBase *Transaction) *Block {
	return NewBlock([]*Transaction{coinBase}, []byte{}, 0)
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
func DeserializeBlock(b []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(b))
	err := decoder.Decode(&block)
	if err != nil {
		fmt.Println("Error in Deserialize err: ", err)
		return nil
	}
	return &block
}

//hashTransactions 用于计算交易的hash
/*
比特币使用了一种更复杂的技术：
它将包含在一个块中的所有交易表示为Merkle 树，
并在工作量证明系统中使用树的根哈希。
这种方法允许快速检查一个块是否包含某个交易，
只有根哈希，而无需下载所有交易。
*/
func (block *Block) hashTransactions() []byte {
	//var txHashes [][]byte
	//var txHash [32]byte
	var transactions [][]byte

	for _, tx := range block.Transaction {
		//txHashes = append(txHashes, tx.ID)
		transactions = append(transactions, tx.Serialize())
	}
	//txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	mTree := NewMerkleTree(transactions)

	//return txHash[:]
	return mTree.RootNode.Data
}
