package webui

import (
	"blockchain_go/blockchain"
	"blockchain_go/log"
	"blockchain_go/net"
	"blockchain_go/utils"
	"blockchain_go/wallet"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	olog "log"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type msg struct {
	Status bool
	Data   interface{}
}

const maxBlocksCount = 10 // 一次最大允许返回的区块数量
var webUIAddr string;

var nodeID string

var upgrader = websocket.Upgrader{} // use default options

func init() {
	port, err := utils.GetAvailablePort()
	if err != nil {
		olog.Panic("upgrade:", err)
	}
	webUIAddr = fmt.Sprintf("127.0.0.1:%d", port)

}

func StartWebUIServer(htmlPath string, nodeId_ string) {
	nodeID = nodeId_

	fileHandler := http.FileServer(http.Dir(htmlPath))
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/send", sendHandler)
	http.HandleFunc("/log", getlogHandler)
	http.Handle("/", fileHandler)

	go utils.OpenBrowser("http://" + webUIAddr)
	olog.Fatal(http.ListenAndServe(webUIAddr, nil))
}

func getlogHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		olog.Print("upgrade:", err)
		return
	}
	defer c.Close()

	c.WriteMessage(websocket.TextMessage, []byte("Connect to server\n"))

	for {
		msg, _ := <-log.LogChan
		err = c.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			olog.Println("write:", err)
			break
		}
	}
}

// status 返回客户端状态
func statusHandler(w http.ResponseWriter, r *http.Request) {
	type statusMsg struct {
		Wallet []struct {
			Addr    string
			Balance int
		}
		Peers []struct {
			Addr       string
			LastActive string
		}
		Blocks   []blockchain.Block
		NodeAddr string
	}

	setJsonHeader(w)

	var res statusMsg

	res.NodeAddr = net.GetNodeAddr()

	bc := blockchain.GetBlockchain()
	UTXOSet := blockchain.UTXOSet{bc}
	wallets, _ := wallet.NewWallets(nodeID)
	addresses := wallets.GetAddresses()
	for _, address := range addresses {
		pubKeyHash := utils.Base58Decode([]byte(address))
		pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
		UTXOs := UTXOSet.FindUTXO(pubKeyHash)

		balance := 0
		for _, out := range UTXOs {
			balance += out.Value
		}
		res.Wallet = append(res.Wallet, struct {
			Addr    string
			Balance int
		}{Addr: address, Balance: balance})
	}
	sort.Slice(res.Wallet, func(i, j int) bool {
		return res.Wallet[i].Balance > res.Wallet[j].Balance || (res.Wallet[i].Balance == res.Wallet[j].Balance && res.Wallet[i].Addr >= res.Wallet[j].Addr)
	})

	for addr, st := range net.GetActivePeers() {
		tm := time.Unix(st, 0)
		res.Peers = append(res.Peers, struct {
			Addr       string
			LastActive string
		}{Addr: addr, LastActive: tm.Format("15:04:05")})
	}
	sort.Slice(res.Peers, func(i, j int) bool {
		return res.Peers[i].LastActive < res.Peers[j].LastActive
	})

	bci := bc.Iterator()
	for {
		block := bci.Next()

		res.Blocks = append(res.Blocks, *block)

		if len(block.PrevBlockHash) == 0 || len(res.Blocks) > maxBlocksCount {
			break
		}
	}

	setRespose(w, msg{true, res})
}

// send 发起交易
func sendHandler(w http.ResponseWriter, r *http.Request) {
	setJsonHeader(w)
	r.ParseForm()
	if r.Form["amount"] == nil || r.Form["to"] == nil || r.Form["from"] == nil {
		log.MISC.Println("ERROR: Can't get 'amount' or 'to' or 'from' field in sendHandler")
		setRespose(w, msg{false, "缺少参数！"})
		return
	}

	from := r.Form["from"][0]
	to := r.Form["to"][0]
	amount, _ := strconv.Atoi(r.Form["amount"][0])

	if !wallet.ValidateAddress(from) || !wallet.ValidateAddress(to) {
		setRespose(w, msg{false, "钱包地址不合法！"})
		return
	}

	bc := blockchain.GetBlockchain()
	UTXOSet := blockchain.UTXOSet{bc}

	wallets, err := wallet.NewWallets(nodeID)
	if err != nil {
		log.MISC.Println("ERROR: call wallet.NewWallets() failed")
		setRespose(w, msg{false, "获取钱包失败！"})
		return
	}
	wallet_ := wallets.GetWallet(from)

	tx := wallet.NewUTXOTransaction(&wallet_, to, amount, &UTXOSet)

	UTXOSet.UpdateForTx(tx, false)
	net.MemPool[hex.EncodeToString(tx.ID)] = *tx

	net.BroadcastTx(tx)
	if err != nil {
		log.MISC.Println("ERROR: call net.SendTx() failed")
		setRespose(w, msg{false, "发送交易消息失败！"})
		return
	}

	log.MISC.Printf("%s => %s\n %dBTC\n发送交易消息成功，等待交易确认中...\n", from, to, amount)
	setRespose(w, msg{true, "发送交易消息成功，等待交易确认中..."})
}

func setJsonHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json") // normal header
	w.WriteHeader(http.StatusOK)
}

func setRespose(w http.ResponseWriter, msg_ msg) {
	data, _ := json.Marshal(msg_)
	w.Write(data)
}
