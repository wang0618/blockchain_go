package cli

import (
	"blockchain_go/wallet"
	"fmt"
	"log"
)

func (cli *CLI) recoverWallet(code []string, nodeID string) {
	wallets, err := wallet.NewWallets(nodeID)
	if err != nil {
		log.Panic(err)
	}
	address := wallets.RecoverWallet(code)
	fmt.Printf("Your address: %s\n", address)
	wallets.SaveToFile(nodeID)
}
