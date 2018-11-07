package Blockchain

import (
	"goBlockchain/Blocks"
	"goBlockchain/Transactions"
)

type Blockchain struct {
	chain       []Blocks.Block
	numOfblocks int
	difficulty  int
	pendingTx   []Transactions.Transaction
}

var maxSizeOfTx = 8

func InitBlockchain(block Blocks.Block) Blockchain {
	bc := Blockchain{chain: []Blocks.Block{}, numOfblocks: 0, difficulty: 4, pendingTx: []Transactions.Transaction{}}
	bc.chain = append(bc.chain, Blocks.MineGenesisBlock())
	bc.numOfblocks++
	return bc
}

func (bc Blockchain) InsertToChain(block Blocks.Block) {
	var remainingTx []Transactions.Transaction
	for i := 0; i < len(bc.pendingTx); i += maxSizeOfTx {
		remainingTx = bc.pendingTx[i : i+maxSizeOfTx]
		bc.chain = append(bc.chain, Blocks.MineBlock(bc.difficulty, bc.GetLastBlock().GetHash(), remainingTx))
		bc.numOfblocks++
	}
}

func (bc Blockchain) SearchBlock(hash string) Blocks.Block {
	var b Blocks.Block
	for i := 0; i < bc.numOfblocks; i++ {
		if hash == bc.chain[i].GetHash() {
			return bc.chain[i]
		}
	}
	return b
}

func (bc Blockchain) GetLastBlock() Blocks.Block {
	return bc.chain[bc.numOfblocks-1]
}
