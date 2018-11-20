package DataStructures

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"goBlockchain/Transactions"
)

//Content represents the data that is stored and verified by the tree. A type that
//implements this interface can be used as an item in the tree.

//type Content interface {
//	CalculateHash() ([]byte, error)
//	Equals(other Content) (bool, error)
//}

//MerkleTree is the container for the tree. It holds a pointer to the Root of the tree,
//a list of pointers to the leaf nodes, and the merkle Root.
type MerkleTree struct {
	Root    *Node
	MerRoot string
	Leafs   []*Node
}

//Node represents a node, Root, or leaf in the tree. It stores pointers to its immediate
//relationships, a hash, the content stored if it is a leaf, and other metadata.
type Node struct {
	Parent *Node
	Left   *Node
	Right  *Node
	Leaf   bool
	Dup    bool
	Hash   string
	Tx     Transactions.Transaction
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

func (n *Node) verifyNode() string {
	if n.Leaf {
		return Transactions.CalcHash(n.Tx)
	}
	rightBytes := n.Right.verifyNode()
	leftBytes := n.Left.verifyNode()

	h := sha256.New()
	h.Write([]byte(leftBytes + rightBytes))

	return hex.EncodeToString(h.Sum(nil))
}

//calculateNodeHash is a helper function that calculates the hash of the node.

func (n *MerkleTree) GetTransactions() []Transactions.Transaction { // TODO - FIX BUF
	var transactions []Transactions.Transaction
	for _, j := range n.Leafs {
		if !j.Dup {
			transactions = append(transactions, j.Tx)
		}
	}
	return transactions
}

func (n *Node) calculateNodeHash() ([]byte, error) {
	if n.Leaf {
		return []byte(Transactions.CalcHash(n.Tx)), nil
	}

	h := sha256.New()
	if _, err := h.Write([]byte(n.Left.Hash + n.Right.Hash)); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

//NewTree creates a new Merkle Tree using the content cs.
func NewTree(cs []Transactions.Transaction) (*MerkleTree, error) {
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
func buildWithContent(cs []Transactions.Transaction) (*Node, []*Node, error) {
	if len(cs) == 0 {
		return nil, nil, errors.New("error: cannot construct tree with no content")
	}
	var Leafs []*Node
	for _, tx := range cs {
		hash := Transactions.CalcHash(tx)

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

//RebuildTree is a helper function that will rebuild the tree reusing only the content that
//it holds in the leaves.
func (m *MerkleTree) RebuildTree() error {
	var cs []Transactions.Transaction
	for _, c := range m.Leafs {
		cs = append(cs, c.Tx)
	}
	Root, Leafs, err := buildWithContent(cs)
	if err != nil {
		return err
	}
	m.Root = Root
	m.Leafs = Leafs
	m.MerRoot = Root.Hash
	return nil
}

//RebuildTreeWith replaces the content of the tree and does a complete rebuild; while the Root of
//the tree will be replaced the MerkleTree completely survives this operation. Returns an error if the
//list of content cs contains no entries.
func (m *MerkleTree) RebuildTreeWith(cs []Transactions.Transaction) error {
	Root, Leafs, err := buildWithContent(cs)
	if err != nil {
		return err
	}
	m.Root = Root
	m.Leafs = Leafs
	m.MerRoot = Root.Hash
	return nil
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

//VerifyContent indicates whether a given content is in the tree and the hashes are valid for that content.
//Returns true if the expected Merkle Root is equivalent to the Merkle Root calculated on the critical path
//for a given content. Returns true if valid and false otherwise.
func (m *MerkleTree) VerifyContent(content Transactions.Transaction) (bool, error) {
	for _, l := range m.Leafs {
		ok := Transactions.Equals(l.Tx, content)
		if ok {
			currentParent := l.Parent
			for currentParent != nil {
				h := sha256.New()
				rightBytes, err := currentParent.Right.calculateNodeHash()
				if err != nil {
					return false, err
				}

				leftBytes, err := currentParent.Left.calculateNodeHash()
				if err != nil {
					return false, err
				}
				if currentParent.Left.Leaf && currentParent.Right.Leaf {
					if _, err := h.Write(append(leftBytes, rightBytes...)); err != nil {
						return false, err
					}
					if hex.EncodeToString(h.Sum(nil)) == currentParent.Hash {
						return false, nil
					}
					currentParent = currentParent.Parent
				} else {
					if _, err := h.Write(append(leftBytes, rightBytes...)); err != nil {
						return false, err
					}
					if hex.EncodeToString(h.Sum(nil)) == currentParent.Hash {
						return false, nil
					}
					currentParent = currentParent.Parent
				}
			}
			return true, nil
		}
	}
	return false, nil
}

func (m *MerkleTree) GetTransactionsWithTag(tag string) []Transactions.Transaction {
	var result []Transactions.Transaction
	for _, l := range m.Leafs {
		if l.Tx.GetTag() == tag {
			result = append(result, l.Tx)
		}
	}

	return result
}

//String returns a string representation of the tree. Only leaf nodes are included
//in the output.
func (m *MerkleTree) String() string {
	s := ""
	for _, l := range m.Leafs {
		s += fmt.Sprint(l)
		s += "\n"
	}
	return s
}

func (m *MerkleTree) HexString() string {
	s := ""
	for _, l := range m.Leafs {
		s += string(l.Hash)
		s += "\n"
	}
	return s
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
	var tx []Transactions.Transaction

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&tx)

	tree, _ = NewTree(tx)
	if err != nil {
		fmt.Println(err)
	}

	return tree
}
