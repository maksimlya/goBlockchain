package Blocks

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"goBlockchain/DataStructures"
	"goBlockchain/Transactions"
	"io"
	"strconv"
	"time"
)

type Header struct {
	Index        int
	Timestamp    string
	Hash         string
	PreviousHash string
	Nonce        int
	MerkleRoot   string
	BloomFilter  string
}
type Block struct {
	BlockHeader Header
	merkleTree  *DataStructures.MerkleTree
}

func (b *Block) CheckBloomFilter(txHash string) bool {
	if DataStructures.CheckExist(txHash, b.BlockHeader.BloomFilter) {
		return true
	}
	return false
}

func MineGenesisBlock() Block {
	hasher := sha256.New()
	tStamp := time.Now().Format("02-01-2006 15:04:05")
	hasher.Write([]byte(tStamp))
	b := Block{BlockHeader: Header{Index: 0, Timestamp: tStamp, Hash: hex.EncodeToString(hasher.Sum(nil)), PreviousHash: strconv.Itoa(0), Nonce: 0, MerkleRoot: "Genesis Block", BloomFilter: ""}, merkleTree: nil}
	return b
}

func MineBlock(id int, difficulty int, previousHash string, txs []Transactions.Transaction) Block {
	tStamp := time.Now().Format("02-01-2006 15:04:05")
	nonce := 0
	merkle, _ := DataStructures.NewTree(txs)
	merkleRoot := ""
	if merkle != nil {
		merkleRoot = merkle.Root.Hash
	}
	bloomFilter := ""
	if len(txs) > 0 {
		bloomFilter = DataStructures.CreateBloom(txs)
	}

	hasher := sha256.New()
	hasher.Write([]byte(strconv.Itoa(id)))
	hasher.Write([]byte(tStamp))
	hasher.Write([]byte(previousHash))
	hasher.Write([]byte(merkleRoot))
	hasher.Write([]byte(bloomFilter))
	hasher.Write([]byte(strconv.Itoa(0)))
	hash := hex.EncodeToString(hasher.Sum(nil))
	isValid := ValidateHash(hash, difficulty)

	for !isValid {
		nonce++
		hasher := sha256.New()
		hasher.Write([]byte(strconv.Itoa(id)))
		hasher.Write([]byte(tStamp))
		hasher.Write([]byte(previousHash))
		hasher.Write([]byte(merkleRoot))
		hasher.Write([]byte(strconv.Itoa(nonce)))
		hash = hex.EncodeToString(hasher.Sum(nil))
		isValid = ValidateHash(hash, difficulty)
	}
	b := Block{BlockHeader: Header{Index: id, Timestamp: tStamp, Hash: hash, PreviousHash: previousHash, Nonce: nonce, MerkleRoot: merkleRoot, BloomFilter: bloomFilter}, merkleTree: merkle}
	return b
}

func (b Block) GetTransactions() []Transactions.Transaction {
	if b.merkleTree != nil {
		return b.merkleTree.GetTransactions()
	}
	return nil
}

func (b Block) PrintTime() {
	fmt.Println(b.BlockHeader.Timestamp)
}
func (b Block) PrintIdx() {
	fmt.Println(b.BlockHeader.Index)
}
func (b Block) GetId() int {
	return b.BlockHeader.Index
}
func (b Block) GetPreviousHash() string {
	return b.BlockHeader.PreviousHash
}
func (b Block) GetNonce() int {
	return b.BlockHeader.Nonce
}
func (b Block) GetTimestamp() string {
	return b.BlockHeader.Timestamp
}
func (b Block) GetMerkleRoot() string {
	return b.BlockHeader.MerkleRoot
}
func (b Block) PrintHash() {
	fmt.Println(b.BlockHeader.Hash)
}

func ValidateHash(hash string, diff int) bool {
	checkStr := string(hash[0:diff])

	hashBytes := []byte(checkStr)
	j := 0

	for i := 0; i < len(hashBytes); i++ {
		if hashBytes[i] != 48 {
			return false
		}
		j++
	}

	return true
}

func (b Block) GetHash() string {
	return b.BlockHeader.Hash
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if b.merkleTree != nil {
		err = encoder.Encode(b.merkleTree.GetTransactions())
		if err != nil {
			fmt.Println(err)
		}
	}

	return result.Bytes()
}

func DeserializeBlock(d []byte) *Block {
	var block Block
	var txs []Transactions.Transaction
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	err = decoder.Decode(&txs)

	if err != nil && err != io.EOF {
		fmt.Println(err)
	}

	tree, _ := DataStructures.NewTree(txs)

	block.merkleTree = tree

	return &block
}

func (b *Block) String() string {
	s := ""
	s += "{\n"
	s += "Block Id: " + strconv.Itoa(b.GetId()) + "\n"
	s += "Block hash: " + b.GetHash() + "\n"
	s += "Previous Hash: " + b.GetPreviousHash() + "\n"
	s += "Nonce: " + strconv.Itoa(b.GetNonce()) + "\n"
	s += "Timestamp: " + b.GetTimestamp() + "\n"
	s += "Merkle Root: " + b.GetMerkleRoot() + "\n"
	s += "Transactions: {\n"
	for _, tx := range b.GetTransactions() {
		s += tx.String()
	}
	s += "}\n"
	s += "};\n"

	return s
}
