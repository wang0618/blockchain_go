package cli

import (
	"blockchain_go/blockchain"
	"blockchain_go/miner"
	"blockchain_go/transaction"
	"blockchain_go/utils"
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"time"
)

// createGenesisBlock 创建一个创世区块, 可以根据传入的区块创建间隔生成合适的难度值，并以此难度值计算出需要硬编码进程序的创世区块的一些字段
// delta_sec为区块链平均出块时间（秒）
func (cli *CLI) createGenesisBlock(delta_sec int) {
	const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"
	const genesisCoinbaseAddress = "168C3RJbprmpnxNry49ftjWGfFGQTNeDsU"

	coinbase := transaction.NewCoinbaseTX(genesisCoinbaseAddress, genesisCoinbaseData)
	ts := time.Now().Unix()
	block := &blockchain.Block{ts, []*transaction.Transaction{coinbase}, []byte{}, []byte{}, 0, 0, make([]byte, 32)}
	difficulty := getDifficult(block, delta_sec)
	block.Difficulty = difficulty

	fmt.Printf("Difficulty: %x\n\n", difficulty)
	fmt.Printf("Start proof of work\n")

	pow := miner.NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce

	fmt.Printf("============== Genesis Block Info ====================\n")
	fmt.Printf("CoinbaseData: %s\n", genesisCoinbaseData)
	fmt.Printf("CoinbaseAddress: %s\n", genesisCoinbaseAddress)
	fmt.Printf("Hash: %x\n", block.Hash)
	fmt.Printf("Timestamp: %d\n", block.Timestamp)
	fmt.Printf("Difficulty: %x\n", block.Difficulty)
	fmt.Printf("Nonce: %d\n", block.Nonce)
}

// getDifficult 返回可以在平均sec秒内出块的难度值
func getDifficult(block *blockchain.Block, sec int) []byte {
	var hashInt, difficulty big.Int
	hashInt.Exp(big.NewInt(2), big.NewInt(256), nil) // hashInt = 2**256

	const repeat = 2
	cnt := mockPow(block, sec*repeat) / repeat
	fmt.Printf("Repeat time: %d\n", repeat)
	fmt.Printf("Avg hash per %d sec: %d\n", sec, cnt)

	difficulty.Div(&hashInt, big.NewInt(cnt))

	difficultyBytes := difficulty.Bytes()
	difficultyBytes = append(make([]byte, 32-len(difficultyBytes)), difficultyBytes...)

	return difficultyBytes
}

// mockPow 模拟挖矿运算，返回在sec秒内进行的哈希运算的次数
func mockPow(block *blockchain.Block, sec int) int64 {
	start_ts := time.Now().Unix()
	var cnt int64
	fmt.Printf("Running hash, please wait %d seconds\n", sec)
	for cnt < math.MaxInt64 {
		// 模拟 ProofOfWork.prepareData()
		data := bytes.Join(
			[][]byte{
				block.PrevBlockHash,
				block.HashTransactions(),
				utils.IntToHex(block.Timestamp),
				block.Difficulty,
				utils.IntToHex(int64(cnt)),
			},
			[]byte{},
		)

		sha256.Sum256(data)
		cnt++
		if cnt%100000 == 0 {
			if time.Now().Unix()-start_ts > int64(sec) {
				return cnt
			}
		}
	}
	return cnt
}
