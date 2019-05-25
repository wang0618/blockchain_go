package miner

import (
	"blockchain_go/blockchain"
	"blockchain_go/transaction"
	"log"
)

// MineBlock mines a new block with the provided transactions
func MineBlock(bc *blockchain.Blockchain, transactions []*transaction.Transaction, mempool map[string]transaction.Transaction) *blockchain.Block {
	var lastHash []byte
	var lastHeight int

	for _, tx := range transactions {
		// TODO: ignore transaction if it's not valid
		if bc.VerifyTransactionSig(tx, mempool) != true {
			log.Panic("ERROR: Invalid transaction")
		}
	}

	lastBlock := bc.LastBlockInfo()
	lastHash = lastBlock.Hash
	lastHeight = lastBlock.Height

	newBlock := blockchain.NewBlock(transactions, lastHash, lastHeight+1, blockchain.GetBlockchain().GetCurrentDifficult())

	pow := NewProofOfWork(newBlock)
	nonce, hash := pow.Run()

	newBlock.Hash = hash[:]
	newBlock.Nonce = nonce

	bc.AddBlock(newBlock)

	return newBlock
}
