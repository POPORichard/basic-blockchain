package address

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
)

const walletFile = "wallet.dat"

type Wallets struct {
	Wallets map[string]*Wallet
}

//创建钱包并尝试从文件中读取
func NewWallets() (*Wallets, error){
	wallets := Wallets{}

	wallets.Wallets = make(map[string]*Wallet)

	err := wallets.LoadFromFile()

	return &wallets,err
}

//向wallets中添加一个wallet
func (ws *Wallets)CreateWallet() string{
	wallet := NewWallet()
	address := fmt.Sprintf("%s", wallet.GetAddress())

	ws.Wallets[address] = wallet

	return address
}

//返回wallet中存储的地址
func (ws *Wallets) GetAddresses() []string{
	var addresses []string
	for address := range ws.Wallets{
		addresses = append(addresses, address)
	}

	return addresses
}

//根据地址返回wallet
func (ws Wallets) GetWallet (address string) Wallet{
	return *ws.Wallets[address]
}

//从文件读取wallets
func (ws *Wallets)LoadFromFile() error{

	if _,err := os.Stat(walletFile); os.IsNotExist(err){
		return err
	}

	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil{
		panic(err)
	}

	var wallets Wallets

	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err !=nil{
		panic(err)
	}

	ws.Wallets = wallets.Wallets

	return nil
}

// 将wallets存储到文件
func (ws Wallets)SaveToFile(){
	var content bytes.Buffer

	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	if err != nil{
		panic(err)
	}

	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
}
