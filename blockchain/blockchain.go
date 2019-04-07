package blockchain

import (
	"blockchain_go/utils"
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"

	"blockchain_go/transaction"
	"github.com/boltdb/bolt"
)

const dbFile = "blockchain_%s.db"
const blocksBucket = "blocks"

var bc Blockchain

// Blockchain implements interactions with a DB
type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

// 首次运行时创建本地区块数据库
// 可通过NODE_ID环境变量来制定不同的数据库文件
func init() {
	nodeID := os.Getenv("NODE_ID")
	dbFile := fmt.Sprintf(dbFile, nodeID)

	dbExist := dbExists(dbFile)

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	if dbExist {
		// 数据库已存在，读取最顶端区块哈希

		var tip []byte

		err = db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(blocksBucket))
			tip = b.Get([]byte("l"))

			return nil
		})
		if err != nil {
			log.Panic(err)
		}

		bc = Blockchain{tip, db}

	} else {
		// 数据库不存在，将创世区块存入数据库

		genesis := newGenesisBlock()

		bc = Blockchain{genesis.Hash, db}

		err = db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic(err)
			}

			err = b.Put(genesis.Hash, genesis.Serialize())
			if err != nil {
				log.Panic(err)
			}

			err = b.Put([]byte("l"), genesis.Hash)
			if err != nil {
				log.Panic(err)
			}

			return nil
		})
		if err != nil {
			log.Panic(err)
		}

		UTXOSet := UTXOSet{&bc}
		UTXOSet.Reindex()
	}
}

func GetBlockchain() *Blockchain {
	return &bc
}

// LastBlockInfo 返回本地区块链中最新的区块
func (bc *Blockchain) LastBlockInfo() *Block {
	var block *Block
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash := b.Get([]byte("l"))

		blockData := b.Get(lastHash)
		block = DeserializeBlock(blockData)

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return block
}

var orphanBlocks = map[string]*Block{} // 游离区块， string(区块前驱哈希)->区块

// AddBlock saves the block into the blockchain
// 支持游离区块管理
func (bc *Blockchain) AddBlock(block *Block) {
	err := bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		blockInDb := b.Get(block.Hash)
		if blockInDb != nil {
			return nil
		}

		prevInDb := b.Get(block.PrevBlockHash)
		if prevInDb == nil {
			// block为游离区块
			orphanBlocks[string(block.PrevBlockHash)] = block
			return nil
		}
		err := b.Put(block.Hash, block.Serialize())
		utils.PanicIfError(err)

		currBlock := block
		for h := block.Hash; orphanBlocks[string(h)] != nil; h = orphanBlocks[string(h)].Hash {
			currBlock = orphanBlocks[string(h)]

			err := b.Put(currBlock.Hash, currBlock.Serialize())
			utils.PanicIfError(err)
		}

		lastHash := b.Get([]byte("l"))
		lastBlockData := b.Get(lastHash)
		lastBlock := DeserializeBlock(lastBlockData)

		if currBlock.Height > lastBlock.Height {
			err = b.Put([]byte("l"), currBlock.Hash)
			utils.PanicIfError(err)
			bc.tip = currBlock.Hash
		}

		return nil
	})
	utils.PanicIfError(err)
}

// FindTransaction finds a transaction by its ID
func (bc *Blockchain) FindTransaction(ID []byte) (transaction.Transaction, error) {
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return *tx, nil
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return transaction.Transaction{}, errors.New("Transaction is not found")
}

// FindUTXO 返回一个 txID -> []UTXO 的map
func (bc *Blockchain) findUTXO() map[string][]UTXOItem {
	UTXOs := make(map[string][]UTXOItem)       // txID ->  UTXO slice
	spentTXOs := make(map[string]map[int]bool) // txID ->  map(output_idx -> is_spent)
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

			for outIdx, out := range tx.Vout {
				// Was the output spent?
				if spentTXOs[txID] != nil && spentTXOs[txID][outIdx] {
					continue
				}

				UTXOs[txID] = append(UTXOs[txID], UTXOItem{outIdx, out})
			}

			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin {
					inTxID := hex.EncodeToString(in.Txid)
					if spentTXOs[inTxID] == nil {
						spentTXOs[inTxID] = map[int]bool{}
					}
					spentTXOs[inTxID][in.Vout] = true
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return UTXOs
}

// Iterator returns a BlockchainIterat
func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}

	return bci
}

// GetBestHeight returns the height of the latest block
func (bc *Blockchain) GetBestHeight() int {
	var lastBlock Block

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash := b.Get([]byte("l"))
		blockData := b.Get(lastHash)
		lastBlock = *DeserializeBlock(blockData)

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return lastBlock.Height
}

// GetBlock finds a block by its hash and returns it
func (bc *Blockchain) GetBlock(blockHash []byte) (Block, error) {
	var block Block

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		blockData := b.Get(blockHash)

		if blockData == nil {
			return errors.New("Block is not found.")
		}

		block = *DeserializeBlock(blockData)

		return nil
	})
	if err != nil {
		return block, err
	}

	return block, nil
}

// GetBlockHashes 获取本地区块链哈希列表
func (bc *Blockchain) GetBlockHashes() [][]byte {
	var blocks [][]byte
	bci := bc.Iterator()

	for {
		block := bci.Next()

		blocks = append(blocks, block.Hash)

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return blocks
}

// GetCurrentDifficult 返回当前区块链的挖矿难度值
// 当前，挖矿难度值固定，和创世区块难度值一致
// todo 参照比特币动态调整难度值
func (bc *Blockchain) GetCurrentDifficult() []byte {
	return newGenesisBlock().Difficulty
}

// SignTransaction signs inputs of a Transaction
func (bc *Blockchain) SignTransaction(tx *transaction.Transaction, privKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]transaction.Transaction)

	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	tx.Sign(privKey, prevTXs)
}

// VerifyTransaction verifies transaction input signatures
// 仅校验了签名，没有检测是否双花
// TODO 失败时返回具体错误原因
func (bc *Blockchain) VerifyTransactionSig(tx *transaction.Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	prevTXs := make(map[string]transaction.Transaction)

	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	return tx.VerifySig(prevTXs)
}

// VerifyTransaction verifies transaction input signatures
func (bc *Blockchain) Close() {
	bc.db.Close()
}

func dbExists(dbFile string) bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}
