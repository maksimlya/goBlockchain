package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"goBlockchain/Blockchain"
	"goBlockchain/Security"
	"goBlockchain/Transactions"
)

//CalculateHash hashes the values of a TestContent

func main() {

	hasher := sha256.New()
	hasher.Write([]byte("MainNodeId1"))
	nodeHash := hex.EncodeToString(hasher.Sum(nil))

	nodeKey := Security.GenerateKey(nodeHash)

	tx1 := Transactions.Tx(nodeKey, "Tomer", 1, "mm")
	tx2 := Transactions.Tx(nodeKey, "Yaki", 1, "mm")
	tx3 := Transactions.Tx(nodeKey, "Koko", 1, "mm")
	tx4 := Transactions.Tx(nodeKey, "Momo", 1, "mm")

	sign1 := Security.Sign(tx1.GetId(), nodeHash)
	sign2 := Security.Sign(tx2.GetId(), nodeHash)
	sign3 := Security.Sign(tx3.GetId(), nodeHash)
	sign4 := Security.Sign(tx4.GetId(), nodeHash)

	valid1 := Security.VerifySignature(sign1, tx1.GetId(), nodeKey)
	valid2 := Security.VerifySignature(sign2, tx2.GetId(), nodeKey)
	valid3 := Security.VerifySignature(sign3, tx3.GetId(), nodeKey)
	valid4 := Security.VerifySignature(sign4, tx4.GetId(), nodeKey)

	fmt.Println(valid1)
	fmt.Println(valid2)
	fmt.Println(valid3)
	fmt.Println(valid4)

	//tx1 := Transactions.Tx("Yaki", "Tomer", 10, "Wow")
	//tx2 := Transactions.Tx("Yaki", "Mas Hahnasa", 10000, "Arnona")
	//tx3 := Transactions.Tx("Yaki", "Tomer", 5, "Arnona")
	//tx4 := Transactions.Tx("Yaki", "Adi", 10, "Takataka")
	//tx5 := Transactions.Tx("Tomer", "Momo", 10, "Wow")
	//tx6 := Transactions.Tx("Yaki", "Popo", 10000, "Arnona")
	//tx7 := Transactions.Tx("Yaki", "Zozo", 51, "Arnona")
	//tx8 := Transactions.Tx("Yaki", "Koko", 10, "Takataka")
	//
	//var list []Transactions.Transaction
	//
	//list = append(list, tx1)
	//list = append(list, tx2)
	//list = append(list, tx3)
	//list = append(list, tx4)
	//
	////merkle, _ := DataStructures.NewTree(list)
	//
	////	merkle.Root.PrintHash()
	//
	////g := Security.GenerateKey("795A433949D3340E7CBA7971DE1B428830C15D901B65303B3A65C0A45EE3F498")
	////fmt.Println(g)
	////
	////t := "daklfwklwklkdlcl asdca cascascac"
	////
	////tt := []byte(t)
	////var pp []int
	////
	////for i := 0; i < len(tt); i++ {
	////	pp = append(pp, Security.Encrypt(strconv.Itoa(int(tt[i])), g))
	////}
	////
	////privKey := Security.GeneratePrivKey("795A433949D3340E7CBA7971DE1B428830C15D901B65303B3A65C0A45EE3F498")
	////
	////for i := 0; i < len(tt); i++ {
	////	pp[i] = Security.Encrypt(strconv.Itoa(int(pp[i])), privKey)
	////}
	////
	////str := []string{}
	////
	////for i := range pp {
	////	bytea := byte(pp[i])
	////
	////	str = append(str, string(bytea))
	////}
	////
	////res := strings.Join(str, "")
	////fmt.Println(res)
	//
	blockchain := Blockchain.InitBlockchain()
	//
	//blockchain.AddTransaction(tx1, sign1)
	//blockchain.AddTransaction(tx2, sign2)
	//blockchain.AddTransaction(tx3, sign3)
	//blockchain.AddTransaction(tx4, sign4)
	//
	//blockchain.MineNextBlock()

	//blockchain.MineNextBlock()
	//blockchain.MineNextBlock()
	//blockchain.MineNextBlock()
	//blockchain.MineNextBlock()
	//

	//bl := blockchain.GetLastBlock()

	j := blockchain.GetSignature(tx3.Id)
	fmt.Println(j)
	//k := make(map[string]string)
	//json.Unmarshal(j,&k)
	//fmt.Println(k)
	//fmt.Println(blockchain)

	//for {
	//	fmt.Println("Hello")
	//	time.Sleep(time.Second)
	//}
	//pubKey := Security.GenerateKey("A034B1566E979D3C5FE487BF3CF721FF3517570E1151DFC67D0329A54A48F9F8")
	//fmt.Println(pubKey)
	//signature := Security.Sign("muhahsada", "A034B1566E979D3C5FE487BF3CF721FF3517570E1151DFC67D0329A54A48F9F8")
	//
	//fmt.Println(signature)
	//fmt.Println(Security.VerifySignature(signature, "muhahsada", pubKey))

	//	cli := CommandInterface.CLI{blockchain}
	//	cli.Run()

	fmt.Println(blockchain.GetBalanceForAddress("Tomer"))

}
