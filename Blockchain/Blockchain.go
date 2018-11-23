package Blockchain

import (
	"fmt"
	"github.com/boltdb/bolt"
	"goBlockchain/Blocks"
	"goBlockchain/Security"
	"goBlockchain/Transactions"
	"strconv"
)

type Blockchain struct {
	chain       []Blocks.Block
	tip         []byte
	lastId      int
	db          *bolt.DB
	numOfBlocks int
	difficulty  int
	pendingTx   []Transactions.Transaction
	signatures  map[string]string
}

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}
type BlockchainForwardIterator struct {
	currentId []byte
	db        *bolt.DB
}

var maxSizeOfTx = 4

func InitBlockchain() *Blockchain {
	var tip []byte
	var lastId int

	db, err := bolt.Open("Blockchain.db", 0600, nil)

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocks"))

		if b == nil {

			genesis := Blocks.MineGenesisBlock()
			b, err := tx.CreateBucket([]byte("blocks"))
			s, err := tx.CreateBucket([]byte("signatures"))
			err = s.Put([]byte("Genesis"), []byte("0"))
			err = b.Put([]byte(genesis.GetHash()), genesis.Serialize())
			err = b.Put([]byte("l"), []byte(genesis.GetHash()))
			err = b.Put([]byte("0"), []byte(genesis.GetHash()))
			tip = []byte(genesis.GetHash())
			lastId = 0

			if err != nil {
				fmt.Println(err)
			}

		} else {
			tip = b.Get([]byte("l"))
			lastId = Blocks.DeserializeBlock(b.Get(tip)).GetId()
		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	bc := Blockchain{lastId: lastId, tip: tip, db: db, difficulty: 4, signatures: make(map[string]string)}
	db.Close()
	return &bc

}

func (bc Blockchain) getSignature(key string) string {
	return bc.signatures[key]
}

func (bc Blockchain) GetSignature(txid string) string {
	sigData := ""
	bc.db, _ = bolt.Open("Blockchain.db", 0600, nil)
	err := bc.db.View(func(tx *bolt.Tx) error {
		s := tx.Bucket([]byte("signatures"))
		sigData = string(s.Get([]byte(txid)))

		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	bc.db.Close()
	return sigData
}

func (bc Blockchain) InsertToChain(block Blocks.Block) {
	var remainingTx []Transactions.Transaction
	for i := 0; i < len(bc.pendingTx); i += maxSizeOfTx {
		remainingTx = bc.pendingTx[i : i+maxSizeOfTx]
		bc.chain = append(bc.chain, Blocks.MineBlock(0, bc.difficulty, bc.GetLastBlock().GetHash(), remainingTx))
		bc.numOfBlocks++
	}
}

func (bc *Blockchain) MineNextBlock() {
	var lastHash []byte
	var lastId int
	bc.db, _ = bolt.Open("Blockchain.db", 0600, nil)
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
		if bc.pendingTx[i].IsNil() {
			amountOfTx++
			continue
		}
		if !Security.VerifySignature(bc.signatures[bc.pendingTx[i].GetId()], bc.pendingTx[i].GetId(), bc.pendingTx[i].GetSender()) {
			bc.pendingTx[i] = Transactions.GetNil()
		}
		transactios = append(transactios, bc.pendingTx[i])
		amountOfTx++
	}
	bc.pendingTx = bc.pendingTx[amountOfTx:] // TODO -  improve for dynamic use
	newBlock := Blocks.MineBlock(lastId+1, bc.difficulty, string(lastHash), transactios)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocks"))
		s := tx.Bucket([]byte("signatures"))
		err := b.Put([]byte(newBlock.GetHash()), newBlock.Serialize())
		err = b.Put([]byte(strconv.Itoa(newBlock.GetId())), []byte(newBlock.GetHash()))
		for i := range transactios {
			err := s.Put([]byte(transactios[i].GetId()), []byte(bc.signatures[transactios[i].GetId()]))
			if err != nil {
				fmt.Println(err)
			}
		}
		err = b.Put([]byte("l"), []byte(newBlock.GetHash()))
		bc.tip = []byte(newBlock.GetHash())

		if err != nil {
			fmt.Println(err)
		}

		return nil
	})
	bc.db.Close()
}

func (bc Blockchain) SearchBlock(hash string) Blocks.Block {
	var b Blocks.Block
	for i := 0; i < bc.numOfBlocks; i++ {
		if hash == bc.chain[i].GetHash() {
			return bc.chain[i]
		}
	}
	return b
}

func (bc Blockchain) GetBlockById(id int) *Blocks.Block {
	bc.db, _ = bolt.Open("Blockchain.db", 0600, nil)
	var block *Blocks.Block
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocks"))
		blockHash := b.Get([]byte(strconv.Itoa(id)))
		encodedBlock := b.Get(blockHash)
		block = Blocks.DeserializeBlock(encodedBlock)

		return nil
	})
	if err != nil {
		fmt.Println(err)
	}

	bc.db.Close()
	return block
}

