package cli

import (
	"blockchain_go/blockchain"
	"blockchain_go/miner"
	"blockchain_go/net"
	"blockchain_go/transaction"
	ws "blockchain_go/wallet"
	"fmt"
	"log"
)

func (cli *CLI) send(from, to string, amount int, nodeID string, mineNow bool) {
	if (!ws.ValidateAddress(from)) && (from != "wallets") {
		log.Panic("ERROR: Sender address is not valid")
	}
	if !ws.ValidateAddress(to) {
		log.Panic("ERROR: Recipient address is not valid")
	}

	bc := blockchain.GetBlockchain()
	UTXOSet := blockchain.UTXOSet{bc}
	defer bc.Close()

	wallets, err := ws.NewWallets(nodeID)
	if err != nil {
		log.Panic(err)
	}

	txs := make([]*transaction.Transaction, 0)
	if from != "wallets" {
		wallet := wallets.GetWallet(from)
		txs = append(txs, ws.NewUTXOTransaction(&wallet, to, amount, &UTXOSet))
	} else {
		txs = append(txs, ws.NewUTXOTransactionByWallets(wallets, to, amount, &UTXOSet)...)
	}
	//tx := w.NewUTXOTransaction(&wallet, to, amount, &UTXOSet)

	if mineNow {
		cbTx := transaction.NewCoinbaseTX(from, "")
		//txs := []*transaction.Transaction{cbTx, txs[0]}
		txs = append([]*transaction.Transaction{cbTx}, txs...)

		newBlock := miner.MineBlock(bc, txs)
		UTXOSet.Update(newBlock)
	} else {
		for _, tx := range txs {
			//_ = net.SendTx(net.CenterNode, tx)
			net.BroadcastTx(tx)
		}
	}

	fmt.Println("Success!")
}
