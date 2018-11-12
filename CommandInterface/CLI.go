package CommandInterface

import (
	"flag"
	"fmt"
	"goBlockchain/Blockchain"
	"os"
)

type CLI struct {
	Bc *Blockchain.Blockchain
}

func (cli *CLI) Run() {
	cli.validateArgs()

	mineBlock := flag.NewFlagSet("mineblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	switch os.Args[1] {
	case "addblock":
		err := mineBlock.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if mineBlock.Parsed() {
		cli.mineBlock()
	}

	if printChainCmd.Parsed() {
		cli.printChain()

	}
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  addblock -data BLOCK_DATA - add a block to the blockchain")
	fmt.Println("  printchain - print all the blocks of the blockchain")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) mineBlock() {
	cli.Bc.MineNextBlock()
	fmt.Println("Success!")
}

func (cli *CLI) printChain() {
	bci := cli.Bc.Iterator()

	for {
		block := bci.Next()

		fmt.Println(block)

		if block.GetPreviousHash() == "0" {
			break
		}
	}
}
