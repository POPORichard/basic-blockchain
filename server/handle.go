package server

import (
	"basic-blockchain/handel"
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
)

var blocksInTransit = [][]byte{}	//跟踪下载到哪个块赖

func handleVersion(request []byte, bc *handel.BlockChain) {
	var buff bytes.Buffer
	var payload verzion

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	myBestHeight := bc.GetBestHeight()
	foreignerBestHeight := payload.BestHeight

	//TODO:验证合法性
	if myBestHeight < foreignerBestHeight {
		sendGetBlocks(payload.AddrFrom)
	} else if myBestHeight > foreignerBestHeight {
		sendVersion(payload.AddrFrom, bc)
	}

	// sendAddr(payload.AddrFrom)
	if !nodeIsKnown(payload.AddrFrom) {
		KnownNodes = append(KnownNodes, payload.AddrFrom)
	}
}

//获取所有blockChain
func handleGetBlocks(request []byte, bc *handel.BlockChain) {
	var buff bytes.Buffer
	var payload getblocks

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blocks := bc.GetBlockHashes()
	sendInv(payload.AddrFrom, "block", blocks)
}

//处理Inv
func handleInv(request []byte, bc *handel.BlockChain) {
	var buff bytes.Buffer
	var payload inv

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Recevied inventory with %d %s\n", len(payload.Items), payload.Type)

	if payload.Type == "block" {
		blocksInTransit = payload.Items

		blockHash := payload.Items[0]
		sendGetData(payload.AddrFrom, "block", blockHash)

		newInTransit := [][]byte{}
		for _, b := range blocksInTransit {
			if bytes.Compare(b, blockHash) != 0 {
				newInTransit = append(newInTransit, b)
			}
		}
		blocksInTransit = newInTransit
	}

	if payload.Type == "tx" {
		txID := payload.Items[0]

		if mempool[hex.EncodeToString(txID)].ID == nil {
			sendGetData(payload.AddrFrom, "tx", txID)
		}
	}
}

//处理getData的请求
//TODO：检查自己是否有这个这个block或tx
func handleGetData(request []byte, bc *handel.BlockChain) {
	var buff bytes.Buffer
	var payload getdata

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	if payload.Type == "block" {
		block, err := bc.GetBlock([]byte(payload.ID))
		if err != nil {
			return
		}

		sendBlock(payload.AddrFrom, &block)
	}

	if payload.Type == "tx" {
		txID := hex.EncodeToString(payload.ID)
		tx := mempool[txID]

		SendTx(payload.AddrFrom, &tx)
		// delete(mempool, txID)
	}
}

//处理接受到的block数据
func handleBlock(request []byte, bc *handel.BlockChain) {
	var buff bytes.Buffer
	var payload block

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blockData := payload.Block
	block := handel.DeserializeBlock(blockData)

	fmt.Println("Recevied a new block!")
	//TODO：验证块是否合法
	bc.AddBlock(block)

	fmt.Printf("Added block %x\n", block.Hash)

	//判断是否还有需要下载的块
	if len(blocksInTransit) > 0 {
		blockHash := blocksInTransit[0]
		sendGetData(payload.AddrFrom, "block", blockHash)

		blocksInTransit = blocksInTransit[1:]
	} else {
		//TODO:考虑在添加新块时就updateUTXOSet
		//从新建立UTXOSet
		UTXOSet := handel.UTXOSet{bc}
		UTXOSet.Reindex()
	}
}

//处理接受到的TX数据
func handleTx(request []byte, bc *handel.BlockChain) {
	var buff bytes.Buffer
	var payload tx

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	//将交易解码后放入内存池
	txData := payload.Transaction
	tx := handel.DeserializeTransaction(txData)
	mempool[hex.EncodeToString(tx.ID)] = tx


	//如果当前节点是中心节点，则会将新交易转发到其他节点
	if nodeAddress == KnownNodes[0] {
		for _, node := range KnownNodes {
			if node != nodeAddress && node != payload.AddFrom {
				sendInv(node, "tx", [][]byte{tx.ID})
			}
		}
	} else {
		//若当前节点不是中心节点，就进行挖矿
		//当当前节点为矿工节点（设置了miningaddress）
		//并且内存池中有两个或以上的交易时开始挖矿
		if len(mempool) >= 2 && len(miningAddress) > 0 {
		MineTransactions:
			var txs []*handel.Transaction

			//验证交易是否合法 将合法交易添加到txs中
			for id := range mempool {
				tx := mempool[id]
				if bc.VerifyTransaction(&tx) {
					txs = append(txs, &tx)
				}
			}

			//若无合法交易 返回
			if len(txs) == 0 {
				fmt.Println("All transactions are invalid! Waiting for new ones...")
				return
			}

			//将本块的奖励同样放入交易中
			cbTx := handel.NewCoinbaseTX(miningAddress, "")
			txs = append(txs, cbTx)

			//出块后对UTXOSet重新建立
			//TODO：UpdateUTXOS
			newBlock := bc.MineBlock(txs)
			UTXOSet := handel.UTXOSet{bc}
			UTXOSet.Reindex()

			fmt.Println("New block is mined!")

			//从交易池中删去已经生成的交易
			for _, tx := range txs {
				txID := hex.EncodeToString(tx.ID)
				delete(mempool, txID)
			}

			//对所有已知节点广播
			for _, node := range KnownNodes {
				if node != nodeAddress {
					sendInv(node, "block", [][]byte{newBlock.Hash})
				}
			}

			//若池中仍有交易继续挖掘
			if len(mempool) > 0 {
				goto MineTransactions
			}
		}
	}
}

//处理新增地址的请求
func handleAddr(request []byte) {
	var buff bytes.Buffer
	var payload addr

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	KnownNodes = append(KnownNodes, payload.AddrList...)
	fmt.Printf("There are %d known nodes now!\n", len(KnownNodes))
	requestBlocks()
}

func handleGetStart(request []byte){
	var buff bytes.Buffer
	var addrFrom string


	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&addrFrom)
	if err != nil {
		log.Panic(err)
	}

	content, err := ioutil.ReadFile("blockChain.db_genesis")
	if err != nil {
		log.Panic(err)
	}
	payload := gobEncode(storefile{FileData:content})
	re := append(commandToBytes("storeNode"), payload...)

	sendData(addrFrom, re)

}

func handleStoreFirstNode(request []byte){
	var buff bytes.Buffer
	var payload storefile


	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	//if err != nil {
	//	log.Panic(err)
	//}
	err = ioutil.WriteFile("blockChain.db_genesis", payload.FileData, 0644)
	if err != nil{
		panic(err)
	}


	err = ioutil.WriteFile(handel.DbFile, payload.FileData, 0644)
	if err != nil{
		panic(err)
	}

}
