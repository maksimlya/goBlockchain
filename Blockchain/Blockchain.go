package Blockchain

import (
	"goBlockchain/Blocks"
	"goBlockchain/Transactions"
)

type Blockchain struct {
	chain       []Blocks.Block
	numOfblocks int
	difficulty  int
	pendingTx	[]Transactions.Transaction
}

var maxSizeOfTx = 8

func InitBlockchain(block Blocks.Block) Blockchain {
	bc := Blockchain{chain: []Blocks.Block{}, numOfblocks: 0}
	bc.chain = append(bc.chain, block.GetGenisis())
	bc.numOfblocks++
	bc.difficulty = 4;
	return bc
}

func (bc Blockchain) InsertToChain(block Blocks.Block) {
	var transactios []Transactions.Transaction
	for
	bc.chain = append(bc.chain, Blocks.MineBlock(bc.difficulty, bc.GetLastBlock().GetHash(), ))
	bc.chain = append(bc.chain, block)
	bc.numOfblocks++
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
