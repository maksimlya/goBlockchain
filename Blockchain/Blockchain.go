package Blockchain

import (
	"goBlockchain/Blocks"
	"goBlockchain/Transactions"
	"strconv"
)

type Blockchain struct {
	chain       []Blocks.Block
	numOfblocks int
	difficulty  int
	pendingTx   []Transactions.Transaction
}

var maxSizeOfTx = 4

func InitBlockchain() Blockchain {
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

func (bc *Blockchain) MineNextBlock() {
	var transactios []Transactions.Transaction
	amountOfTx := 0
	for i := 0; i < maxSizeOfTx && i < len(bc.pendingTx); i++ {
		transactios = append(transactios, bc.pendingTx[i])
		amountOfTx++
	}
	bc.pendingTx = bc.pendingTx[amountOfTx:] // TODO -  improve for dynamic use
	bc.chain = append(bc.chain, Blocks.MineBlock(bc.difficulty, bc.GetLastBlock().GetHash(), transactios))
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

func (bc *Blockchain) AddTransaction(tx Transactions.Transaction) {
	bc.pendingTx = append(bc.pendingTx, tx)
}

func (bc Blockchain) GetLastBlock() Blocks.Block {
	return bc.chain[bc.numOfblocks-1]
}

func (bc *Blockchain) String() string {
	s := ""
	for _, l := range bc.chain {
		s += "{\n"
		s += "Block Id: " + strconv.Itoa(l.GetId()) + "\n"
		s += "Block hash: " + l.GetHash() + "\n"
		s += "Previous Hash: " + l.GetPreviousHash() + "\n"
		s += "Nonce: " + strconv.Itoa(l.GetNonce()) + "\n"
		s += "Timestamp: " + l.GetTimestamp() + "\n"
		s += "Merkle Root: " + l.GetMerkleRoot() + "\n"
		s += "Transactions: {\n"
		for _, j := range l.GetTransactions() {
			s += j.String()
		}
		s += "}\n"
		s += "};\n"
	}
	return s
}
