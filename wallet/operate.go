package wallet

import (
	"blockchain_go/blockchain"
	"blockchain_go/utils"
	"encoding/hex"
	"fmt"
	"log"
	"strings"

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
	UTXOSet.Blockchain.SignTransaction(&tx, wallet_.PrivateKey)

	return &tx
}

// NewUTXOTransactionByWallets creates a new transaction from node wallets
func NewUTXOTransactionByWallets(wallets *Wallets, to string, amount int, UTXOSet *blockchain.UTXOSet) []*transaction.Transaction {
	transactionSet := make([]*transaction.Transaction, 0)
	// 查看wallets余额是否充足
	totalBalance := 0
	walletCount := 0
	balance := make([]string, 0)
	for address := range wallets.Wallets {
		balance = append(balance, address)
		walletBalance := UTXOSet.GetAddressBalance(address)
		totalBalance += walletBalance
		walletCount += 1
		if totalBalance >= amount {
			break
		}
	}

	if totalBalance < amount {
		log.Panic("ERROR: Not enough funds")
	}

	index := 0
	// 构造交易
	for address, wallet := range wallets.Wallets {
		amount -= UTXOSet.GetAddressBalance(address)
		// 去除余额为零的wallet
		if UTXOSet.GetAddressBalance(address) == 0 || strings.Compare(address, to) == 0 {
			index++
			continue
		}

		if amount <= 0 {
			transactionSet = append(transactionSet, NewUTXOTransaction(wallet, to, amount+UTXOSet.GetAddressBalance(address), UTXOSet))
			break
		} else {
			transactionSet = append(transactionSet, NewUTXOTransaction(wallet, to, UTXOSet.GetAddressBalance(address), UTXOSet))
		}
		index++
	}

	return transactionSet
}
