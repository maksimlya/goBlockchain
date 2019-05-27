package blockchain

import (
	"fmt"
	"goBlockchain/database"
	"goBlockchain/p2p"
	"goBlockchain/transactions"
	"goBlockchain/utility"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	difficulty  int = 16   // Amount of 0's in binary format to aim for when mining a block.
	maxSizeOfTx int = 1000 // Amount of tx's per block.
)

var (
	nc           *p2p.NetworkController         // Network controller to interact with p2p system.
	instance     *Blockchain                    // Instance of this blockchain object
	once         sync.Once                      // To sync goroutines at critical parts
	pendingBlock bool                   = false // TODO - test
)

// Main system's structure. Contains last block's hash and last block's id, reference to database, mempool of pending transactions, their signatures and some other data.
type Blockchain struct {
	lastHash               string
	lastId                 int
	db                     *database.Database
	pendingTx              []transactions.Transaction // Mempool of pending transactions.
	signatures             map[string]string          // Signatures to verify the transactions in network.
	authorizedTokenSources []string
}

// Since the data stored on a database, we will traverse through it
// using Iterators.
type BlockchainIterator struct {
	currentHash string
	db          *database.Database
}
type BlockchainForwardIterator struct {
	currentId int
	db        *database.Database
}

// Singleton design pattern to have only one instance of the blockchain
// across the program
func GetInstance() *Blockchain {
	once.Do(func() {
		instance = initBlockchain()
		go instance.DataListener()
	})
	return instance
}

// Function that creates the blockchain if it does not exists yet
func initBlockchain() *Blockchain {
	var bc Blockchain
	nc = p2p.GetInstance()
	db := database.GetInstance()
	if !database.IsBlockchainExists() { // Checks whether the blockchain stored on database. If not, mine genesis block and create the blockchain on it.
		genesis := MineGenesisBlock()
		bc = Blockchain{lastHash: genesis.GetHash(), lastId: genesis.GetId(), db: db, signatures: make(map[string]string)}
		bc.db.StoreNewBlockchain(genesis.GetHash(), genesis.GetId(), genesis.Serialize())
	} else {
		lastBlock := DeserializeBlock(db.GetLastBlock())                                                                                                              // If blockchain does exists on database, create one in memory based on the read data.
		bc = Blockchain{lastHash: lastBlock.GetHash(), lastId: lastBlock.GetId(), db: db, signatures: make(map[string]string)}                                        // TODO - rework signatures mempool to sync over multiple peers.
		bc.authorizedTokenSources = append(bc.authorizedTokenSources, "33b02183dba1d072dc7f337013b6bb191fb168b86971feb48f5b5ca3a7da1952c75558bea8b7d1bdf5396fcc7099") // Add the authorized token generator.
	}
	return &bc
}

// Function to aim mining next block. It packs {maxSizeOfTx} transactions from mempool and tries to find block hash based on required difficulty.
func (bc *Blockchain) MineNextBlock() {
	lastBlock := DeserializeBlock(bc.db.GetLastBlock())
	var transactios []transactions.Transaction
	pendingBlock = false

	amountOfTx := 0
	for i := 0; i < maxSizeOfTx && i < len(bc.pendingTx); i++ { // TODO - improve pending transactions model/data structure (Maybe Mutex is needed??)
		if bc.pendingTx[i].IsNil() {
			amountOfTx++
			continue
		}
		//if !security.VerifySignature(bc.signatures[bc.pendingTx[i].GetHash()], bc.pendingTx[i].GetHash(), bc.pendingTx[i].GetSender()) {
		//	bc.pendingTx[i] = transactions.GetNil()
		//}	// TODO - check transaction's signature
		transactios = append(transactios, bc.pendingTx[i])
		amountOfTx++
	}
	bc.pendingTx = bc.pendingTx[amountOfTx:] // TODO -  improve for dynamic use
	newBlock := MineBlock(lastBlock.GetId()+1, difficulty, string(lastBlock.GetHash()), transactios)

	nc.BroadcastBlock(newBlock.Serialize())

	bc.db.StoreBlock(newBlock.GetHash(), newBlock.GetId(), newBlock.Serialize())
	bc.lastId = newBlock.GetId()
	bc.lastHash = newBlock.GetHash()
	//for i := range transactios {
	//	txHash := transactios[i].GetHash()
	//	bc.db.StoreSignature(txHash, bc.signatures[txHash])
	//}
}

