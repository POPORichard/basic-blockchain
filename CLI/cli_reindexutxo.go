package CLI

import (
	"basic-blockchain/handel"
	"fmt"
)

func (cli *CLI) reindexUTXO(nodeID string) {
	bc := handel.NewBlockchainLink(nodeID)
	UTXOSet := handel.UTXOSet{BlockChain: bc}
	UTXOSet.Reindex()

	count := UTXOSet.CountTransactions()
	fmt.Printf("Done! There are %d transactions in the UTXO set.\n", count)
}