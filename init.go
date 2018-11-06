package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"goBlockchain/Blocks"
	"goBlockchain/Security"
	"goBlockchain/Transactions"
)

func main() {

	b := Blocks.MineGenesisBlock()

	b.PrintIdx()
	b.PrintHash()
	b.PrintTime()

	t1 := Transactions.Tx("yaki", "max", 10, "shit")
	fmt.Println(t1)
	t2 := Transactions.Tx("max", "yaki", 5, "shit")
	fmt.Println(t2)

	enc := sha256.New()
	enc.Write([]byte("abcd"))
	j := hex.EncodeToString(enc.Sum(nil))

	a := Security.GenerateKey(j)

	Security.Encrypt("gogo", a)

	fmt.Println(a)

}
