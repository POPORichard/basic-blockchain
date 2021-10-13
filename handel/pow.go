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
			pow.block.Data,
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)
	return data
}

func (pow *ProofOfWork) Run() (int, []byte){
	var hashInt big.Int
	var hash[32]byte
	nonce := 0

	for {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1{
			break
		}else{
			if nonce == math.MaxInt32{
				fmt.Println("can not get answer! data : ",data)
				panic("pow Error")
			}
			nonce++
		}
	}
	return nonce, hash[:]
}



func IntToHex(n int64) []byte {
	return []byte(strconv.FormatInt(n, 16))
}

func (pow *ProofOfWork) Validate() bool{
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(pow.target) == -1
}


