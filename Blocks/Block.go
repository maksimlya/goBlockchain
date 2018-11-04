package Blocks

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

type Block struct {
	index			int
	timestamp 		int64
	hash			string
	previousHash	string
	nonce			int
}

var x = 0;

//func New(timestamp int, previousHash int, nonce int) block {
//	b := block {x, time.Now().UnixNano(), 0,0}
//	x++
//	return b
//}

func MineBlock(difficulty int, previousHash string) Block{
	tstamp := time.Now().UnixNano()
	nonce := 0
	hasher := sha256.New()
	hasher.Write([]byte (strconv.Itoa(x+1)))
	hasher.Write([]byte (strconv.FormatInt(tstamp,10)))
	hasher.Write([]byte (previousHash))
	hasher.Write([]byte (strconv.Itoa(0)))
	hash := hex.EncodeToString(hasher.Sum(nil))
	isValid := ValidateHash(hash,difficulty)

	for !isValid {
		nonce++
		hasher := sha256.New()
		hasher.Write([]byte (strconv.Itoa(x+1)))
		hasher.Write([]byte (strconv.FormatInt(tstamp,10)))
		hasher.Write([]byte (previousHash))
		hasher.Write([]byte (strconv.Itoa(nonce)))
		hash = hex.EncodeToString(hasher.Sum(nil))
		isValid = ValidateHash(hash,difficulty)
	}
	x++
	b:= Block{index:x,timestamp:tstamp,hash: hash,previousHash:previousHash,nonce:nonce}
	return b;
	}

func MineGenesisBlock() Block{
	hasher := sha256.New()
	tstamp := time.Now().UnixNano()
	hasher.Write([]byte (strconv.FormatInt(tstamp,10)))
	b := Block{index:0,timestamp:tstamp,hash: hex.EncodeToString(hasher.Sum(nil)),previousHash:strconv.Itoa(0),nonce:0}
	return b
	}

func (b Block) PrintTime(){
	fmt.Println(b.timestamp)
}
func (b Block) PrintIdx(){
	fmt.Println(b.index)
}
func (b Block) PrintHash(){
	fmt.Println(b.hash)
}

func ValidateHash(hash string, diff int) bool{
	checkStr := string(hash[0:diff])

	bytes := []byte(checkStr)
	j := 0;

	for i:=0 ; i < len(bytes) ; i++ {
		if(bytes[i] != 48){
			return false;
		}
		j++;
	}

	return true;
}

func (b Block) GetHash() string{
	return b.hash
}