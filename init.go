package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"goBlockchain/p2p"
	"goBlockchain/security"
	"goBlockchain/transactions"
)

//CalculateHash hashes the values of a TestContent

func main() {

	hasher := sha256.New()
	hasher.Write([]byte("MainNodeId1"))
	nodeHash := hex.EncodeToString(hasher.Sum(nil))

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

	//tx1 := transactions.Tx("Yaki", "Tomer", 10, "Wow")
	//tx2 := transactions.Tx("Yaki", "Mas Hahnasa", 10000, "Arnona")
	//tx3 := transactions.Tx("Yaki", "Tomer", 5, "Arnona")
	//tx4 := transactions.Tx("Yaki", "Adi", 10, "Takataka")
	//tx5 := transactions.Tx("Tomer", "Momo", 10, "Wow")
	//tx6 := transactions.Tx("Yaki", "Popo", 10000, "Arnona")
	//tx7 := transactions.Tx("Yaki", "Zozo", 51, "Arnona")
	//tx8 := transactions.Tx("Yaki", "Koko", 10, "Takataka")
	//
	//var list []transactions.Transaction
	////
	//list = append(list, tx1)
	//list = append(list, tx2)
	//list = append(list, tx3)
	////list = append(list, tx4)
	////
	//merkle, _ := datastructures.NewTree(list)
	//
	////	merkle.Root.PrintHash()
	//
	////g := security.GenerateKey("795A433949D3340E7CBA7971DE1B428830C15D901B65303B3A65C0A45EE3F498")
	////fmt.Println(g)
	////
	////t := "daklfwklwklkdlcl asdca cascascac"
	////
	////tt := []byte(t)
	////var pp []int
	////
	////for i := 0; i < len(tt); i++ {
	////	pp = append(pp, security.Encrypt(strconv.Itoa(int(tt[i])), g))
	////}
	////
	////privKey := security.GeneratePrivKey("795A433949D3340E7CBA7971DE1B428830C15D901B65303B3A65C0A45EE3F498")
	////
	////for i := 0; i < len(tt); i++ {
	////	pp[i] = security.Encrypt(strconv.Itoa(int(pp[i])), privKey)
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
	//blockchain := blockchain.GetInstance()
	////////
	//blockchain.AddTransaction(tx1, sign1)
	//blockchain.AddTransaction(tx2, sign2)
	//blockchain.AddTransaction(tx3, sign3)
	//blockchain.AddTransaction(tx4, sign4)
	//////
	//blockchain.MineNextBlock()
	//
	//blockchain.MineNextBlock()
	//blockchain.MineNextBlock()
	//blockchain.MineNextBlock()
	//blockchain.MineNextBlock()
	//

	//bl := blockchain.GetLastBlock()

	//j := blockchain.GetSignature(tx3.Id)
	//	fmt.Println(j)
	//k := make(map[string]string)
	//json.Unmarshal(j,&k)
	//fmt.Println(k)
	//fmt.Println(blockchain)

	//for {
	//	fmt.Println("Hello")
	//	time.Sleep(time.Second)
	//}
	//pubKey := security.GenerateKey("A034B1566E979D3C5FE487BF3CF721FF3517570E1151DFC67D0329A54A48F9F8")
	//fmt.Println(pubKey)
	//signature := security.Sign("muhahsada", "A034B1566E979D3C5FE487BF3CF721FF3517570E1151DFC67D0329A54A48F9F8")
	//
	//fmt.Println(signature)
	//fmt.Println(security.VerifySignature(signature, "muhahsada", pubKey))

	//	cli := cli.CLI{blockchain}
	//	cli.Run()

	//fmt.Println(blockchain.GetBalanceForAddress("Tomer"))

	//k := merkle.GetProofElements(tx2)
	//
	//test := datastructures.VerifyContent(tx3, k, merkle.MerRoot)
	//
	//fmt.Println(test)
	//fmt.Println(blockchain.ValidateChain())
	//blockchain := blockchain.InitBlockchain()
	//fmt.Println(blockchain.GetBlockById(1).GetMerkleTree().PrintLevels())

	p2p.StartServer("3000", "")

	//go func() {
	//
	//}()
	//log.Fatal(webserver.Run())

	//
	//pubKey := "MTAxMTAxMDAwMTAxMDEwMTExAAAAAAAAAAAAAAAAAAAxMDExMTExMTAwMDAwMDAwMDAxAAAAAAAAAAAAAAAAAA=="
	//txId := "d200f24e285bdbc9e48448030aedbf9188b927908cf693e62e167048f1165722"
	//signature := "nfhy5bH5dbCZtWp0fXNAOpIUQjZhGkdd5kcQ5CMQE-I=,MTA5NzAyLDE2OTA0OSw3MzA0LDEwOTcwMiwzMzgyOSwxNDUxMjMsMzYxODkxLDM3NzM1NCwxNjY2MzcsMjQ5MDgyLDI3ODk3OCwyNDkwODIsMzIxMTAsMTY2NjM3LDQyODc1LDQyODc1LDMyMTEwLDMzODI5LDMyMTEwLDEwOTcwMiwyNDkwODIsMzM4MjksMTY5MDQ5LDE4OTE0LDE2OTA0OSwzNjE4OTEsMzIxMTAsMTQ3NDI3LDMyMTEwLDE0NTEyMywxMDk3MDIsNDI4NzUsMTQ3NDI3LDM1MDgxNywyNzg5NzgsNzMwNCwxNjkwNDksMTY5MDQ5LDMyMTEwLDE2NjYzNywzMjExMCwxNjkwNDksNzMwNCw3MzA0LDE0NTEyMywzNTA4MTcsMzc3MzU0LDkzMzkyLDE0NTEyMywzMzgyOSwzMzgyOSwzNzczNTQsNzMwNCwzNTA4MTcsOTMzOTIsMzUwODE3LDE0NzQyNywxODkxNCw0Mjg3NSwxNjY2MzcsMTg5MTQsNzMwNCwxODkxNCwzNTA4MTc="
	//
	//
	//check := security.VerifySignature(signature,txId,pubKey)
	//
	//fmt.Println(check)

}
