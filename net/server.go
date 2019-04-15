package net

import (
	"blockchain_go/blockchain"
	"blockchain_go/log"
	"blockchain_go/utils"
	"fmt"
	"io/ioutil"
	"net"
	"sync"
	"time"
)

const (
	maxConnectPeer           = 20 // 最大允许连接的节点数量
	rePeerDiscoveryThreshold = 1  // 对等节点数量小于rePeerDiscoveryThreshold时，再次运行节点发现
)
const protocol = "tcp"
const nodeVersion = 1
const commandLength = 12

var nodeAddress string
var miningAddress string
var CenterNode = "localhost:3000"

// TODO: 使用connectingPeersMutex互斥访问connectingPeers和activePeers
var mutex sync.Mutex

// 正在连接的节点 节点地址->连接时间戳
//var connectingPeers = map[string]int64{}
var connectingPeers = sync.Map{}

// 已连接的节点 节点地址->上次接收到节点消息时的时间戳
//var activePeers = map[string]int64{}
var activePeers = sync.Map{}

// initPeerDiscovery 初始化节点发现
func initPeerDiscovery(bc *blockchain.Blockchain) {
	nowTS := time.Now().Unix()
	peers := dnsSeedPeerDiscovery()
	log.Net.Printf("Find peers %#v\n", peers)
	for _, addr := range peers {
		sendVersion(addr, bc)
		//connectingPeers[addr] = nowTS
		connectingPeers.Store(addr, nowTS)
	}

}

// connnectionCheck 周期性检查对等节点连接状况
// 依赖 handleConnection 函数在每次收到对等节点消息时，更新activePeers[peerAddr]时间戳
func connnectionCheck(bc *blockchain.Blockchain) {
	const (
		checkIntervalSec  = 600  // 检查间隔 10min
		sendPingThreshold = 600  // 不活跃多长时间触发发送ping消息 10min
		inactiveThreshold = 1800 // 不活跃多长时间删除连接 30min
	)

	for {
		time.Sleep(checkIntervalSec * time.Second)
		// 检查已连接节点
		activePeers.Range(func(peerAddr, ts interface{}) bool {
			if time.Now().Unix()-ts.(int64) > inactiveThreshold {
				activePeers.Delete(peerAddr)
			} else if time.Now().Unix()-ts.(int64) > sendPingThreshold {
				sendPing(peerAddr.(string))
			}
			return true
		})

		if lenSycnMap(&activePeers) < rePeerDiscoveryThreshold {
			initPeerDiscovery(bc)
		}
	}
}

func handleConnection(conn net.Conn, bc *blockchain.Blockchain) {
	defer conn.Close()
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Net.Panic(err)
	}
	command := bytesToCommand(request[:commandLength])

	fromAddr := fmt.Sprintf("localhost:%s", request[len(request)-4:])
	//fromAddr := conn.RemoteAddr().String()
	request = request[:len(request)-4]

	// 远程节点不在当前建立连接的对等节点当中时，仅允许接收version、verack消息
	if !existInSyncMap(&activePeers, fromAddr) && command != "version" && command != "verack" {
		return
	}

	if existInSyncMap(&activePeers, fromAddr) {
		activePeers.Store(fromAddr, time.Now().Unix())
	}

	log.Net.Printf("Received %s command from %s\n", command, fromAddr)

	switch command {
	case "addr":
		var msg addr
		utils.GobDecode(request[commandLength:], &msg)
		msg.handleMsg(bc, fromAddr)
	case "block":
		var msg block
		utils.GobDecode(request[commandLength:], &msg)
		msg.handleMsg(bc, fromAddr)
	case "inv":
		var msg inv
		utils.GobDecode(request[commandLength:], &msg)
		msg.handleMsg(bc, fromAddr)
	case "getblocks":
		var msg getblocks
		utils.GobDecode(request[commandLength:], &msg)
		msg.handleMsg(bc, fromAddr)
	case "getdata":
		var msg getdata
		utils.GobDecode(request[commandLength:], &msg)
		msg.handleMsg(bc, fromAddr)
	case "tx":
		var msg tx
		utils.GobDecode(request[commandLength:], &msg)
		msg.handleMsg(bc, fromAddr)
	case "version":
		var msg version
		utils.GobDecode(request[commandLength:], &msg)
		msg.handleMsg(bc, fromAddr)
	case "verack":
		var msg verack
		utils.GobDecode(request[commandLength:], &msg)
		msg.handleMsg(bc, fromAddr)
	case "ping":
		var msg ping
		utils.GobDecode(request[commandLength:], &msg)
		msg.handleMsg(bc, fromAddr)
	default:
		log.Net.Println("Unknown command!")
	}
}

// StartServer starts a node
func StartServer(nodeID, minerAddress string) {
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	miningAddress = minerAddress

	ln, err := net.Listen(protocol, nodeAddress)
	if err != nil {
		log.Net.Panic(err)
	}
	defer ln.Close()

	bc := blockchain.GetBlockchain()

	go initPeerDiscovery(bc)

	go connnectionCheck(bc)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Net.Panic(err)
		}
		go handleConnection(conn, bc)
	}
}

func GetActivePeers() map[string]int64 {
	peers := map[string]int64{}
	activePeers.Range(func(key, value interface{}) bool {
		peers[key.(string)] = value.(int64)
		return true
	})
	return peers
}

func GetNodeAddr() string {
	return nodeAddress
}
