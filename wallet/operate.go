package wallet

import (
	"blockchain_go/blockchain"
	"blockchain_go/net"
	"blockchain_go/utils"
	"encoding/hex"
	"fmt"
	"log"

	"blockchain_go/transaction"
)

// NewUTXOTransaction creates a new transaction
func NewUTXOTransaction(wallet_ *Wallet, to string, amount int, UTXOSet *blockchain.UTXOSet) *transaction.Transaction {
	var inputs []transaction.TXInput
	var outputs []transaction.TXOutput

	pubKeyHash := utils.HashPubKey(wallet_.PublicKey)
	acc, validOutputs := UTXOSet.FindSpendableOutputs(pubKeyHash, amount)

	if acc < amount {
		log.Panic("ERROR: Not enough funds")
	}

	// Build a list of inputs
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			input := transaction.TXInput{txID, out, nil, wallet_.PublicKey}
			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs
	from := fmt.Sprintf("%s", wallet_.GetAddress())
	outputs = append(outputs, *transaction.NewTXOutput(amount, to))
	if acc > amount {
		outputs = append(outputs, *transaction.NewTXOutput(acc-amount, from)) // a change
	}

	tx := transaction.Transaction{nil, inputs, outputs}
	tx.ID = tx.Hash()
	UTXOSet.Blockchain.SignTransaction(&tx, wallet_.PrivateKey, net.MemPool)

	return &tx
}
