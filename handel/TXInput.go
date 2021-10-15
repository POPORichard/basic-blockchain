package handel

import (
	add "basic-blockchain/address"
	"bytes"
)

//交易输入
type TXInput struct {
	Txid      []byte
	Vout      int		//交易输出的索引
	//ScriptSig string	//地址
	Signature []byte
	PubKey    []byte
}

//用公钥检测是否是其发起了交易
func (in *TXInput)UsesKey(pubKeyHash []byte) bool{
	lockingHash := add.HashPubKey(in.PubKey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}