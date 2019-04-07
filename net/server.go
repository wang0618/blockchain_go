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
	rePeerDiscoveryThreshold = 3  // 对等节点数量小于rePeerDiscoveryThreshold时，再次运行节点发现
)
const protocol = "tcp"
const nodeVersion = 1
const commandLength = 12

var nodeAddress string
var miningAddress string
var CenterNode string = "localhost:3000"
var knownNodes = []string{CenterNode}

type connectStatus int

type connectingPeerStatus struct {
	status     connectStatus
	timestamp  int64
	versionMsg *version
}

const (
	needSendVer connectStatus = iota //  未发送version
	waitVer                          //  发送version等待回应
	waitVerAck                       // 等待VerAck
)

// TODO: 使用connectingPeersMutex互斥访问connectingPeers和activePeers
var mutex sync.Mutex
// 正在连接的节点 节点地址->连接状态
var connectingPeers = map[string]*connectingPeerStatus{}
// 已连接的节点 节点地址->上次接收到节点消息时的时间戳
var activePeers = map[string]int64{}

// initPeerDiscovery 初始化节点发现
func initPeerDiscovery(bc *blockchain.Blockchain) {
	nowTS := time.Now().Unix()
	peers := dnsSeedPeerDiscovery()
	log.Net.Printf("Find peers %#v\n", peers)
	for _, addr := range peers {
		connectingPeers[addr] = &connectingPeerStatus{status: needSendVer, timestamp: nowTS}
		go func() {
			sendVersion(addr, bc)
			connectingPeers[addr] = &connectingPeerStatus{status: waitVer, timestamp: nowTS}
		}()
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
		for peerAddr, ts := range activePeers {
			if time.Now().Unix()-ts > inactiveThreshold {
				delete(activePeers, peerAddr)
			} else if time.Now().Unix()-ts > sendPingThreshold {
				sendPing(peerAddr)
			}
		}

		if len(activePeers) < rePeerDiscoveryThreshold {
			initPeerDiscovery(bc)
		}

		// 检查正在连接的节点
		for peerAddr, s := range connectingPeers {
			if time.Now().Unix()-s.timestamp > inactiveThreshold {
				delete(connectingPeers, peerAddr)
			}
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

	log.Net.Printf("Received %s command from %s\n", command, fromAddr)

	// 远程节点不在当前建立连接的对等节点当中时，仅允许接收version、verack消息
	if activePeers[fromAddr] == 0 && command != "version" &&  command != "verack"{
		return
	}

	if activePeers[fromAddr] != 0 {
		activePeers[fromAddr] = time.Now().Unix()
	}

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

	initPeerDiscovery(bc)

	go connnectionCheck(bc)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Net.Panic(err)
		}
		go handleConnection(conn, bc)
	}
}
