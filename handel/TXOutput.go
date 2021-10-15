package handel

import (
	add "basic-blockchain/address"
	"bytes"
)

//交易输出
type TXOutput struct {
	Value        int		//数量
	//ScriptPubKey string		//地址
	PubKeyHash []byte
}

//创建新output
func NewTXOutput(value int, address string) *TXOutput{
	txo := &TXOutput{
		Value:      value,
		PubKeyHash: nil,
	}
	txo.Lock([]byte(address))

	return txo
}

//对output签名
func (out *TXOutput)Lock(address []byte){
	pubKeyHash := add.Base58Decode(address)
	pubKeyHash = pubKeyHash[1:len(pubKeyHash)-4]
	out.PubKeyHash = pubKeyHash
}

//检查output是否可以用公钥解锁
func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool{
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}
