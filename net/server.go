package net

import (
	"blockchain_go/blockchain"
	"fmt"
	"io/ioutil"
	"log"
	"net"
)

const protocol = "tcp"
const nodeVersion = 1
const commandLength = 12

var nodeAddress string
var miningAddress string
var CenterNode string = "localhost:3000"
var knownNodes = []string{CenterNode}

func handleConnection(conn net.Conn, bc *blockchain.Blockchain) {
	defer conn.Close()
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}
	command := bytesToCommand(request[:commandLength])
	fmt.Printf("Received %s command\n", command)

	fromAddr := conn.RemoteAddr().String()
	switch command {
	case "addr":
		var msg addr
		gobDecode(request[commandLength:], &msg)
		msg.handleMsg(bc, fromAddr)
	case "block":
		var msg block
		gobDecode(request[commandLength:], &msg)
		msg.handleMsg(bc, fromAddr)
	case "inv":
		var msg inv
		gobDecode(request[commandLength:], &msg)
		msg.handleMsg(bc, fromAddr)
	case "getblocks":
		var msg getblocks
		gobDecode(request[commandLength:], &msg)
		msg.handleMsg(bc, fromAddr)
	case "getdata":
		var msg getdata
		gobDecode(request[commandLength:], &msg)
		msg.handleMsg(bc, fromAddr)
	case "tx":
		var msg tx
		gobDecode(request[commandLength:], &msg)
		msg.handleMsg(bc, fromAddr)
	case "version":
		var msg version
		gobDecode(request[commandLength:], &msg)
		msg.handleMsg(bc, fromAddr)
	default:
		fmt.Println("Unknown command!")
	}
}

// StartServer starts a node
func StartServer(nodeID, minerAddress string) {
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	miningAddress = minerAddress

	ln, err := net.Listen(protocol, nodeAddress)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	bc := blockchain.GetBlockchain()

	if nodeAddress != knownNodes[0] {
		sendVersion(knownNodes[0], bc)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleConnection(conn, bc)
	}
}

func nodeIsKnown(addr string) bool {
	for _, node := range knownNodes {
		if node == addr {
			return true
		}
	}

	return false
}
