package blockchain

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"log"

	"blockchain_go/transaction"
	"github.com/boltdb/bolt"
)

const utxoBucket = "chainstate"

// UTXO represents an unspent output from a transaction
type UTXOItem struct {
	Idx int
	transaction.TXOutput
}

// SerializeUTXOItems serializes UTXOItem slice
func SerializeUTXOItems(UTXOs []UTXOItem) []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(UTXOs)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

// DeserializeUTXOItems deserializes UTXOItem slice
func DeserializeUTXOItems(data []byte) []UTXOItem {
	var outputs []UTXOItem

	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&outputs)
	if err != nil {
		log.Panic(err)
	}

	return outputs
}

// UTXOSet represents UTXO set
type UTXOSet struct {
	Blockchain *Blockchain
}

// FindSpendableOutputs finds and returns unspent outputs to reference in inputs
func (u UTXOSet) FindSpendableOutputs(pubkeyHash []byte, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	accumulated := 0
	db := u.Blockchain.db

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			txID := hex.EncodeToString(k)
			UTXOItems := DeserializeUTXOItems(v)

			for _, UTXO := range UTXOItems {
				if UTXO.IsLockedWithKey(pubkeyHash) {
					if accumulated >= amount {
						return nil
					}

					accumulated += UTXO.Value
					unspentOutputs[txID] = append(unspentOutputs[txID], UTXO.Idx)
				}
			}
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return accumulated, unspentOutputs
}

// FindUTXO finds UTXO for a public key hash
func (u UTXOSet) FindUTXO(pubKeyHash []byte) []transaction.TXOutput {
	var UTXOs []transaction.TXOutput
	db := u.Blockchain.db

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			UTXOItems := DeserializeUTXOItems(v)

			for _, UTXO := range UTXOItems {
				if UTXO.IsLockedWithKey(pubKeyHash) {
					UTXOs = append(UTXOs, UTXO.TXOutput)
				}
			}
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return UTXOs
}

// CountTransactions returns the number of transactions in the UTXO set
func (u UTXOSet) CountTransactions() int {
	db := u.Blockchain.db
	counter := 0

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			counter++
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return counter
}

// Reindex rebuilds the UTXO set
func (u UTXOSet) Reindex() {
	db := u.Blockchain.db
	bucketName := []byte(utxoBucket)

	err := db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(bucketName)
		if err != nil && err != bolt.ErrBucketNotFound {
			log.Panic(err)
		}

		_, err = tx.CreateBucket(bucketName)
		if err != nil {
			log.Panic(err)
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	UTXO := u.Blockchain.findUTXO()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)

		for txID, UTXOItems := range UTXO {
			key, err := hex.DecodeString(txID)
			if err != nil {
				log.Panic(err)
			}

			err = b.Put(key, SerializeUTXOItems(UTXOItems))
			if err != nil {
				log.Panic(err)
			}
		}

		return nil
	})
}

func updateTx(b *bolt.Bucket, tx *transaction.Transaction) {
	// 删除交易输入的UTXO
	if tx.IsCoinbase() == false {
		for _, vin := range tx.Vin {
			var updatedOuts []UTXOItem
			outsBytes := b.Get(vin.Txid)
			outs := DeserializeUTXOItems(outsBytes)

			for _, out := range outs {
				if out.Idx != vin.Vout {
					updatedOuts = append(updatedOuts, out)
				}
			}

			if len(updatedOuts) == 0 {
				err := b.Delete(vin.Txid)
				if err != nil {
					log.Panic(err)
				}
			} else {
				err := b.Put(vin.Txid, SerializeUTXOItems(updatedOuts))
				if err != nil {
					log.Panic(err)
				}
			}

		}
	}
	// 新增新的UTXO
	var newOutputs []UTXOItem
	for idx, out := range tx.Vout {
		newOutputs = append(newOutputs, UTXOItem{idx, out})
	}

	err := b.Put(tx.ID, SerializeUTXOItems(newOutputs))
	if err != nil {
		log.Panic(err)
	}
}

// Update updates the UTXO set with transactions from the Block
// The Block is considered to be the tip of a blockchain
func (u UTXOSet) Update(block *Block) {
	db := u.Blockchain.db

	err := db.Update(func(dbtx *bolt.Tx) error {
		b := dbtx.Bucket([]byte(utxoBucket))

		for _, tx := range block.Transactions {
			updateTx(b, tx)
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

// UpdateForTx 根据交易更新UTXO
func (u UTXOSet) UpdateForTx(tx *transaction.Transaction) {
	db := u.Blockchain.db

	err := db.Update(func(dbtx *bolt.Tx) error {
		b := dbtx.Bucket([]byte(utxoBucket))
		updateTx(b, tx)
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}