func (bc *Blockchain) AddTransaction(transaction Transactions.Transaction, signature string) {
	if !Security.VerifySignature(signature, transaction.GetId(), transaction.GetSender()) {
		return
	}
	bc.signatures[transaction.GetId()] = signature
	bc.pendingTx = append(bc.pendingTx, transaction)
}

func (bc *Blockchain) GetLastBlock() *Blocks.Block {
	it := bc.Iterator()
	return it.Next()
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}

	return bci
}
func (bc *Blockchain) ForwardIterator() *BlockchainForwardIterator {
	bci := &BlockchainForwardIterator{[]byte("0"), bc.db}

	return bci
}

func (i *BlockchainIterator) Next() *Blocks.Block {
	var block *Blocks.Block
	i.db, _ = bolt.Open("Blockchain.db", 0600, nil)
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
	i.db.Close()
	return block
}
func (i *BlockchainForwardIterator) Next() *Blocks.Block {
	var block *Blocks.Block
	i.db, _ = bolt.Open("Blockchain.db", 0600, nil)
	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocks"))
		blockHash := b.Get(i.currentId)
		encodedBlock := b.Get(blockHash)
		block = Blocks.DeserializeBlock(encodedBlock)

		return nil
	})

	i.currentId = []byte(strconv.Itoa(block.GetId() + 1))

	if err != nil {
		fmt.Println(err)
	}
	i.db.Close()
	return block
}

//func (bc *Blockchain) String() string {
//	it := bc.Iterator()
//	s := ""
//	for {
//		block := it.Next()
//		s += "{\n"
//		s += "Block Id: " + strconv.Itoa(block.GetId()) + "\n"
//		s += "Block hash: " + block.GetHash() + "\n"
//		s += "Previous Hash: " + block.GetPreviousHash() + "\n"
//		s += "Nonce: " + strconv.Itoa(block.GetNonce()) + "\n"
//		s += "Timestamp: " + block.GetTimestamp() + "\n"
//		s += "Merkle Root: " + block.GetMerkleRoot() + "\n"
//		s += "Transactions: {\n"
//		for _, j := range block.GetTransactions() {
//			s += j.String()
//		}
//		s += "}\n"
//		s += "};\n"
//
//		if block.GetPreviousHash() == "0" {
//			break
//		}
//	}
//	return s
//}

func (bc *Blockchain) ValidateChain() bool {
	it := bc.Iterator()
	previousIterator := bc.Iterator()
	previousIterator.Next()

	for {
		block := it.Next()

		if block.GetId() == 0 {
			return true
		}

		prevBlock := previousIterator.Next()

		if block.GetId() != prevBlock.GetId()+1 {
			return false
		}
		if block.GetPreviousHash() != prevBlock.GetHash() {
			return false
		}

		if !block.ValidateBlock() {
			return false
		}
		if block.GetPreviousHash() == "0" {
			break
		}
	}
	return true
}

func (bc *Blockchain) TraverseBlockchain() []*Blocks.Block {
	var blocks []*Blocks.Block
	it := bc.Iterator()
	for {
		block := it.Next()
		blocks = append(blocks, block)
		if block.GetPreviousHash() == "0" {
			break
		}
	}
	return blocks
}
func (bc *Blockchain) TraverseForwardBlockchain() []*Blocks.Block {

	var blocks []*Blocks.Block
	it := bc.ForwardIterator()
	for {
		block := it.Next()
		blocks = append(blocks, block)

		if block.GetId() == bc.lastId {
			break
		}
	}
	return blocks
}

// TODO - rework function for polls use ( aka check balance for a given poll tag )
func (bc *Blockchain) GetBalanceForAddress(address string) int {
	bc.db, _ = bolt.Open("Blockchain.db", 0600, nil)
	var amount = 0
	it := bc.Iterator()
	err := bc.db.View(func(tx *bolt.Tx) error {
		s := tx.Bucket([]byte("signatures"))

		for {
			block := it.Next()
			for _, element := range block.GetTransactions() {
				signature := string(s.Get([]byte(element.GetId())))
				if element.GetReceiver() == address {
					if Security.VerifySignature(signature, element.GetId(), element.GetSender()) {
						amount += element.GetAmount()
					}
				}
				if element.GetSender() == address {
					if Security.VerifySignature(signature, element.GetId(), element.GetSender()) {
						amount -= element.GetAmount()
					}
				}
			}

			if block.GetPreviousHash() == "0" {
				break
			}
		}

		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	bc.db.Close()
	return amount

}
