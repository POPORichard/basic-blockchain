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
	"errors"
	"fmt"
	"log"
	"math/big"
	"strings"
)
//奖励
const subsidy = 50

//交易输出
type TXOutput struct {
	Value        int		//数量
	//ScriptPubKey string		//地址
	PubKeyHash []byte
}

//交易输入
type TXInput struct {
	Txid      []byte
	Vout      int		//交易输出的索引
	//ScriptSig string	//地址
	Signature []byte
	PubKey    []byte
}

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
	bc.SignTransaction(&tx, wallet.PrivateKey)

	return &tx
}

func (bc *BlockChain) FindSpendableOutPuts(pubKeyHash []byte, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTXs := bc.FindUnspentTransactions(pubKeyHash)
	accumulated := 0

Work:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.VOut {
			if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOutputs
}

func (bc *BlockChain)FindTransaction(ID []byte) (Transaction, error){
	bci := bc.Iterator()

	for{
		block := bci.Next()

		for _,tx := range block.Transaction{
			if bytes.Compare(tx.ID, ID) == 0{
				return *tx, nil
			}
		}

		if len(block.PrevBlockHash) == 0{
			break
		}
	}

	return Transaction{}, errors.New("Transaction is not found")
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

//找到有未花费的输出交易
func (bc *BlockChain) FindUnspentTransactions(pubKeyHash []byte) []Transaction {
	var unspentTXs []Transaction
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transaction {
			txID := hex.EncodeToString(tx.ID)

		Output:
			for outIdx, out := range tx.VOut {
				if spentTXOs[txID] != nil {
					for _, spentOutIdx := range spentTXOs[txID] {
						if spentOutIdx == outIdx {
							continue Output
						}
					}
				}

				if out.IsLockedWithKey(pubKeyHash) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}

			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin {
					if in.UsesKey(pubKeyHash) {
						inTxId := hex.EncodeToString(in.Txid)
						spentTXOs[inTxId] = append(spentTXOs[inTxId], in.Vout)
					}
				}
			}
		}
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return unspentTXs
}

func (bc *BlockChain) FindUTXO(pubKeyHash []byte) []TXOutput {
	var UTXOs []TXOutput
	unspentTransactions := bc.FindUnspentTransactions(pubKeyHash)

	for _, tx := range unspentTransactions {
		for _, out := range tx.VOut {
			if out.IsLockedWithKey(pubKeyHash) {
				UTXOs = append(UTXOs, out)
			}
		}
	}
	return UTXOs
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

//用公钥检测是否是其发起了交易
func (in *TXInput)UsesKey(pubKeyHash []byte) bool{
	lockingHash := add.HashPubKey(in.PubKey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
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

//创建新output
func NewTXOutput(value int, address string) *TXOutput{
	txo := &TXOutput{
		Value:      value,
		PubKeyHash: nil,
	}
	txo.Lock([]byte(address))

	return txo
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
