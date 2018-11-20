package Transactions

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
)

type Transaction struct {
	Id        string
	From      string
	To        string
	Amount    int
	Tag       string
	Timestamp string
}

func Tx(from string, to string, amount int, tag string) Transaction {
	timestamp := time.Now().Format("02-01-2006 15:04:05")
	shaHasher := sha256.New()
	shaHasher.Write([]byte(from + to + strconv.Itoa(amount) + tag + timestamp))
	tx := Transaction{Id: hex.EncodeToString(shaHasher.Sum(nil)), From: from, To: to, Amount: amount, Tag: tag, Timestamp: timestamp}

	return tx
}

func GetNil() Transaction {
	timestamp := "0"
	tx := Transaction{Id: "0", From: "nil", To: "nil", Amount: 0, Tag: "nil", Timestamp: timestamp}
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
	hasher.Write([]byte(t.Id + t.From + t.To + string(t.Amount) + t.Tag + t.Timestamp))
	f := hex.EncodeToString(hasher.Sum(nil))
	return f
}

func Equals(first Transaction, second Transaction) bool {
	hasher := sha256.New()
	hasher.Write([]byte(string(first.Id) + first.From + first.To + string(first.Amount) + first.Tag + first.Timestamp))
	firstHash := hex.EncodeToString(hasher.Sum(nil))
	hasher = sha256.New()
	hasher.Write([]byte(string(second.Id) + second.From + second.To + string(second.Amount) + second.Tag + second.Timestamp))
	secondHash := hex.EncodeToString(hasher.Sum(nil))

	return firstHash == secondHash
}

func (t Transaction) GetTag() string {
	return t.Tag
}

func (t Transaction) GetId() string {
	return t.Id
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
	s += "Tx Id: " + tx.Id + "\n"
	s += "From Address: " + tx.From + "\n"
	s += "To Address: " + tx.To + "\n"
	s += "Amount: " + strconv.Itoa(tx.Amount) + "\n"
	s += "Timestamp: " + tx.Timestamp + "\n"
	s += "Tag: " + tx.Tag + "\n"
	s += "}\n"

	return s
}
