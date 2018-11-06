package DataStructures

import (
	"crypto/sha256"
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

//MerkleTree is the container for the tree. It holds a pointer to the root of the tree,
//a list of pointers to the leaf nodes, and the merkle root.
type MerkleTree struct {
	Root       *Node
	merkleRoot string
	Leafs      []*Node
}

//Node represents a node, root, or leaf in the tree. It stores pointers to its immediate
//relationships, a hash, the content stored if it is a leaf, and other metadata.
type Node struct {
	Parent *Node
	Left   *Node
	Right  *Node
	leaf   bool
	dup    bool
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
	if n.leaf {
		return Transactions.CalcHash(n.Tx)
	}
	rightBytes := n.Right.verifyNode()
	leftBytes := n.Left.verifyNode()

	h := sha256.New()
	h.Write([]byte(leftBytes + rightBytes))

	return hex.EncodeToString(h.Sum(nil))
}

//calculateNodeHash is a helper function that calculates the hash of the node.
func (n *Node) calculateNodeHash() ([]byte, error) {
	if n.leaf {
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
	root, leafs, err := buildWithContent(cs)
	if err != nil {
		return nil, err
	}
	t := &MerkleTree{
		Root:       root,
		merkleRoot: root.Hash,
		Leafs:      leafs,
	}
	return t, nil
}

//buildWithContent is a helper function that for a given set of Contents, generates a
//corresponding tree and returns the root node, a list of leaf nodes, and a possible error.
//Returns an error if cs contains no Contents.
func buildWithContent(cs []Transactions.Transaction) (*Node, []*Node, error) {
	if len(cs) == 0 {
		return nil, nil, errors.New("error: cannot construct tree with no content")
	}
	var leafs []*Node
	for _, tx := range cs {
		hash := Transactions.CalcHash(tx)

		leafs = append(leafs, &Node{
			Hash: hash,
			Tx:   tx,
			leaf: true,
		})
	}
	if len(leafs)%2 == 1 {
		duplicate := &Node{
			Hash: leafs[len(leafs)-1].Hash,
			Tx:   leafs[len(leafs)-1].Tx,
			leaf: true,
			dup:  true,
		}
		leafs = append(leafs, duplicate)
	}
	root, err := buildIntermediate(leafs)
	if err != nil {
		return nil, nil, err
	}

	return root, leafs, nil
}

//buildIntermediate is a helper function that for a given list of leaf nodes, constructs
//the intermediate and root levels of the tree. Returns the resulting root node of the tree.
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

//MerkleRoot returns the unverified Merkle Root (hash of the root node) of the tree.
func (m *MerkleTree) MerkleRoot() string {
	return m.merkleRoot
}

//RebuildTree is a helper function that will rebuild the tree reusing only the content that
//it holds in the leaves.
func (m *MerkleTree) RebuildTree() error {
	var cs []Transactions.Transaction
	for _, c := range m.Leafs {
		cs = append(cs, c.Tx)
	}
	root, leafs, err := buildWithContent(cs)
	if err != nil {
		return err
	}
	m.Root = root
	m.Leafs = leafs
	m.merkleRoot = root.Hash
	return nil
}

//RebuildTreeWith replaces the content of the tree and does a complete rebuild; while the root of
//the tree will be replaced the MerkleTree completely survives this operation. Returns an error if the
//list of content cs contains no entries.
func (m *MerkleTree) RebuildTreeWith(cs []Transactions.Transaction) error {
	root, leafs, err := buildWithContent(cs)
	if err != nil {
		return err
	}
	m.Root = root
	m.Leafs = leafs
	m.merkleRoot = root.Hash
	return nil
}

//VerifyTree verify tree validates the hashes at each level of the tree and returns true if the
//resulting hash at the root of the tree matches the resulting root hash; returns false otherwise.
func (m *MerkleTree) VerifyTree() (bool, error) {
	calculatedMerkleRoot := m.Root.verifyNode()

	if m.merkleRoot == calculatedMerkleRoot {
		return true, nil
	}
	return false, nil
}

//VerifyContent indicates whether a given content is in the tree and the hashes are valid for that content.
//Returns true if the expected Merkle Root is equivalent to the Merkle root calculated on the critical path
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
				if currentParent.Left.leaf && currentParent.Right.leaf {
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
