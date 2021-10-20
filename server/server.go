package server

import (
	"basic-blockchain/handel"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const protocol = "tcp"
const nodeVersion = 1
const commandLength = 12

var nodeAddress string
var miningAddress string
var KnownNodes = []string{"localhost:3000"}
var mempool = make(map[string]handel.Transaction)

type addr struct {
	AddrList []string
}

type verzion struct {
	Version    int
	BestHeight int
	AddrFrom   string
}

func extractCommand(request []byte) []byte {
	return request[:commandLength]
}

//处理收到的数据
func handleConnection(conn net.Conn, bc *handel.BlockChain) {
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}
	command := bytesToCommand(request[:commandLength])
	fmt.Printf("Received %s command\n", command)

	switch command {
	case "addr":
		handleAddr(request)
	case "block":
		handleBlock(request, bc)
	case "inv":
		handleInv(request, bc)
	case "getblocks":
		handleGetBlocks(request, bc)
	case "getdata":
		handleGetData(request, bc)
	case "tx":
		handleTx(request, bc)
	case "version":
		handleVersion(request, bc)
	case "getstart":
		handleGetStart(request)
	default:
		fmt.Println("Unknown command!")
	}

	conn.Close()
}

// StartServer starts a node
func StartServer(nodeID, minerAddress string) {
	printIP()
	fmt.Println("port is : ",nodeID)
	nodeAddress = fmt.Sprintf("172.16.10.49:%s", nodeID)
	miningAddress = minerAddress
	ln, err := net.Listen(protocol, nodeAddress)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	//启动节点后检查本地是否存在数据库
	//没有从其他节点获取得创世块
	dbFile := fmt.Sprintf(handel.DbFile, nodeID)
	if handel.DbExists(dbFile) == false{
		sendGetStart(KnownNodes[0])
		for{
			firstConn, err := ln.Accept()
			if err != nil{
				panic(err)
			}
			request, err := ioutil.ReadAll(firstConn)
			if err != nil {
				log.Panic(err)
			}
			command := bytesToCommand(request[:commandLength])
			fmt.Printf("Received First %s command\n", command)
			if command == "storeNode"{
				handleStoreFirstNode(request, nodeID)
				break
			}
			fmt.Println("waiting for first node")
		}
	}

	bc := handel.NewBlockchainLink(nodeID)
	//TODO:该defer无法执行，将导致数据库错误
	defer bc.Db.Close()

	if nodeAddress != KnownNodes[0] {
		sendVersion(KnownNodes[0], bc)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleConnection(conn, bc)
		//监听退出指令
		go func (){
			c := make(chan os.Signal)
			defer close(c)
			signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
			for s := range c{
				switch s {
				case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT: // ctrl + c
					bc.Db.Close()
					conn.Close()
					ln.Close()
					fmt.Println("退出", s)
				}
			}
		}()
	}
}

//判断节点是否已经存储
func nodeIsKnown(addr string) bool {
	for _, node := range KnownNodes {
		if node == addr {
			return true
		}
	}

	return false
}

func printIP(){
	addrs, err := net.InterfaceAddrs()
	if err != nil{
		panic(err)
	}
	for _,address := range addrs{
		if ipnet,ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback(){
			if ipnet.IP.To4() != nil{
				fmt.Println("IP: ",ipnet.IP.String())
			}
		}
	}

}