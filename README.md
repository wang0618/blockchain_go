# Blockchain in Go

A toy blockchain implementation in Go:

Base on [Jeiwan/blockchain_go](https://github.com/Jeiwan/blockchain_go/)


## 教程译文

转载自 [liuchengxu/blockchain-tutorial](https://github.com/liuchengxu/blockchain-tutorial)

* [基本原型](https://github.com/liuchengxu/blockchain-tutorial/blob/master/content/part-1/basic-prototype.md)
* [工作量证明](https://github.com/liuchengxu/blockchain-tutorial/blob/master/content/part-2/proof-of-work.md)
* [持久化和命令行接口](https://github.com/liuchengxu/blockchain-tutorial/blob/master/content/part-3/persistence-and-cli.md)
* [交易（1）](https://github.com/liuchengxu/blockchain-tutorial/blob/master/content/part-4/transactions-1.md)
* [地址](https://github.com/liuchengxu/blockchain-tutorial/blob/master/content/part-5/address.md)
* [交易（2）](https://github.com/liuchengxu/blockchain-tutorial/blob/master/content/part-6/transactions-2.md)
* [网络](https://github.com/liuchengxu/blockchain-tutorial/blob/master/content/part-7/network.md)

## 运行

```
go get github.com/1746199054/blockchain_go
cd $GOPATH/github.com/1746199054/blockchain_go
go build
NODE_ID=3000 ./blockchain_go createwallet  # You can use this command create many wallets
NODE_ID=3000 ./blockchain_go startnode -webui &
NODE_ID=3001 ./blockchain_go startnode -miner 168C3RJbprmpnxNry49ftjWGfFGQTNeDsU  # to genesis transaction output address 
```

### Usage

Run `./blockchain_go` to get this help:

```
Usage:
  createwallet - Generates a new key-pair and saves it into the wallet file
  getbalance -address ADDRESS - Get balance of ADDRESS
  listaddresses - Lists all addresses from the wallet file
  printchain - Print all the blocks of the blockchain
  reindexutxo - Rebuilds the UTXO set
  send -from FROM -to TO -amount AMOUNT -mine - Send AMOUNT of coins from FROM address to TO. Mine on the same node, when -mine is set.
  startnode -miner ADDRESS - Start a node with ID specified in NODE_ID env. var. -miner enables mining
  creategenesisblock -delta SECOND - Create a genesis block which timestamp is now with given 区块生成间隔
```

### Screenshot

![Screenshot](https://user-images.githubusercontent.com/8682073/56465515-5a241500-6431-11e9-937d-44e675ac5c03.png)
