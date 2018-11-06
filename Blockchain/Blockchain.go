package Blockchain

import (
	"goBlockchain/Blocks"
)

type Blockchain struct {
	chain       []Blocks.Block
	numOfblocks int
}

//getlastblock

func InitBlockchain(block Blocks.Block) Blockchain {
	bc := Blockchain{chain: []Blocks.Block{}, numOfblocks: 0}
	bc.chain = append(bc.chain, block.GetGenisis())
	bc.numOfblocks++
	return bc
}

func (bc Blockchain) InsertToChain(block Blocks.Block) {
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
