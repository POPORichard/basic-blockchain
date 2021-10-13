package main

import (
	"basic-blockchain/CLI"
	"basic-blockchain/database"
)

func main(){
	blockChain := database.NewBlockChain()
	defer blockChain.Db.Close()

	cli := CLI.CLI{Bc:blockChain}
	cli.Run()

}
