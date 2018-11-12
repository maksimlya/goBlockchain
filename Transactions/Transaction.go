package Transactions

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
)

type Transaction struct {
	Id        int
	From      string
	To        string
	Amount    int
	Tag       string
	Timestamp string
}

var idx = 0

func Tx(from string, to string, amount int, tag string) Transaction {
	tx := Transaction{Id: idx, From: from, To: to, Amount: amount, Tag: tag, Timestamp: time.Now().Format("02-01-2006 15:04:05")}
	idx++
	return tx
}

func CalcHash(t Transaction) string {
	hasher := sha256.New()
	hasher.Write([]byte(string(t.Id) + t.From + t.To + string(t.Amount) + t.Tag + t.Timestamp))
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

func (tx *Transaction) String() string {
	s := ""
	s += "{\n" //fmt.Sprint(l)
	s += "Tx Id: " + strconv.Itoa(tx.Id) + "\n"
	s += "From Address: " + tx.From + "\n"
	s += "To Address: " + tx.To + "\n"
	s += "Amount: " + strconv.Itoa(tx.Amount) + "\n"
	s += "Timestamp: " + tx.Timestamp + "\n"
	s += "Tag: " + tx.Tag + "\n"
	s += "}\n"

	return s
}