func (bc *Blockchain) MineControlBlock(transaction transactions.Transaction) int {
	lastBlock := DeserializeBlock(bc.db.GetLastBlock())
	var txs []transactions.Transaction

	txs = append(txs, transaction)

	newBlock := MineBlock(lastBlock.GetId()+1, difficulty, string(lastBlock.GetHash()), txs)

	nc.BroadcastBlock(newBlock.Serialize()) // TODO - improve

	bc.db.StoreBlock(newBlock.GetHash(), newBlock.GetId(), newBlock.Serialize())
	bc.lastId = newBlock.GetId()
	bc.lastHash = newBlock.GetHash()

	return newBlock.GetId()
}

func (bc *Blockchain) BlockFound(block *Block) {

}

/*==============================================Getters for various blockchain members================================*/
func (bc *Blockchain) GetBlocksAmount() int {
	return bc.lastId + 1
}
func (bc *Blockchain) GetPendingTransactions() []transactions.Transaction {
	return bc.pendingTx
}
func (bc Blockchain) GetSignature(txHash string) string {
	signature := bc.db.GetSignatureByHash(txHash)
	return signature
}

func (bc Blockchain) GetAuthorizedTokenGenerators() []string {
	return bc.authorizedTokenSources
}
func (bc *Blockchain) GetAllSignatures() map[string]string {
	sigs := make(map[string]string, bc.lastId*4)
	it := bc.ForwardIterator()
	for {
		block := it.Next()
		for i := range block.GetTransactions() {
			tx := block.GetTransactions()[i]
			fmt.Println(tx)
			sigs[block.GetTransactions()[i].Hash] = bc.GetSignature(block.GetTransactions()[i].Hash)
		}
		if block.GetId() == bc.lastId {
			break
		}
	}
	return sigs
}
func (bc Blockchain) GetBlockById(id int) *Block {
	block := DeserializeBlock(bc.db.GetBlockById(id))
	return block
}
func (bc *Blockchain) GetLastBlock() *Block {
	it := bc.Iterator()
	return it.Next()
}
func (bc *Blockchain) GetBlockByHash(blockHash string) *Block {
	block := DeserializeBlock(bc.db.GetBlockByHash(blockHash))
	return block
}
func (bc *Blockchain) GetPendingTransactionByHash(txHash string) transactions.Transaction {
	for _, tx := range bc.GetPendingTransactions() {
		if tx.GetHash() == txHash {
			return tx
		}
	}
	return transactions.GetNil()
}
func (bc *Blockchain) GetBlockHashes() [][]byte { // TODO - simplify function by simply traversing through block keys in database
	it := bc.Iterator()
	var hashes = make([][]byte, bc.GetLastBlock().GetId()+1)
	for {
		block := it.Next()
		hashes[block.GetId()] = append(hashes[block.GetId()], []byte(block.GetHash())...)
		if block.GetPreviousHash() == "0" {
			break
		}
	}
	return hashes
}

func (bc *Blockchain) GetControlSig(tag string) transactions.Transaction {
	it := bc.Iterator()

	for {
		block := it.Next()
		txs := block.GetTransactions()
		for _, tx := range txs {
			if tx.GetTag() == tag && tx.GetAmount() == 0 {
				return tx
			}
		}
		if block.GetPreviousHash() == "0" {
			break
		}
	}
	return transactions.GetNil()
}

// Function to check whether Generator transaction is valid by checking control block for proper receipents
func (bc *Blockchain) CheckControlBlock(blockId string, pubKey string) bool {
	if blockId == "" {
		return false // If signature didn't exist, therefore no valid control block was found, and the transaction won't be valid.
	}
	bId, _ := strconv.Atoi(blockId)
	block := bc.GetBlockById(bId)

	controlArr := strings.Split(block.GetTransactions()[0].GetReceiver(), ",")

	for _, value := range controlArr {
		if value == pubKey {
			return true
		}
	}
	return false

}

func (bc *Blockchain) GetTxAmount() int { // TODO - simplify function by simply traversing through block keys in database
	it := bc.Iterator()
	var amount = 0
	for {
		block := it.Next()
		amount += len(block.GetTransactions())
		if block.GetPreviousHash() == "0" {
			break
		}
	}
	return amount
}

/*====================================================================================================================*/

//func (bc Blockchain) InsertToChain(block Block) {	// TODO - Maybe such function can be used to replace single faulty block in the database
//	var remainingTx []transactions.Transaction
//	for i := 0; i < len(bc.pendingTx); i += maxSizeOfTx {
//		remainingTx = bc.pendingTx[i : i+maxSizeOfTx]
//		bc.chain = append(bc.chain, MineBlock(0, bc.difficulty, bc.GetLastBlock().GetHash(), remainingTx))
//		bc.numOfBlocks++
//	}
//}

