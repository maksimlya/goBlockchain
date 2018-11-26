package database

import (
	"fmt"
	"goBlockchain/imports/bolt"
	"strconv"
	"sync"
)

var instance *Database // Instance of this database object
var once sync.Once     // To sync goroutines at critical parts

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

func GetInstance() *Database {
	once.Do(func() {
		instance = GetDatabase()
	})
	return instance
}

func (d *Database) GetBlockById(blockId int) []byte {
	d.db, _ = bolt.Open("Blockchain.db", 0600, nil)
	var blockData []byte
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocks"))
		blockHash := b.Get([]byte(strconv.Itoa(blockId)))
		data := b.Get(blockHash)
		blockData = make([]byte, len(data))
		copy(blockData, data)

		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	d.db.Close()
	return blockData
}
func (d *Database) GetBlockByHash(blockHash string) []byte {
	d.db, _ = bolt.Open("Blockchain.db", 0600, nil)
	var blockData []byte
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocks"))
		data := b.Get([]byte(blockHash))
		blockData = make([]byte, len(data))
		copy(blockData, data)

		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	d.db.Close()
	return blockData
}

func (d *Database) GetLastBlock() []byte {
	var blockData []byte
	d.db, _ = bolt.Open("Blockchain.db", 0600, nil)
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocks"))

		lastHash := b.Get([]byte("l"))
		data := b.Get(lastHash)
		blockData = make([]byte, len(data))
		copy(blockData, data)

		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	d.db.Close()

	return blockData
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

func (d *Database) StoreBlock(blockHash string, blockId int, blockData []byte) {
	d.db, _ = bolt.Open("Blockchain.db", 0600, nil)
	err := d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocks"))
		err := b.Put([]byte(blockHash), blockData)
		err = b.Put([]byte(strconv.Itoa(blockId)), []byte(blockHash))
		err = b.Put([]byte("l"), []byte(blockHash))
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

func GetDatabase() *Database {
	db, _ := bolt.Open("Blockchain.db", 0600, nil)
	db.Close()
	return &Database{db: db}
}

func IsBlockchainExists() bool {
	exists := true
	db, err := bolt.Open("Blockchain.db", 0600, nil)

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocks"))
		if b == nil {
			exists = false
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	db.Close()
	return exists
}

func (d *Database) StoreNewBlockchain(lastHash string, lastId int, blockData []byte) {
	d.db, _ = bolt.Open("Blockchain.db", 0600, nil)
	err := d.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte("blocks"))
		s, err := tx.CreateBucket([]byte("signatures"))
		err = s.Put([]byte("Genesis"), []byte("0"))
		err = b.Put([]byte(lastHash), blockData)
		err = b.Put([]byte("l"), []byte(lastHash))
		err = b.Put([]byte("0"), []byte(lastHash))

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
