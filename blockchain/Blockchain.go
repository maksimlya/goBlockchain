package blockchain

import (
	"goBlockchain/database"
	"goBlockchain/security"
	"goBlockchain/transactions"
)

type Blockchain struct {
	chain       []Block
	LastHash    string
	LastId      int
	Db          database.Database
	numOfBlocks int
	Difficulty  int
	pendingTx   []transactions.Transaction
	Signatures  map[string]string
}

type BlockchainIterator struct {
	currentHash string
	db          database.Database
}
type BlockchainForwardIterator struct {
	currentId int
	db        database.Database
}

var maxSizeOfTx = 4

func InitBlockchain() *Blockchain {
	bc := database.GetBlockchain()
	return bc

}

func (bc Blockchain) getSignature(key string) string {
	return bc.Signatures[key]
}

func (bc Blockchain) GetSignature(txHash string) string {
	signature := bc.Db.GetSignatureByHash(txHash)
	return signature
}

func (bc *Blockchain) GetAllSignatures() map[string]string {
	sigs := make(map[string]string, bc.LastId*4)
	it := bc.ForwardIterator()
	for {
		block := it.Next()
		for i := range block.GetTransactions() {
			sigs[block.GetTransactions()[i].Hash] = bc.GetSignature(block.GetTransactions()[i].Hash)
		}
		if block.GetId() == bc.LastId {
			break
		}
	}

	return sigs
}

//func (bc Blockchain) InsertToChain(block Block) {
//	var remainingTx []transactions.Transaction
//	for i := 0; i < len(bc.pendingTx); i += maxSizeOfTx {
//		remainingTx = bc.pendingTx[i : i+maxSizeOfTx]
//		bc.chain = append(bc.chain, MineBlock(0, bc.difficulty, bc.GetLastBlock().GetHash(), remainingTx))
//		bc.numOfBlocks++
//	}
//}

func (bc *Blockchain) MineNextBlock() {

	lastBlock := bc.Db.GetLastBlock()

	var transactios []transactions.Transaction
	amountOfTx := 0
	for i := 0; i < maxSizeOfTx && i < len(bc.pendingTx); i++ { // TODO - improve pending transactions model/data structure
		if bc.pendingTx[i].IsNil() {
			amountOfTx++
			continue
		}
		if !security.VerifySignature(bc.Signatures[bc.pendingTx[i].GetHash()], bc.pendingTx[i].GetHash(), bc.pendingTx[i].GetSender()) {
			bc.pendingTx[i] = transactions.GetNil()
		}
		transactios = append(transactios, bc.pendingTx[i])
		amountOfTx++
	}
	bc.pendingTx = bc.pendingTx[amountOfTx:] // TODO -  improve for dynamic use
	newBlock := MineBlock(lastBlock.GetId()+1, bc.Difficulty, string(lastBlock.GetHash()), transactios)
	bc.Db.StoreBlock(newBlock)
	for i := range transactios {
		txHash := transactios[i].GetHash()
		bc.Db.StoreSignature(txHash, bc.Signatures[txHash])
	}
}

func (bc Blockchain) SearchBlock(hash string) Block {
	var b Block
	for i := 0; i < bc.numOfBlocks; i++ {
		if hash == bc.chain[i].GetHash() {
			return bc.chain[i]
		}
	}
	return b
}

func (bc Blockchain) GetBlockById(id int) *Block {
	block := bc.Db.GetBlockById(id)
	return block
}

func (bc *Blockchain) AddTransaction(transaction transactions.Transaction, signature string) {
	if !security.VerifySignature(signature, transaction.GetHash(), transaction.GetSender()) {
		return
	}
	bc.Signatures[transaction.GetHash()] = signature
	bc.pendingTx = append(bc.pendingTx, transaction)
}

func (bc *Blockchain) GetLastBlock() *Block {
	it := bc.Iterator()
	return it.Next()
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.LastHash, bc.Db}

	return bci
}
func (bc *Blockchain) ForwardIterator() *BlockchainForwardIterator {
	bcfi := &BlockchainForwardIterator{0, bc.Db}

	return bcfi
}

func (i *BlockchainIterator) Next() *Block {

	block := i.db.GetBlockByHash(i.currentHash)
	i.currentHash = block.GetPreviousHash()

	return block
}
func (i *BlockchainForwardIterator) Next() *Block {

	block := i.db.GetBlockById(i.currentId)
	i.currentId++

	return block
}

//func (bc *blockchain) String() string {
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
//		s += "transactions: {\n"
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

		if block.GetId() == bc.LastId {
			break
		}
	}
	return blocks
}

//// TODO - rework function for polls use ( aka check balance for a given poll tag )
//func (bc *Blockchain) GetBalanceForAddress(address string) int {
//	bc.db, _ = bolt.Open("blockchain.db", 0600, nil)
//	var amount = 0
//	it := bc.Iterator()
//	err := bc.db.View(func(tx *bolt.Tx) error {
//		s := tx.Bucket([]byte("signatures"))
//
//		for {
//			block := it.Next()
//			for _, element := range block.GetTransactions() {
//				signature := string(s.Get([]byte(element.GetHash())))
//				if element.GetReceiver() == address {
//					if security.VerifySignature(signature, element.GetHash(), element.GetSender()) {
//						amount += element.GetAmount()
//					}
//				}
//				if element.GetSender() == address {
//					if security.VerifySignature(signature, element.GetHash(), element.GetSender()) {
//						amount -= element.GetAmount()
//					}
//				}
//			}
//
//			if block.GetPreviousHash() == "0" {
//				break
//			}
//		}
//
//		return nil
//	})
//	if err != nil {
//		fmt.Println(err)
//	}
//	bc.db.Close()
//	return amount
//
//}
