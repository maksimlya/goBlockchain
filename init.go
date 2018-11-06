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

	c := Blocks.MineBlock(5, b.GetHash())
	d := Blocks.MineBlock(5, c.GetHash())
	e := Blocks.MineBlock(5, d.GetHash())
	f := Blocks.MineBlock(5, e.GetHash())

	f.PrintHash()

	enc := sha256.New()
	enc.Write([]byte("abcd"))
	j := hex.EncodeToString(enc.Sum(nil))

	a := Security.GenerateKey(j)

	Security.Encrypt("gogo", a)

	fmt.Println(a)

}
