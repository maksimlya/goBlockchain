package Blocks

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"goBlockchain/DataStructures"
	"goBlockchain/Transactions"
	"strconv"
	"time"
)

type Header struct {
	index        int
	timestamp    string
	hash         string
	previousHash string
	nonce        int
	merkleRoot   string
}
type Block struct {
	blockHeader Header
	merkleTree  *DataStructures.MerkleTree
}

var idx = 0

func MineGenesisBlock() Block {
	hasher := sha256.New()
	tStamp := time.Now().Format("02-01-2006 15:04:05")
	hasher.Write([]byte(tStamp))
	b := Block{blockHeader: Header{index: 0, timestamp: tStamp, hash: hex.EncodeToString(hasher.Sum(nil)), previousHash: strconv.Itoa(0), nonce: 0, merkleRoot: "Genesis Block"}, merkleTree: nil}
	return b
}

func MineBlock(difficulty int, previousHash string, txs []Transactions.Transaction) Block {
	tStamp := time.Now().Format("02-01-2006 15:04:05")
	nonce := 0
	merkle, _ := DataStructures.NewTree(txs)
	merkleRoot := ""
	if merkle != nil {
		merkleRoot = merkle.Root.Hash
	}

	hasher := sha256.New()
	hasher.Write([]byte(strconv.Itoa(idx + 1)))
	hasher.Write([]byte(tStamp))
	hasher.Write([]byte(previousHash))
	hasher.Write([]byte(merkleRoot))
	hasher.Write([]byte(strconv.Itoa(0)))
	hash := hex.EncodeToString(hasher.Sum(nil))
	isValid := ValidateHash(hash, difficulty)

	for !isValid {
		nonce++
		hasher := sha256.New()
		hasher.Write([]byte(strconv.Itoa(idx + 1)))
		hasher.Write([]byte(tStamp))
		hasher.Write([]byte(previousHash))
		hasher.Write([]byte(merkleRoot))
		hasher.Write([]byte(strconv.Itoa(nonce)))
		hash = hex.EncodeToString(hasher.Sum(nil))
		isValid = ValidateHash(hash, difficulty)
	}
	idx++
	b := Block{blockHeader: Header{index: idx, timestamp: tStamp, hash: hash, previousHash: previousHash, nonce: nonce, merkleRoot: merkleRoot}, merkleTree: merkle}
	return b
}

func (b Block) GetTransactions() []Transactions.Transaction {
	if b.merkleTree != nil {
		return b.merkleTree.GetTransactions()
	}
	return nil
}

func (b Block) PrintTime() {
	fmt.Println(b.blockHeader.timestamp)
}
func (b Block) PrintIdx() {
	fmt.Println(b.blockHeader.index)
}
func (b Block) GetId() int {
	return b.blockHeader.index
}
func (b Block) GetPreviousHash() string {
	return b.blockHeader.previousHash
}
func (b Block) GetNonce() int {
	return b.blockHeader.nonce
}
func (b Block) GetTimestamp() string {
	return b.blockHeader.timestamp
}
func (b Block) GetMerkleRoot() string {
	return b.blockHeader.merkleRoot
}
func (b Block) PrintHash() {
	fmt.Println(b.blockHeader.hash)
}

func ValidateHash(hash string, diff int) bool {
	checkStr := string(hash[0:diff])

	bytes := []byte(checkStr)
	j := 0

	for i := 0; i < len(bytes); i++ {
		if bytes[i] != 48 {
			return false
		}
		j++
	}

	return true
}

func (b Block) GetHash() string {
	return b.blockHeader.hash
}
