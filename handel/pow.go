package handel

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"strconv"
)

const targetBits = 12

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	//fmt.Printf("\n%x\n", target)

	pow := &ProofOfWork{
		block:  block,
		target: target,
	}

	return pow
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.hashTransactions(),
			intToHex(pow.block.Timestamp),
			intToHex(int64(targetBits)),
			intToHex(int64(nonce)),
		},
		[]byte{},
	)
	return data
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	for {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			if nonce == math.MaxInt32 {
				fmt.Println("can not get answer! data : ", data)
				panic("pow Error")
			}
			nonce++
		}
	}
	return nonce, hash[:]
}

func intToHex(n int64) []byte {
	return []byte(strconv.FormatInt(n, 16))
}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(pow.target) == -1
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
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range block.Transaction {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}
