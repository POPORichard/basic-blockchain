package CLI

import (
	"basic-blockchain/address"
	"basic-blockchain/handel"
	"basic-blockchain/server"
	"fmt"
	"log"
)

//打钱
func (cli *CLI) send(from, to string, amount int, nodeID string, mineNow bool) {
	if !address.ValidateAddress(from) {
		log.Panic("ERROR: Sender address is not valid")
	}
	if !address.ValidateAddress(to) {
		log.Panic("ERROR: Recipient address is not valid")
	}

	bc := handel.NewBlockchainLink(nodeID)
	defer bc.Db.Close()
	UTXOSet := handel.UTXOSet{BlockChain: bc}

	wallets, err := address.NewWallets(nodeID)
	if err != nil{
		panic(err)
	}

	wallet := wallets.GetWallet(from)

	tx := handel.NewUTXOTransaction(&wallet, to, amount, &UTXOSet)

	if mineNow{
		cbTx := handel.NewCoinbaseTX(from, "")
		txs := []*handel.Transaction{cbTx, tx}

		newBlock := bc.MineBlock(txs)
		UTXOSet.Update(newBlock)
	}else {
		server.SendTx(server.KnownNodes[0], tx)
	}

	fmt.Println("Success!")
}