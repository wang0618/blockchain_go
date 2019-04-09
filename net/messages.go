package net

import (
	"blockchain_go/blockchain"
	"blockchain_go/transaction"
	"blockchain_go/utils"
)

type version struct {
	Version    int
	BestHeight int
}

type verack struct{}

type addr struct {
	AddrList []string
}

type block struct {
	Block []byte
}

type getblocks struct {
	StartHash []byte
}

type getdata struct {
	Type string
	ID   []byte
}

type inv struct {
	Type  string
	Items [][]byte
}

type tx struct {
	Transaction transaction.Transaction
}

type ping struct{}
type pong struct{}

func sendInv(address, kind string, items [][]byte) error {
	inventory := inv{kind, items}
	payload := utils.GobEncode(inventory)
	request := append(commandToBytes("inv"), payload...)

	return sendData(address, request)
}

func sendGetBlocks(address string, startHash []byte) error {
	payload := utils.GobEncode(getblocks{startHash})
	request := append(commandToBytes("getblocks"), payload...)

	return sendData(address, request)
}

func sendGetData(address, kind string, id []byte) error {
	payload := utils.GobEncode(getdata{kind, id})
	request := append(commandToBytes("getdata"), payload...)

	return sendData(address, request)
}

func SendTx(addr string, tnx *transaction.Transaction) error {
	data := tx{*tnx}
	payload := utils.GobEncode(data)
	request := append(commandToBytes("tx"), payload...)

	return sendData(addr, request)
}

func sendVersion(addr string, bc *blockchain.Blockchain) error {
	bestHeight := bc.GetBestHeight()
	payload := utils.GobEncode(version{nodeVersion, bestHeight})

	request := append(commandToBytes("version"), payload...)

	return sendData(addr, request)
}

func sendAddr(address string) error {
	addrs := make([]string, 0, lenSycnMap(&activePeers))
	activePeers.Range(func(addr, value interface{}) bool {
		addrs = append(addrs, addr.(string))
		return true
	})
	nodes := addr{addrs}
	payload := utils.GobEncode(nodes)
	request := append(commandToBytes("addr"), payload...)

	return sendData(address, request)
}

func sendBlock(addr string, b *blockchain.Block) error {
	data := block{b.Serialize()}
	payload := utils.GobEncode(data)
	request := append(commandToBytes("block"), payload...)

	return sendData(addr, request)
}

func sendVerack(addr string) error {
	payload := utils.GobEncode(verack{})
	request := append(commandToBytes("verack"), payload...)

	return sendData(addr, request)
}

func sendPing(addr string) error {
	payload := utils.GobEncode(ping{})
	request := append(commandToBytes("ping"), payload...)

	return sendData(addr, request)
}

func sendPong(addr string) error {
	payload := utils.GobEncode(pong{})
	request := append(commandToBytes("pong"), payload...)

	return sendData(addr, request)
}