func (bc *Blockchain) AppendSignature(txHash string, signature string) {
	for idx := range bc.GetPendingTransactions() {
		if bc.GetPendingTransactions()[idx].GetHash() == txHash {
			bc.pendingTx[idx].AddSignature(signature)
		}
	}
}

func (bc *Blockchain) AddTransaction(transaction transactions.Transaction) []string {
	var response = make([]string, 2)
	response[0] = transaction.GetHash()

	if transaction.GetSender() != "Generator" && transaction.GetAmount() > 0 {
		hash := utility.PostRequest(transaction.GetSender(), transaction.GetSignature()) // Testing Signature and therefore identity of the sender...
		if hash != transaction.GetHash() {                                               // Verify properly signed transaction.
			log.Println("Log Err: Signature from " + transaction.GetSender() + " could not be verified.... rejecting (trying to hack???)")

			response[1] = "Bad Signature"
			return response
		}
		if bc.GetBalanceForAddress(transaction.GetSender(), transaction.GetTag()) < 1 {
			log.Println("Log Err: Balance of " + transaction.GetSender() + " equals to 0 in poll tag " + transaction.GetTag())

			response[1] = "Balance equals 0 in poll " + transaction.GetTag() + " from " + transaction.GetSender() + " to " + transaction.GetReceiver()
			return response

		}
	} else if !bc.CheckControlBlock(transaction.GetSignature(), transaction.GetReceiver()) {
		log.Println("Log Err: no proper control block was found for transaction from " + transaction.GetSender() + ".... rejecting.....")

		response[1] = "Signature error w/control block"
		return response

	}

	bc.pendingTx = append(bc.pendingTx, transaction)

	response[1] = "Success : TxHash: " + transaction.GetHash()
	return response
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.lastHash, bc.db}

	return bci
}
func (bc *Blockchain) ForwardIterator() *BlockchainForwardIterator {
	bcfi := &BlockchainForwardIterator{0, bc.db}

	return bcfi
}

func (i *BlockchainIterator) Next() *Block {

	block := DeserializeBlock(i.db.GetBlockByHash(i.currentHash))
	i.currentHash = block.GetPreviousHash()

	return block
}
func (i *BlockchainForwardIterator) Next() *Block {

	block := DeserializeBlock(i.db.GetBlockById(i.currentId))
	i.currentId++

	return block
}

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

func (bc *Blockchain) TraverseBlockchain() []*Block {
	var blocks []*Block
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
func (bc *Blockchain) TraverseForwardBlockchain() []*Block {

	var blocks []*Block
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

// Checks vote balance for given user in certain poll
func (bc *Blockchain) GetBalanceForAddress(address string, pollTag string) int {
	amount := 0
	it := bc.Iterator()
	for {
		block := it.Next()
		txs := block.GetTransactions()
		for _, tx := range txs {
			if tx.GetTag() == pollTag {
				if tx.GetReceiver() == address {

					amount += tx.GetAmount()
				}
				if tx.GetSender() == address {
					amount -= tx.GetAmount()
				}
			}
		}

		if block.GetPreviousHash() == "0" {
			break
		}
	}
	return amount

}

func (bc *Blockchain) AddBlock(block *Block) bool { // TODO - Rework that function (should work now)

	if block.GetId() == 0 {
		bc.db.StoreBlock(block.GetHash(), block.GetId(), block.Serialize())
		bc.lastId = 0
		bc.lastHash = block.GetHash()

		return true
	}

	if !(block.GetId() == bc.GetLastBlock().GetId()+1) {
		return false
	}
	if !(block.GetPreviousHash() == bc.GetLastBlock().GetHash()) {
		return false
	}
	isValid := block.ValidateBlock()

	if isValid {
		bc.db.StoreBlock(block.GetHash(), block.GetId(), block.Serialize())
		bc.lastId = block.GetId()
		bc.lastHash = block.GetHash()
	}
	return isValid
}

// Funtion that runs in it's own goroutine and checks whether the blockchain still valid.
// If by any reason the fail will check, the whole blockchain will be deleted and a copy
// of it will be requested from nearby node //	TODO - replace only the faulty blocks to reduce network bandwidth
func (bc *Blockchain) DataListener() {
	for {
		time.Sleep(15 * time.Second)
		if !database.IsBlockchainExists() {
			fmt.Println("Database deleted alert")
			bc = initBlockchain() // TODO - Add auto-restart on db deletion
		}
		if !bc.ValidateChain() {
			fmt.Println("Blockchain data was compromised... requesting new copy from neighbor peer...")
			bc.lastId = 0
			for _, node := range nc.GetKnownNodes() {
				if nc.GetNodeAddress() != node {
					nc.SendVersion(node, bc.lastId)
				}
			}
		}
	}
}
