package main

import (
	"fmt"
	"goBlockchain/Security"
	"strconv"
	"strings"
)

//CalculateHash hashes the values of a TestContent

func main() {

	//tx1 := Transactions.Tx("Yaki", "Tomer", 10, "Wow")
	//tx2 := Transactions.Tx("Yaki", "Mas Hahnasa", 10000, "Arnona")
	//tx3 := Transactions.Tx("Yaki", "Zona", 5, "Arnona")
	//tx4 := Transactions.Tx("Yaki", "Adi", 10, "Takataka")
	//
	//var list []Transactions.Transaction
	//
	//list = append(list, tx1)
	//list = append(list, tx2)
	//list = append(list, tx3)
	//list = append(list, tx4)
	//
	//merkle, _ := DataStructures.NewTree(list)
	//
	//fmt.Println(merkle.GetTransactionsWithTag("Arnona"))
	//
	//merkle.Root.PrintHash()

	g := Security.GenerateKey("795A433949D3340E7CBA7971DE1B428830C15D901B65303B3A65C0A45EE3F498")
	fmt.Println(g)

	t := "795A433949D3340E7CBA7971DE1B428830C15D901B65303B3A65C0A45EE3F498"

	tt := []byte(t)
	var pp []int

	for i := 0; i < len(tt); i++ {
		pp = append(pp, Security.Encrypt(strconv.Itoa(int(tt[i])), g))
	}

	privKey := Security.GeneratePrivKey("795A433949D3340E7CBA7971DE1B428830C15D901B65303B3A65C0A45EE3F498")

	for i := 0; i < len(tt); i++ {
		pp[i] = Security.Encrypt(strconv.Itoa(int(pp[i])), privKey)
	}

	str := []string{}

	for i := range pp {
		bytea := byte(pp[i])

		str = append(str, string(bytea))
	}

	res := strings.Join(str, "")
	fmt.Println(res)

}
