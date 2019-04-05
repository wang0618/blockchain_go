package cli

import (
	"blockchain_go/blockchain"
	"blockchain_go/miner"
	"fmt"
	"strconv"
)

func (cli *CLI) printChain(nodeID string) {
	bc := blockchain.GetBlockchain()
	defer bc.Close()

	bci := bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("============ Block %x ============\n", block.Hash)
		fmt.Printf("Height: %d\n", block.Height)
		fmt.Printf("Prev. block: %x\n", block.PrevBlockHash)
		pow := miner.NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Printf("Difficulty: %x\n\n", block.Difficulty)
		for _, tx := range block.Transactions {
			fmt.Println(tx)
		}
		fmt.Printf("\n\n")

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
