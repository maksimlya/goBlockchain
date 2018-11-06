package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"goBlockchain/Blocks"
	"goBlockchain/Security"
)

func main() {
	b := Blocks.MineGenesisBlock()

	b.PrintIdx()
	b.PrintHash()
	b.PrintTime()

	// c := Blocks.MineBlock(5,b.GetHash())

	// c.PrintHash()

	enc := sha256.New()
	enc.Write([]byte("abcd"))
	j := hex.EncodeToString(enc.Sum(nil))

	a := Security.GenerateKey(j)

	Security.Encrypt("gogo", a)

	fmt.Println(a)

}
