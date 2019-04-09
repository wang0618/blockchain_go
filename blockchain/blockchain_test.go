package blockchain

import (
	"blockchain_go/transaction"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBlockchain_AddBlock_GetBlockHashes(t *testing.T) {
	genesis := newGenesisBlock()

	blocks := []*Block{genesis}

	for i := 0; i < 10; i++ {
		last := blocks[len(blocks)-1]
		b := NewBlock([]*transaction.Transaction{}, last.Hash, last.Height+1, last.Difficulty)
		b.Hash = make([]byte, 32)
		copy(b.Hash, last.Hash)
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

	bc.AddBlock(blocks[9])
	assert.Equal(t, bc.tip, blocks[9].Hash)

	hashs := bc.GetBlockHashes(blocks[2].Hash, 3)
	assert.Equal(t, hashs, [][]byte{blocks[3].Hash, blocks[4].Hash, blocks[5].Hash})

	hashs = bc.GetBlockHashes(blocks[0].Hash, 5)
	assert.Equal(t, hashs, [][]byte{blocks[1].Hash, blocks[2].Hash, blocks[3].Hash, blocks[4].Hash, blocks[5].Hash})

	hashs = bc.GetBlockHashes(blocks[8].Hash, 5)
	assert.Equal(t, hashs, [][]byte{blocks[9].Hash})

}
