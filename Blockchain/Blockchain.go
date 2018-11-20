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

var maxSizeOfTx = 4

func InitBlockchain() *Blockchain {
	var tip []byte
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

	bc := Blockchain{tip: tip, db: db, difficulty: 4, signatures: make(map[string]string)}
	return &bc
}

func (bc Blockchain) getSignature(key string) string {
	return bc.signatures[key]
}

func (bc Blockchain) GetSignature(txid string) string {
	sigData := ""
	err := bc.db.View(func(tx *bolt.Tx) error {
		s := tx.Bucket([]byte("signatures"))
		sigData = string(s.Get([]byte(txid)))

		return nil
	})
	if err != nil {
		fmt.Println(err)
	}

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

func (bc Blockchain) GetBlockById(id int) Blocks.Block {
	for i := 0; i < len(bc.chain); i++ {
		if bc.chain[i].GetId() == id {
			return bc.chain[i]
		}
	}
	return Blocks.Block{}
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
	it := bc.Iterator()
	s := ""
	for {
		block := it.Next()
		s += "{\n"
		s += "Block Id: " + strconv.Itoa(block.GetId()) + "\n"
		s += "Block hash: " + block.GetHash() + "\n"
		s += "Previous Hash: " + block.GetPreviousHash() + "\n"
		s += "Nonce: " + strconv.Itoa(block.GetNonce()) + "\n"
		s += "Timestamp: " + block.GetTimestamp() + "\n"
		s += "Merkle Root: " + block.GetMerkleRoot() + "\n"
		s += "Transactions: {\n"
		for _, j := range block.GetTransactions() {
			s += j.String()
		}
		s += "}\n"
		s += "};\n"

		if block.GetPreviousHash() == "0" {
			break
		}
	}
	return s
}

func (bc *Blockchain) GetBalanceForAddress(address string) int {
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

	return amount

}
