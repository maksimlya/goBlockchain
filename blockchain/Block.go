package blockchain

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"goBlockchain/datastructures"
	"goBlockchain/transactions"
	"goBlockchain/utility"
	"io"
	"strconv"
)

// Block header stores necessary information to verify the block's intergity
type Header struct {
	Index        int
	Timestamp    string
	Hash         string
	PreviousHash string
	Nonce        int
	MerkleRoot   string
	BloomFilter  string
}

// Full nodes hold the whole blocks including merkle trees
type Block struct {
	BlockHeader Header
	merkleTree  *datastructures.MerkleTree // Pointer to the block's merkle tree
}

/*==================================================================Miner Functions=============================================================================*/

// Mines initial genesis block in the blockchain
func MineGenesisBlock() Block {
	tStamp := utility.Time()

	// Genesis block constructor
	b := Block{BlockHeader: Header{Timestamp: tStamp, Hash: utility.Hash(tStamp), PreviousHash: strconv.Itoa(0), MerkleRoot: "Genesis Block"}}
	return b
}

// Provided a slice of transactions, will attempt to mine the next block in chain
func MineBlock(id int, difficulty int, previousHash string, txs []transactions.Transaction) Block {
	var hash string
	isValid := false

	tStamp := utility.Time()
	nonce := 0                               // Starting with nonce of 0
	merkle, _ := datastructures.NewTree(txs) // Build merkle tree based on provided transactions
	merkleRoot := ""
	if merkle != nil {
		merkleRoot = merkle.Root.Hash // If merkle tree isn't nil ( it has transaction ), then store the merkle root hash
	}
	bloomFilter := ""
	if len(txs) > 0 {
		bloomFilter = datastructures.CreateBloom(txs) // Create bloom filter
	}

	for !isValid {
		hash = utility.Hash(strconv.Itoa(id) + tStamp + previousHash + merkleRoot + bloomFilter + strconv.Itoa(nonce))
		// Checks whether the current nonce provides correct hash, and keeps increasing it until is proper
		isValid = ValidateHash(hash, difficulty)
		if !isValid {
			nonce++
		}
	}
	// Once the mining completed Construct the proper block and return it to the blockchain
	b := Block{BlockHeader: Header{Index: id, Timestamp: tStamp, Hash: hash, PreviousHash: previousHash, Nonce: nonce, MerkleRoot: merkleRoot, BloomFilter: bloomFilter}, merkleTree: merkle}
	return b
}

/*===========================================================================================================================================================================================*/

/*=======================================================Getter functions for various block's members==================*/
func (b *Block) GetId() int {
	return b.BlockHeader.Index
}
func (b Block) GetHash() string {
	return b.BlockHeader.Hash
}
func (b *Block) GetPreviousHash() string {
	return b.BlockHeader.PreviousHash
}
func (b *Block) GetNonce() int {
	return b.BlockHeader.Nonce
}
func (b *Block) GetTimestamp() string {
	return b.BlockHeader.Timestamp
}
func (b *Block) GetMerkleRoot() string {
	return b.BlockHeader.MerkleRoot
}
func (b *Block) GetMerkleTree() *datastructures.MerkleTree {
	return b.merkleTree
}
func (b *Block) GetTransactions() []transactions.Transaction {
	if b.merkleTree != nil {
		return b.merkleTree.GetTransactions()
	}
	return nil
}
func (b *Block) GetBloomFilter() string {
	return b.BlockHeader.BloomFilter
}

/*====================================================================================================================*/

/*===========================================================Various Validation Checks===========================================*/

// Function to check if the block's hash is proper, and no data was altered since it's creation
func (b *Block) ValidateBlock() bool {

	if b.GetId() == 0 { // Genesis block is not checked, since it contains no data and always proper
		return true
	}
	// Re-calculate the hash for current data in the block, and compare it to stored block's hash
	hash := utility.Hash(strconv.Itoa(b.GetId()) + b.GetTimestamp() + b.GetPreviousHash() + b.GetMerkleRoot() + b.GetBloomFilter() + strconv.Itoa(b.GetNonce()))

	if b.GetHash() == hash {
		return true
	}
	return false
}

// Function checks block's validity by hashing all it's members and checking
// if the difficulty result was reached.
func ValidateHash(hash string, diff int) bool {
	hexaData, _ := hex.DecodeString(hash)
	var binaryString string
	// Convert hash into binary representation and later check the amount of leading 0's
	// If it matches the difficulty value, then hash is valid and the block is found.
	for i := range hexaData {
		binaryString += utility.Hex2Bin(hexaData[i])
	}
	binaryData := []byte(binaryString[0:diff])

	for _, data := range binaryData {
		if data != 48 { // 48 == "0" In ASCII alphabet
			return false
		}
	}
	return true
}

// Function that checks if a given txHash possibly can belong to the block,
// thus reducing time in searching for the transaction deep in merkle trees
func (b *Block) CheckBloomFilter(txHash string) bool {
	if datastructures.CheckExist(txHash, b.BlockHeader.BloomFilter) {
		return true
	}
	return false
}

/*============================================================================================================================*/

/*=======================================================Serialization=======================================================*/

// Serialization needed to store blocks data into database. Boltdb stores key/values in []byte format
func (b *Block) Serialize() []byte {
	var result []byte
	result = append(result, utility.Serialize(b)...)
	if b.merkleTree != nil {
		result = append(result, utility.Serialize(b.merkleTree.GetTransactions())...)
	}
	return result
}

// Deserialization reverses the process to get the actual object from a []byte array
func DeserializeBlock(d []byte) *Block {
	var block Block
	var txs []transactions.Transaction
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	err = decoder.Decode(&txs)

	if err != nil && err != io.EOF {
		fmt.Println(err)
	}
	// TODO - Better merkle tree serialization
	tree, _ := datastructures.NewTree(txs)

	block.merkleTree = tree
	return &block
}

/*====================================================================================================================*/

func (b *Block) String() string {
	s := ""
	s += "{\n"
	s += "Block Id: " + strconv.Itoa(b.GetId()) + "\n"
	s += "Block hash: " + b.GetHash() + "\n"
	s += "Previous Hash: " + b.GetPreviousHash() + "\n"
	s += "Nonce: " + strconv.Itoa(b.GetNonce()) + "\n"
	s += "Timestamp: " + b.GetTimestamp() + "\n"
	s += "Merkle Root: " + b.GetMerkleRoot() + "\n"
	s += "transactions: {\n"
	for _, tx := range b.GetTransactions() {
		s += tx.String()
	}
	s += "}\n"
	s += "};\n"

	return s
}
