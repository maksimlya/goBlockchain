package Blocks

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"goBlockchain/Transactions"
	"strconv"
	"time"
)

type Header struct {
	index        int
	timestamp    int64
	hash         string
	previousHash string
	nonce        int
	merkleRoot   string
}
type Block struct {
	blockHeader Header
	trans       []Transactions.Transaction
}

var idx = 0

func MineGenesisBlock() Block {
	hasher := sha256.New()
	tStamp := time.Now().UnixNano()
	hasher.Write([]byte(strconv.FormatInt(tStamp, 10)))
	b := Block{blockHeader: Header{index: 0, timestamp: tStamp, hash: hex.EncodeToString(hasher.Sum(nil)), previousHash: strconv.Itoa(0), nonce: 0, merkleRoot: ""}}
	return b
}

func MineBlock(difficulty int, previousHash string) Block {
	tStamp := time.Now().UnixNano()
	nonce := 0
	hasher := sha256.New()
	hasher.Write([]byte(strconv.Itoa(idx + 1)))
	hasher.Write([]byte(strconv.FormatInt(tStamp, 10)))
	hasher.Write([]byte(previousHash))
	hasher.Write([]byte(strconv.Itoa(0)))
	hash := hex.EncodeToString(hasher.Sum(nil))
	isValid := ValidateHash(hash, difficulty)

	for !isValid {
		nonce++
		hasher := sha256.New()
		hasher.Write([]byte(strconv.Itoa(idx + 1)))
		hasher.Write([]byte(strconv.FormatInt(tStamp, 10)))
		hasher.Write([]byte(previousHash))
		hasher.Write([]byte(strconv.Itoa(nonce)))
		hash = hex.EncodeToString(hasher.Sum(nil))
		isValid = ValidateHash(hash, difficulty)
	}
	idx++
	b := Block{index: idx, timestamp: tStamp, hash: hash, previousHash: previousHash, nonce: nonce}
	return b
}

func (b Block) PrintTime() {
	fmt.Println(b.timestamp)
}
func (b Block) PrintIdx() {
	fmt.Println(b.index)
}
func (b Block) PrintHash() {
	fmt.Println(b.hash)
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
	return b.hash
}
