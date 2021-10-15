package handel

import (
	add "basic-blockchain/address"
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"
)
//奖励
const subsidy = 50

type Transaction struct {
	ID   []byte
	Vin  []TXInput
	VOut []TXOutput
}

// 创建一个新的coinbase交易
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

 	txin := TXInput{
		Txid:      []byte{},
		Vout:      -1,
		Signature: nil,
		PubKey:    []byte(data),
	}
	txout := NewTXOutput(subsidy, to)
	tx := Transaction{
		ID:   nil,
		Vin:  []TXInput{txin},
		VOut: []TXOutput{*txout},
	}
	tx.ID = tx.Hash()

	return &tx
}

//创建一个新的交易
func NewUTXOTransaction(from, to string, amount int, bc *BlockChain) *Transaction {
	var inputs []TXInput
	var outPuts []TXOutput

	wallets, err := add.NewWallets()
	if err != nil{
		panic(err)
	}

	wallet := wallets.GetWallet(from)
	pubKeyHash := add.HashPubKey(wallet.PublicKey)
	acc, validOutputs := bc.FindSpendableOutPuts(pubKeyHash, amount)

	if acc < amount {
		log.Panic("Error: Not enough funds")
	}

	//Build a list of inputs
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			panic("Error : Build a list of inputs")
		}

		for _, out := range outs {
			input := TXInput{txID, out, nil, wallet.PublicKey}
			inputs = append(inputs, input)
		}
	}

	//Build a list of outPuts
	outPuts = append(outPuts, *NewTXOutput(amount, to))
	if acc > amount {
		outPuts = append(outPuts, *NewTXOutput(acc - amount, from)) //as change
	}

	tx := Transaction{
		ID:   nil,
		Vin:  inputs,
		VOut: outPuts,
	}
	tx.ID = tx.Hash()
	//用私钥对交易签名
	bc.SignTransaction(&tx, wallet.PrivateKey)

	return &tx
}

//返回用户可读的交易信息
func (tx Transaction)string() string{
	var lines[]string

	lines = append(lines, fmt.Sprintf("---Transaction %x:", tx.ID))

	for i, input := range tx.Vin {

		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:      %x", input.Txid))
		lines = append(lines, fmt.Sprintf("       Out:       %d", input.Vout))
		lines = append(lines, fmt.Sprintf("       Signature: %x", input.Signature))
		lines = append(lines, fmt.Sprintf("       PubKey:    %x", input.PubKey))
	}

	for i, output := range tx.VOut {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %d", output.Value))
		lines = append(lines, fmt.Sprintf("       Script: %x", output.PubKeyHash))
	}

	return strings.Join(lines, "\n")
}

//创建交易副本用于签名
func(tx *Transaction) TrimmedCopy() Transaction{
	var inputs []TXInput
	var outputs []TXOutput

	for _,vin := range tx.Vin{
		inputs = append(inputs, TXInput{vin.Txid, vin.Vout,nil,nil})
	}

	for _,vout := range tx.VOut{
		outputs = append(outputs, TXOutput{vout.Value, vout.PubKeyHash})
	}

	txCopy := Transaction{tx.ID, inputs, outputs}

	return txCopy
}

//序列化transaction
func (tx Transaction)Serialize() []byte{
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

//返回transaction的hash
func (tx *Transaction)Hash() []byte{
	var hash[32]byte

	txCopy := *tx
	txCopy.ID = []byte{}

	hash = sha256.Sum256(txCopy.Serialize())

	return hash[:]
}

//对transaction中的每个input进行签名
func (tx *Transaction)Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction){
	if tx.IsCoinbase(){
		return
	}
	for _,vin := range tx.Vin{
		if prevTXs[hex.EncodeToString(vin.Txid)].ID == nil{
			panic("Error :previous transaction is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()

	for inID, vin := range txCopy.Vin{
		prevTX := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTX.VOut[vin.Vout].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Vin[inID].PubKey = nil

		r,s,err := ecdsa.Sign(rand.Reader,&privKey,txCopy.ID)
		if err != nil{
			panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)

		tx.Vin[inID].Signature = signature
	}
}


func (tx Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

//验证交易vin的签名
func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool{
	if tx.IsCoinbase(){
		return true
	}

	for _,vin := range tx.Vin{
		if prevTXs[hex.EncodeToString(vin.Txid)].ID == nil{
			panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for inID, vin := range tx.Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTx.VOut[vin.Vout].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Vin[inID].PubKey = nil

		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PubKey)
		x.SetBytes(vin.PubKey[:(keyLen / 2)])
		y.SetBytes(vin.PubKey[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubKey, txCopy.ID, &r, &s) == false {
			return false
		}
	}
	return true
}
