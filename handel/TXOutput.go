package handel

import (
	add "basic-blockchain/address"
	"bytes"
	"encoding/gob"
	"log"
)

//交易输出
type TXOutput struct {
	Value        int		//数量
	//ScriptPubKey string		//地址
	PubKeyHash []byte
}

//TXoutpus的集合
type TXOutputs struct {
	Outputs []TXOutput
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

//序列化TXOutputs
func (outs TXOutputs)Serialize() []byte{
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(outs)
	if err != nil{
		log.Panic(err)
	}

	return buff.Bytes()
}


//反序列化TXOutputs
func DeserializeOutputs(data []byte)TXOutputs{
	var outputs TXOutputs

	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&outputs)
	if err != nil{
		log.Panic(err)
	}

	return outputs
}
