package Blockchain

import (
	"fmt"
	"github.com/boltdb/bolt"
	"goBlockchain/Blocks"
	"goBlockchain/Transactions"
	"strconv"
)

type Blockchain struct {
	chain       []Blocks.Block
	tip         []byte
	db          *bolt.DB
	numOfblocks int
	difficulty  int
	pendingTx   []Transactions.Transaction
}

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

var maxSizeOfTx = 4

func InitBlockchain() *Blockchain {
	var tip []byte
	db, err := bolt.Open("Blockchain.db", 0600, nil)

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocks"))

		if b == nil {

			genesis := Blocks.MineGenesisBlock()
			b, err := tx.CreateBucket([]byte("blocks"))
			err = b.Put([]byte(genesis.GetHash()), genesis.Serialize())
			err = b.Put([]byte("l"), []byte(genesis.GetHash()))
			err = b.Put([]byte("id"), []byte("0"))
			tip = []byte(genesis.GetHash())

			if err != nil {
				fmt.Println(err)
			}

		} else {
			tip = b.Get([]byte("l"))

		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	bc := Blockchain{tip: tip, db: db, difficulty: 4}
	//bc.chain = append(bc.chain, Blocks.MineGenesisBlock())
	//bc.numOfblocks++
	return &bc
}

func (bc Blockchain) InsertToChain(block Blocks.Block) {
	var remainingTx []Transactions.Transaction
	for i := 0; i < len(bc.pendingTx); i += maxSizeOfTx {
		remainingTx = bc.pendingTx[i : i+maxSizeOfTx]
		bc.chain = append(bc.chain, Blocks.MineBlock(0, bc.difficulty, bc.GetLastBlock().GetHash(), remainingTx))
		bc.numOfblocks++
	}
}

func (bc *Blockchain) MineNextBlock() {
	var lastHash []byte
	var lastId int
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocks"))
		lastHash = b.Get([]byte("l"))
		gg := b.Get(lastHash)
		lastBlock := Blocks.DeserializeBlock(gg)

		lastId = lastBlock.GetId()

		return nil
	})
	if err != nil {
		fmt.Println(err)
	}

	var transactios []Transactions.Transaction
	amountOfTx := 0
	for i := 0; i < maxSizeOfTx && i < len(bc.pendingTx); i++ {
		transactios = append(transactios, bc.pendingTx[i])
		amountOfTx++
	}
	bc.pendingTx = bc.pendingTx[amountOfTx:] // TODO -  improve for dynamic use
	newBlock := Blocks.MineBlock(lastId+1, bc.difficulty, string(lastHash), transactios)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocks"))
		err := b.Put([]byte(newBlock.GetHash()), newBlock.Serialize())
		err = b.Put([]byte("l"), []byte(newBlock.GetHash()))
		bc.tip = []byte(newBlock.GetHash())

		if err != nil {
			fmt.Println(err)
		}

		return nil
	})

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

func (bc Blockchain) GetBlockById(id int) Blocks.Block {
	for i := 0; i < len(bc.chain); i++ {
		if bc.chain[i].GetId() == id {
			return bc.chain[i]
		}
	}
	return Blocks.Block{}
}

func (bc *Blockchain) AddTransaction(tx Transactions.Transaction) {
	bc.pendingTx = append(bc.pendingTx, tx)
}

func (bc Blockchain) GetLastBlock() Blocks.Block {
	return bc.chain[bc.numOfblocks-1]
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}

	return bci
}

func (i *BlockchainIterator) Next() *Blocks.Block {
	var block *Blocks.Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocks"))
		encodedBlock := b.Get(i.currentHash)
		block = Blocks.DeserializeBlock(encodedBlock)

		return nil
	})

	i.currentHash = []byte(block.GetPreviousHash())

	if err != nil {
		fmt.Println(err)
	}

	return block
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
