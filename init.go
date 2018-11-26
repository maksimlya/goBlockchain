package main

import (
	"fmt"
	"goBlockchain/p2p/handlers"
	"goBlockchain/security"
	"goBlockchain/transactions"
	"goBlockchain/utility"
)

//CalculateHash hashes the values of a TestContent

func main() {

	nodeHash := utility.Hash("MainNodeId1")

	nodeKey := security.GenerateKey(nodeHash)

	tx1 := transactions.Tx(nodeKey, "Tomer", 1, "mm")
	tx2 := transactions.Tx(nodeKey, "Yaki", 1, "mm")
	tx3 := transactions.Tx(nodeKey, "Koko", 1, "mm")
	tx4 := transactions.Tx(nodeKey, "Momo", 1, "mm")

	sign1 := security.Sign(tx1.GetHash(), nodeHash)
	sign2 := security.Sign(tx2.GetHash(), nodeHash)
	sign3 := security.Sign(tx3.GetHash(), nodeHash)
	sign4 := security.Sign(tx4.GetHash(), nodeHash)

	valid1 := security.VerifySignature(sign1, tx1.GetHash(), nodeKey)
	valid2 := security.VerifySignature(sign2, tx2.GetHash(), nodeKey)
	valid3 := security.VerifySignature(sign3, tx3.GetHash(), nodeKey)
	valid4 := security.VerifySignature(sign4, tx4.GetHash(), nodeKey)

	fmt.Println(valid1)
	fmt.Println(valid2)
	fmt.Println(valid3)
	fmt.Println(valid4)

	//bc := blockchain.GetInstance()
	//
	//bc.AddTransaction(tx1, sign1)
	//bc.AddTransaction(tx2, sign2)
	//bc.AddTransaction(tx3, sign3)
	//bc.AddTransaction(tx4, sign4)
	//////////
	//bc.MineNextBlock()

	//
	//fmt.Println(bc.GetLastBlock())
	//

	go handlers.StartServer("3000")
	//
	go log.Fatal(webserver.Run())

}
