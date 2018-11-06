package main

import (
	"encoding/hex"
	"fmt"
	"goBlockchain/DataStructures"
	"goBlockchain/Transactions"
)

//CalculateHash hashes the values of a TestContent

func main() {

	tx1 := Transactions.Tx("Yaki", "Tomer", 10, "Wow")
	tx2 := Transactions.Tx("Yaki", "Mas Hahnasa", 10000, "Arnona")
	tx3 := Transactions.Tx("Yaki", "Zona", 5, "Arnona")
	tx4 := Transactions.Tx("Yaki", "Adi", 10, "Takataka")

	var list []Transactions.Transaction

	list = append(list, tx1)
	list = append(list, tx2)
	list = append(list, tx3)
	list = append(list, tx4)

	merkle, _ := DataStructures.NewTree(list)

	fmt.Println(merkle.GetTransactionsWithTag("Arnona"))

	fmt.Println(merkle.HexString())
	fmt.Println(hex.EncodeToString(merkle.MerkleRoot()))

}
