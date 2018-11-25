package transactions

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
	"strconv"
	"time"
)

type Transaction struct {
	Hash      string
	From      string
	To        string
	Amount    int
	Tag       string
	Timestamp string
}

func Tx(from string, to string, amount int, tag string) Transaction {
	timestamp := time.Now().Format("02-01-2006 15:04:05")
	shaHasher := sha256.New()
	shaHasher.Write([]byte(from + to + strconv.Itoa(amount) + tag)) // TODO Validate
	tx := Transaction{Hash: hex.EncodeToString(shaHasher.Sum(nil)), From: from, To: to, Amount: amount, Tag: tag, Timestamp: timestamp}

	return tx
}

func (t *Transaction) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(t)
	if err != nil {
		fmt.Println(err)
	}
	return result.Bytes()
}
func DeserializeTransaction(d []byte) *Transaction {
	var tx Transaction

	decoder := gob.NewDecoder(bytes.NewReader(d))

	err := decoder.Decode(&tx)

	if err != nil && err != io.EOF {
		fmt.Println(err)
	}

	return &tx
}

func GetNil() Transaction {
	timestamp := "0"
	tx := Transaction{Hash: "0", From: "nil", To: "nil", Amount: 0, Tag: "nil", Timestamp: timestamp}
	return tx
}

func (t Transaction) IsNil() bool {
	if t.Tag == "nil" {
		return true
	}
	return false
}

func CalcHash(t Transaction) string {
	hasher := sha256.New()
	hasher.Write([]byte(t.From + t.To + string(t.Amount) + t.Tag + t.Timestamp))
	f := hex.EncodeToString(hasher.Sum(nil))
	return f
}

func Equals(first Transaction, second Transaction) bool {
	hasher := sha256.New()
	hasher.Write([]byte(string(first.Hash) + first.From + first.To + string(first.Amount) + first.Tag + first.Timestamp))
	firstHash := hex.EncodeToString(hasher.Sum(nil))
	hasher = sha256.New()
	hasher.Write([]byte(string(second.Hash) + second.From + second.To + string(second.Amount) + second.Tag + second.Timestamp))
	secondHash := hex.EncodeToString(hasher.Sum(nil))

	return firstHash == secondHash
}

func (t Transaction) GetTag() string {
	return t.Tag
}

func (t Transaction) GetHash() string {
	return t.Hash
}

func (t Transaction) GetSender() string {
	return t.From
}

func (t Transaction) GetAmount() int {
	return t.Amount
}

func (t Transaction) GetReceiver() string {
	return t.To
}

func (tx *Transaction) String() string {
	s := ""
	s += "{\n" //fmt.Sprint(l)
	s += "Tx Id: " + tx.Hash + "\n"
	s += "From Address: " + tx.From + "\n"
	s += "To Address: " + tx.To + "\n"
	s += "Amount: " + strconv.Itoa(tx.Amount) + "\n"
	s += "Timestamp: " + tx.Timestamp + "\n"
	s += "Tag: " + tx.Tag + "\n"
	s += "}\n"

	return s
}
