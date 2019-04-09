package net

import (
	"blockchain_go/blockchain"
	"blockchain_go/log"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"sync"
)

// dnsSeedPeerDiscovery 通过DNS Seed获取对等节点地址
// 目前直接返回 [3000, nodeID) 内的地址，所以启动节点时要按照NODE_ID升序启动
// TODO: 使用DNS Seed查询
func dnsSeedPeerDiscovery() []string {
	nodeID := os.Getenv("NODE_ID")
	currID, _ := strconv.Atoi(nodeID)

	var peers []string
	for i := 3000; i < currID; i++ {
		peers = append(peers, fmt.Sprintf("localhost:%d", i))
	}
	return peers
}

func commandToBytes(command string) []byte {
	var bytes [commandLength]byte

	for i, c := range command {
		bytes[i] = byte(c)
	}

	return bytes[:]
}

func bytesToCommand(bytes []byte) string {
	var command []byte

	for _, b := range bytes {
		if b != 0x0 {
			command = append(command, b)
		}
	}

	return fmt.Sprintf("%s", command)
}

func extractCommand(request []byte) []byte {
	return request[:commandLength]
}

func sendData(addr string, data []byte) error {
	command := bytesToCommand(data[:commandLength])
	log.Net.Printf("Send %s msg to %s", command, addr)
	nodeID := os.Getenv("NODE_ID")
	data = append(data, []byte(nodeID)...)

	conn, err := net.Dial(protocol, addr)
	if err != nil {
		log.Net.Printf("%s is not available\n", addr)
		activePeers.Delete(addr)
		if lenSycnMap(&activePeers) < rePeerDiscoveryThreshold {
			go initPeerDiscovery(blockchain.GetBlockchain())
		}
		return err
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(data))
	return err
}

func lenSycnMap(p *sync.Map) (len int) {
	len = 0
	p.Range(func(key, value interface{}) bool {
		len++
		return true
	})
	return len
}

func existInSyncMap(p *sync.Map, key interface{}) bool {
	_, ok := p.Load(key)
	return ok
}
