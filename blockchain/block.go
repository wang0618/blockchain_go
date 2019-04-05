package blockchain

import (
	"blockchain_go/transaction"
	"blockchain_go/transaction/merkletree"
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"log"
	"time"
)

// Block represents a block in the blockchain
type Block struct {
	Timestamp     int64
	Transactions  []*transaction.Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
	Height        int
	Difficulty    []byte
}

// NewBlock 创建一个不含工作量证明的区块
func NewBlock(transactions []*transaction.Transaction, prevBlockHash []byte, height int, difficult []byte) *Block {
	block := &Block{time.Now().Unix(), transactions, prevBlockHash, []byte{}, 0, height, difficult}
	return block
}

// NewGenesisBlock 返回硬编码的创世区块
func NewGenesisBlock() *Block {
	const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"
	const genesisCoinbaseAddress = "168C3RJbprmpnxNry49ftjWGfFGQTNeDsU"

	// 以下常量由 CLI.createGenesisBlock(delta_sec) 函数生成
	const nonce = 242227
	const hashStr = "000018d8e509cbba6cdc57c0574f016c51e76eb47c8c46421b57fb7c1b7061b9"
	const genesisTimestamp = int64(1554389243)
	const genesisDifficultyStr = "00002fef55cb2c5d9377bb5fdbaca0fc0861163f44015b872e0101a66d240e76"

	difficulty, _ := hex.DecodeString(genesisDifficultyStr)
	hash, _ := hex.DecodeString(hashStr)
	coinbase := transaction.NewCoinbaseTX(genesisCoinbaseAddress, genesisCoinbaseData)
	block := &Block{genesisTimestamp, []*transaction.Transaction{coinbase}, []byte{}, hash, nonce, 0, difficulty}

	return block
}

// HashTransactions returns a hash of the transactions in the block
func (b *Block) HashTransactions() []byte {
	var transactions [][]byte

	for _, tx := range b.Transactions {
		transactions = append(transactions, tx.Serialize())
	}
	mTree := merkletree.NewMerkleTree(transactions)

	return mTree.RootNode.Data
}

// Serialize serializes the block
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// DeserializeBlock deserializes a block
func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}
