package cli

import (
	"blockchain_go/wallet"
	"fmt"
)

func (cli *CLI) createWallet(nodeID string) {
	wallets, _ := wallet.NewWallets(nodeID)
	address, mnemonicCode := wallets.CreateWallet()
	wallets.SaveToFile(nodeID)

	fmt.Println("Please store your mnemonicCode in safe place")
	fmt.Println(mnemonicCode)
	fmt.Printf("Your new address: %s\n", address)
}
