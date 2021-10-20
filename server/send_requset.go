package server

func requestBlocks() {
	for _, node := range KnownNodes {
		sendGetBlocks(node)
	}
}

//发送获取blockChain的请求
func sendGetBlocks(address string) {
	payload := gobEncode(getblocks{nodeAddress})
	request := append(commandToBytes("getblocks"), payload...)

	sendData(address, request)
}

//发送getData请求
func sendGetData(address, kind string, id []byte) {
	payload := gobEncode(getdata{nodeAddress, kind, id})
	request := append(commandToBytes("getdata"), payload...)

	sendData(address, request)
}

func sendAddr(address string) {
	nodes := addr{KnownNodes}
	nodes.AddrList = append(nodes.AddrList, nodeAddress)
	payload := gobEncode(nodes)
	request := append(commandToBytes("addr"), payload...)

	sendData(address, request)
}

func sendGetStart(address string){
	payload := gobEncode(nodeAddress)
	request := append(commandToBytes("getstart"), payload...)

	sendData(address, request)
}