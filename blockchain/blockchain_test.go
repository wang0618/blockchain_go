package blockchain

import (
	"blockchain_go/transaction"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBCAddBlock(t *testing.T) {
	genesis := newGenesisBlock()

	blocks := []*Block{genesis}

	for i := 0; i < 10; i++ {
		last := blocks[len(blocks)-1]
		b := NewBlock([]*transaction.Transaction{}, last.Hash, last.Height+1, last.Difficulty)
		b.Hash = make([]byte, 32)
		copy(b.Hash, genesis.Hash)
		b.Hash[0]++
		blocks = append(blocks, b)
	}

	bc := GetBlockchain()
	bc.AddBlock(genesis)
	assert.Equal(t, bc.tip, genesis.Hash)

	bc.AddBlock(blocks[1])
	bc.AddBlock(blocks[2])
	assert.Equal(t, bc.tip, blocks[2].Hash)

	bc.AddBlock(blocks[4])
	bc.AddBlock(blocks[5])
	bc.AddBlock(blocks[7])
	assert.Equal(t, bc.tip, blocks[2].Hash)

	bc.AddBlock(blocks[3])
	assert.Equal(t, bc.tip, blocks[5].Hash)

	bc.AddBlock(blocks[6])
	assert.Equal(t, bc.tip, blocks[7].Hash)

	bc.AddBlock(blocks[4])
	bc.AddBlock(blocks[5])
	bc.AddBlock(blocks[7])
	bc.AddBlock(blocks[8])
	bc.AddBlock(blocks[6])
	assert.Equal(t, bc.tip, blocks[8].Hash)
}
