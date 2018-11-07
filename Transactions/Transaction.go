package Transactions

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
)

type Transaction struct {
	id        int
	from      string
	to        string
	amount    int
	tag       string
	timestamp string
}

var idx = 0

func Tx(from string, to string, amount int, tag string) Transaction {
	tx := Transaction{id: idx, from: from, to: to, amount: amount, tag: tag, timestamp: time.Now().Format("02-01-2006 15:04:05")}
	idx++
	return tx
}

func CalcHash(t Transaction) string {
	hasher := sha256.New()
	hasher.Write([]byte(string(t.id) + t.from + t.to + string(t.amount) + t.tag + t.timestamp))
	f := hex.EncodeToString(hasher.Sum(nil))
	return f
}

func Equals(first Transaction, second Transaction) bool {
	hasher := sha256.New()
	hasher.Write([]byte(string(first.id) + first.from + first.to + string(first.amount) + first.tag + first.timestamp))
	firstHash := hex.EncodeToString(hasher.Sum(nil))
	hasher = sha256.New()
	hasher.Write([]byte(string(second.id) + second.from + second.to + string(second.amount) + second.tag + second.timestamp))
	secondHash := hex.EncodeToString(hasher.Sum(nil))

	return firstHash == secondHash
}

func (t Transaction) GetTag() string {
	return t.tag
}

func (tx *Transaction) String() string {
	s := ""
	s += "{\n" //fmt.Sprint(l)
	s += "Tx Id: " + strconv.Itoa(tx.id) + "\n"
	s += "From Address: " + tx.from + "\n"
	s += "To Address: " + tx.to + "\n"
	s += "Amount: " + strconv.Itoa(tx.amount) + "\n"
	s += "Timestamp: " + tx.timestamp + "\n"
	s += "Tag: " + tx.tag + "\n"
	s += "}\n"

	return s
}
