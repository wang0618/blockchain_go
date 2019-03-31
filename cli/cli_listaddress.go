package cli

import (
	"blockchain_go/wallet"
	"fmt"
	"log"
)

func (cli *CLI) listAddresses(nodeID string) {
	wallets, err := wallet.NewWallets(nodeID)
	if err != nil {
		log.Panic(err)
	}
	addresses := wallets.GetAddresses()

	for _, address := range addresses {
		fmt.Println(address)
	}
}
