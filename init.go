package main

import (
	"goBlockchain/Blocks"
)

func main()  {
	b := Blocks.MineGenesisBlock()

	b.PrintIdx()
	b.PrintHash()
	b.PrintTime()

c := Blocks.MineBlock(5,b.GetHash())

c.PrintHash()




}
