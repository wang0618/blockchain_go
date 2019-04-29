package cli

import (
	"blockchain_go/log"
	"blockchain_go/net"
	"fmt"
	"os"

	"blockchain_go/wallet"
	web "blockchain_go/webui"
)

func (cli *CLI) startNode(nodeID, minerAddress string, webui bool) {
	fmt.Printf("Starting node %s\n", nodeID)
	if len(minerAddress) > 0 {
		if wallet.ValidateAddress(minerAddress) {
			fmt.Println("Mining is on. Address to receive rewards: ", minerAddress)
		} else {
			log.CLI.Panic("Wrong miner address!")
		}
	}

	if webui {
		htmlPath := "./webui/html"
		if _, err := os.Stat(htmlPath); os.IsNotExist(err) {
			log.CLI.Panicf("%s is not exist, can't start web server!!", htmlPath)
		}
		log.LogToChan()
		go web.StartWebUIServer(htmlPath, nodeID)
	}
	net.StartServer(nodeID, minerAddress)
}
