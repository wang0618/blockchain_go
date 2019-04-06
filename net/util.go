package net

import (
	"blockchain_go/log"
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

// dnsSeedPeerDiscovery 通过DNS Seed获取对等节点地址
// 目前直接返回 [3000, nodeID) 内的地址，所以启动节点时要按照NODE_ID升序启动
// TODO: 使用DNS Seed查询
func dnsSeedPeerDiscovery() []string {
	nodeID := os.Getenv("NODE_ID")
	currID, _ := strconv.Atoi(nodeID)

	var peers []string
	for i := 3000; i < currID; {
		peers = append(peers, fmt.Sprintf("localhost:%d", i))
	}
	return peers
}

func gobDecode(data []byte, e interface{}) {
	var buff bytes.Buffer
	buff.Write(data)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(e)
	if err != nil {
		log.Net.Panic(err)
	}
}

func gobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		log.Net.Panic(err)
	}

	return buff.Bytes()
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

func sendData(addr string, data []byte) {
	conn, err := net.Dial(protocol, addr)
	if err != nil {
		log.Net.Printf("%s is not available\n", addr)
		var updatedNodes []string

		for _, node := range knownNodes {
			if node != addr {
				updatedNodes = append(updatedNodes, node)
			}
		}

		knownNodes = updatedNodes

		return
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Net.Panic(err)
	}
}
