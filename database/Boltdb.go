package database

import (
	"fmt"
	"goBlockchain/blockchain"
	"goBlockchain/imports/bolt"
	"strconv"
)

type Database struct {
	db *bolt.DB
}

func (d *Database) GetSignatureByHash(txHash string) string {
	sigData := ""
	d.db, _ = bolt.Open("Blockchain.db", 0600, nil)
	err := d.db.View(func(tx *bolt.Tx) error {
		s := tx.Bucket([]byte("signatures"))
		sigData = string(s.Get([]byte(txHash)))

		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	d.db.Close()

	return sigData
}

func (d *Database) GetBlockById(blockId int) *blockchain.Block {
	d.db, _ = bolt.Open("Blockchain.db", 0600, nil)
	var block *blockchain.Block
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocks"))
		blockHash := b.Get([]byte(strconv.Itoa(blockId)))
		blockData := b.Get(blockHash)
		block = blockchain.DeserializeBlock(blockData)

		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	d.db.Close()
	return block
}
func (d *Database) GetBlockByHash(blockHash string) *blockchain.Block {
	d.db, _ = bolt.Open("Blockchain.db", 0600, nil)
	var block *blockchain.Block
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocks"))
		blockHash := b.Get([]byte(blockHash))
		blockData := b.Get(blockHash)
		block = blockchain.DeserializeBlock(blockData)

		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	d.db.Close()
	return block
}

func (d *Database) GetLastBlock() *blockchain.Block {
	var lastBlock *blockchain.Block
	d.db, _ = bolt.Open("Blockchain.db", 0600, nil)
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocks"))
		lastHash := b.Get([]byte("l"))
		blockData := b.Get(lastHash)
		lastBlock = blockchain.DeserializeBlock(blockData)

		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	d.db.Close()

	return lastBlock
}
func (d *Database) GetLastBlockHash() string {
	var lastBlockHash string
	d.db, _ = bolt.Open("Blockchain.db", 0600, nil)
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocks"))
		lastHash := b.Get([]byte("l"))
		lastBlockHash = string(lastHash[:])

		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	d.db.Close()
	return lastBlockHash
}

func (d *Database) StoreBlock(block blockchain.Block) {
	d.db, _ = bolt.Open("Blockchain.db", 0600, nil)
	err := d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocks"))
		err := b.Put([]byte(block.GetHash()), block.Serialize())
		err = b.Put([]byte(strconv.Itoa(block.GetId())), []byte(block.GetHash()))
		err = b.Put([]byte("l"), []byte(block.GetHash()))
		//bc.tip = []byte(newBlock.GetHash())
		if err != nil {
			fmt.Println(err)
		}

		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	d.db.Close()
}

func (d *Database) StoreSignature(txHash string, signature string) {
	d.db, _ = bolt.Open("Blockchain.db", 0600, nil)
	err := d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("signatures"))
		err := b.Put([]byte(txHash), []byte(signature))
		if err != nil {
			fmt.Println(err)
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	d.db.Close()
}

func GetBlockchain() *blockchain.Blockchain {
	var lastHash string
	var lastId int
	db, err := bolt.Open("Blockchain.db", 0600, nil)

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocks"))

		if b == nil {
			genesis := blockchain.MineGenesisBlock()
			b, err := tx.CreateBucket([]byte("blocks"))
			s, err := tx.CreateBucket([]byte("signatures"))
			err = s.Put([]byte("Genesis"), []byte("0"))
			err = b.Put([]byte(genesis.GetHash()), genesis.Serialize())
			err = b.Put([]byte("l"), []byte(genesis.GetHash()))
			err = b.Put([]byte("0"), []byte(genesis.GetHash()))
			lastHash = genesis.GetHash()
			lastId = 0

			if err != nil {
				fmt.Println(err)
			}

		} else {
			lastHash = string(b.Get([]byte("l"))[:])
			lastId = blockchain.DeserializeBlock(b.Get([]byte(lastHash))).GetId()
		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	bc := blockchain.Blockchain{LastId: lastId, LastHash: lastHash, Db: Database{db: db}, Difficulty: 4, Signatures: make(map[string]string)}
	db.Close()
	return &bc
}
