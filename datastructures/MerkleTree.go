package datastructures

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"goBlockchain/transactions"
	"math"
)

// Enum representation in golang
type Offset int

// Each Node stores it's offset ( aka if it allocated right or left-side of the parent )
const (
	Roof  Offset = 0
	Right Offset = 1
	Left  Offset = 2
)

//MerkleTree is the container for the tree. It holds a pointer to the Root of the tree,
//a list of pointers to the leaf nodes, and the merkle Root.
type MerkleTree struct {
	Root    *Node
	MerRoot string
	Leafs   []*Node
}

// This struct will be used to send to SPV clients ( The blockchain nodes that contain only block headers )
// In order that they will be able to re-build the merkle root and verify that transaction belongs to the block.
type ProofElement struct {
	Hash  string
	Place Offset
}

//Node represents a node, Root, or leaf in the tree. It stores pointers to its immediate
//relationships, a hash, the content stored if it is a leaf, and other metadata.
// Place will store it's position related to the parent (Left or Right) OR Roof value only if the Node is root of the whole tree.
type Node struct {
	Parent *Node
	Left   *Node
	Right  *Node
	Place  Offset
	Leaf   bool
	Dup    bool
	Hash   string
	Tx     transactions.Transaction
}

//verifyNode walks down the tree until hitting a leaf, calculating the hash at each level
//and returning the resulting hash of Node n.
func (n *Node) PrintHash() {
	if n.Right != nil {
		n.Right.PrintHash()
	}
	if n.Left != nil {
		n.Left.PrintHash()
	}
	fmt.Println(n.Hash)
}

func (t *MerkleTree) PrintLevels() map[int][]string {
	if t == nil {
		return nil
	}
	var levels = int(math.Log2(float64(len(t.Leafs))))
	var stages = make(map[int][]string, levels)

	stages = t.Root.visitNode(0, stages)

	return stages
}

func (n *Node) visitNode(level int, stages map[int][]string) map[int][]string {
	stages[level] = append(stages[level], n.Hash)

	if !n.Leaf {
		stages = n.Left.visitNode(level+1, stages)
		stages = n.Right.visitNode(level+1, stages)
	}
	return stages
}

func (n *Node) verifyNode() string {
	if n.Leaf {
		return transactions.CalcHash(n.Tx)
	}
	rightBytes := n.Right.verifyNode()
	leftBytes := n.Left.verifyNode()

	h := sha256.New()
	h.Write([]byte(leftBytes + rightBytes))

	return hex.EncodeToString(h.Sum(nil))
}

//calculateNodeHash is a helper function that calculates the hash of the node.

func (n *MerkleTree) GetTransactions() []transactions.Transaction { // TODO - FIX BUF
	var transactions []transactions.Transaction
	for _, j := range n.Leafs {
		if !j.Dup {
			transactions = append(transactions, j.Tx)
		}
	}
	return transactions
}

