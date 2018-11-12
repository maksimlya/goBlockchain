package main

import (
	"goBlockchain/Blockchain"
	"goBlockchain/CommandInterface"
	"goBlockchain/Transactions"
)

//CalculateHash hashes the values of a TestContent

func main() {

	tx1 := Transactions.Tx("Yaki", "Tomer", 10, "Wow")
	tx2 := Transactions.Tx("Yaki", "Mas Hahnasa", 10000, "Arnona")
	tx3 := Transactions.Tx("Yaki", "Zona", 5, "Arnona")
	tx4 := Transactions.Tx("Yaki", "Adi", 10, "Takataka")
	//tx5 := Transactions.Tx("Yaki", "Momo", 10, "Wow")
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
	blockchain.AddTransaction(tx1)
	blockchain.AddTransaction(tx2)
	blockchain.AddTransaction(tx3)
	blockchain.AddTransaction(tx4)
	//blockchain.AddTransaction(tx5)
	//blockchain.AddTransaction(tx6)
	//blockchain.AddTransaction(tx7)
	//blockchain.AddTransaction(tx8)
	//
	blockchain.MineNextBlock()
	//blockchain.MineNextBlock()
	//blockchain.MineNextBlock()
	//blockchain.MineNextBlock()
	//blockchain.MineNextBlock()
	//blockchain.MineNextBlock()
	//
	//fmt.Println(&blockchain)

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

	cli := CommandInterface.CLI{blockchain}
	cli.Run()

}
