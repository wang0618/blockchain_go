package cli

import (
	"blockchain_go/wallet"
	"fmt"
)

func (cli *CLI) recoverWallet(code []string, nodeID string) {
	wallets, _ := wallet.NewWallets(nodeID)

	address := wallets.RecoverWallet(code)
	fmt.Printf("Your address: %s\n", address)
	wallets.SaveToFile(nodeID)
}
