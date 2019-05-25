package cli

import (
	"blockchain_go/blockchain"
	"blockchain_go/miner"
	"blockchain_go/net"
	"blockchain_go/transaction"
	w "blockchain_go/wallet"
	"encoding/hex"
	"fmt"
	"log"
)

func (cli *CLI) send(from, to string, amount int, nodeID string, mineNow bool) {
	if !w.ValidateAddress(from) {
		log.Panic("ERROR: Sender address is not valid")
	}
	if !w.ValidateAddress(to) {
		log.Panic("ERROR: Recipient address is not valid")
	}

	bc := blockchain.GetBlockchain()
	UTXOSet := blockchain.UTXOSet{bc}
	defer bc.Close()

	wallets, err := w.NewWallets(nodeID)
	if err != nil {
		log.Panic(err)
	}
	wallet := wallets.GetWallet(from)

	tx := w.NewUTXOTransaction(&wallet, to, amount, &UTXOSet)

	if mineNow {
		cbTx := transaction.NewCoinbaseTX(from, "")
		txs := []*transaction.Transaction{cbTx, tx}

		newBlock := miner.MineBlock(bc, txs, map[string]transaction.Transaction{hex.EncodeToString(tx.ID): *tx})
		UTXOSet.Update(newBlock)
	} else {
		net.SendTx(net.CenterNode, tx)
	}

	fmt.Println("Success!")
}
