package cli

import (
	"blockchain_go/blockchain"
	"blockchain_go/wallet"
	"fmt"
	"log"
)

func (cli *CLI) getBalance(address, nodeID string) {
	if !wallet.ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}
	bc := blockchain.GetBlockchain()
	UTXOSet := blockchain.UTXOSet{bc}
	defer bc.Close()

	balance := UTXOSet.GetAddressBalance(address)

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}
