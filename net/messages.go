package net

import (
	"blockchain_go/blockchain"
	"blockchain_go/transaction"
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
	Transaction []byte
}

func sendInv(address, kind string, items [][]byte) {
	inventory := inv{kind, items}
	payload := gobEncode(inventory)
	request := append(commandToBytes("inv"), payload...)

	sendData(address, request)
}

func sendGetBlocks(address string) {
	payload := gobEncode(getblocks{})
	request := append(commandToBytes("getblocks"), payload...)

	sendData(address, request)
}

func sendGetData(address, kind string, id []byte) {
	payload := gobEncode(getdata{kind, id})
	request := append(commandToBytes("getdata"), payload...)

	sendData(address, request)
}

func SendTx(addr string, tnx *transaction.Transaction) {
	data := tx{tnx.Serialize()}
	payload := gobEncode(data)
	request := append(commandToBytes("tx"), payload...)

	sendData(addr, request)
}

func sendVersion(addr string, bc *blockchain.Blockchain) {
	bestHeight := bc.GetBestHeight()
	payload := gobEncode(version{nodeVersion, bestHeight})

	request := append(commandToBytes("version"), payload...)

	sendData(addr, request)
}

func sendAddr(address string) {
	nodes := addr{knownNodes}
	nodes.AddrList = append(nodes.AddrList)
	payload := gobEncode(nodes)
	request := append(commandToBytes("addr"), payload...)

	sendData(address, request)
}

func sendBlock(addr string, b *blockchain.Block) {
	data := block{b.Serialize()}
	payload := gobEncode(data)
	request := append(commandToBytes("block"), payload...)

	sendData(addr, request)
}