func (n *Node) calculateNodeHash() ([]byte, error) {
	if n.Leaf {
		return []byte(transactions.CalcHash(n.Tx)), nil
	}

	h := sha256.New()
	if _, err := h.Write([]byte(n.Left.Hash + n.Right.Hash)); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

//NewTree creates a new Merkle Tree using the content cs.
func NewTree(cs []transactions.Transaction) (*MerkleTree, error) {
	Root, Leafs, err := buildWithContent(cs)
	if err != nil {
		return nil, err
	}
	t := &MerkleTree{
		Root:    Root,
		MerRoot: Root.Hash,
		Leafs:   Leafs,
	}
	return t, nil
}

//buildWithContent is a helper function that for a given set of Contents, generates a
//corresponding tree and returns the Root node, a list of leaf nodes, and a possible error.
//Returns an error if cs contains no Contents.
func buildWithContent(cs []transactions.Transaction) (*Node, []*Node, error) {
	if len(cs) == 0 {
		return nil, nil, errors.New("error: cannot construct tree with no content")
	}
	var Leafs []*Node
	for _, tx := range cs {
		hash := tx.GetHash()

		Leafs = append(Leafs, &Node{
			Hash: hash,
			Tx:   tx,
			Leaf: true,
		})
	}
	if len(Leafs)%2 == 1 {
		duplicate := &Node{
			Hash: Leafs[len(Leafs)-1].Hash,
			Tx:   Leafs[len(Leafs)-1].Tx,
			Leaf: true,
			Dup:  true,
		}
		Leafs = append(Leafs, duplicate)
	}
	Root, err := buildIntermediate(Leafs)
	if err != nil {
		return nil, nil, err
	}

	return Root, Leafs, nil
}

//buildIntermediate is a helper function that for a given list of leaf nodes, constructs
//the intermediate and Root levels of the tree. Returns the resulting Root node of the tree.
func buildIntermediate(nl []*Node) (*Node, error) {
	var nodes []*Node
	for i := 0; i < len(nl); i += 2 {
		h := sha256.New()
		var left, right int = i, i + 1
		if i+1 == len(nl) {
			right = i
		}
		nl[left].Place = Left
		nl[right].Place = Right
		chash := (nl[left].Hash + nl[right].Hash)

		if _, err := h.Write([]byte(chash)); err != nil {
			return nil, err
		}
		n := &Node{
			Left:  nl[left],
			Right: nl[right],
			Hash:  hex.EncodeToString(h.Sum(nil)),
		}
		nodes = append(nodes, n)
		nl[left].Parent = n
		nl[right].Parent = n
		if len(nl) == 2 {
			return n, nil
		}
	}
	return buildIntermediate(nodes)
}

//MerkleRoot returns the unverified Merkle Root (hash of the Root node) of the tree.
func (m *MerkleTree) MerkleRoot() string {
	return m.MerRoot
}

//VerifyTree verify tree validates the hashes at each level of the tree and returns true if the
//resulting hash at the Root of the tree matches the resulting Root hash; returns false otherwise.
func (m *MerkleTree) VerifyTree() (bool, error) {
	calculatedMerkleRoot := m.Root.verifyNode()

	if m.MerRoot == calculatedMerkleRoot {
		return true, nil
	}
	return false, nil
}

// Function that gets Transaction and returns proof element which later can be used to re-build the merkle root
// Used to verify that given transaction belongs to the merkle tree, aka is mined into the block.
// Returns nil if the transaction is not in the tree.
func (m *MerkleTree) GetProofElements(tx transactions.Transaction) []ProofElement {
	var proofs []ProofElement
	var tmpNode *Node
	for _, leaf := range m.Leafs {
		if tx.GetHash() == leaf.Hash {
			tmpNode = leaf
			for tmpNode.Place != Roof {
				switch tmpNode.Place {
				case Right:
					proofs = append(proofs, ProofElement{Hash: tmpNode.Parent.Left.Hash, Place: Left})
					tmpNode = tmpNode.Parent
					break
				case Left:
					proofs = append(proofs, ProofElement{Hash: tmpNode.Parent.Right.Hash, Place: Right})
					tmpNode = tmpNode.Parent
					break
				}
			}
			break
		}
	}
	return proofs
}

// Function that receives transaction, proof elements and merkle root, and tries to re-build
// The merkle root based on the given proof elements.
// Returns true if the result equals to the assumed merkle root.
func VerifyContent(tx transactions.Transaction, proofs []ProofElement, merkleRoot string) bool {
	result := tx.GetHash()

	for i := range proofs {
		hasher := sha256.New()
		if proofs[i].Place == Right {
			hasher.Write([]byte(result + proofs[i].Hash))
			result = hex.EncodeToString(hasher.Sum(nil))
		} else {
			hasher.Write([]byte(proofs[i].Hash + result))
			result = hex.EncodeToString(hasher.Sum(nil))
		}
	}
	return merkleRoot == result
}

func (m *MerkleTree) GetTransactionsWithTag(tag string) []transactions.Transaction {
	var result []transactions.Transaction
	for _, l := range m.Leafs {
		if l.Tx.GetTag() == tag && !l.Dup {
			result = append(result, l.Tx)
		}
	}

	return result
}

func (b *MerkleTree) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b.GetTransactions())
	if err != nil {
		fmt.Println(err.Error())
	}

	return result.Bytes()
}

func Deserialize(d []byte) *MerkleTree {
	var tree *MerkleTree
	var tx []transactions.Transaction

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&tx)

	tree, _ = NewTree(tx)
	if err != nil {
		fmt.Println(err)
	}

	return tree
}
