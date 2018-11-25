package handlers

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"goBlockchain/blockchain"
	"goBlockchain/p2p"
	"goBlockchain/transactions"
	"io/ioutil"
	"log"
	"net"
)

var nc *p2p.NetworkController
var blocksInTransit = [][]byte{}

const (
	protocol      = "tcp"
	version       = 1
	commandLength = 12
)

func HandleConnection(conn net.Conn, chain *blockchain.Blockchain) {
	req, err := ioutil.ReadAll(conn)
	defer conn.Close()

	if err != nil {
		log.Panic(err)
	}
	command := p2p.BytesToCmd(req[:commandLength])
	fmt.Printf("Received %s command \n", command)

	switch command {
	case "version":
		HandleVersion(req, chain)
	case "addr":
		HandleAddr(req)
	case "block":
		HandleBlock(req, chain)
	case "inv":
		HandleInv(req, chain)
	case "getBlocks":
		HandleGetBlocks(req, chain)
	case "getData":
		HandleGetData(req, chain)
	case "tx":
		HandleTx(req, chain)
	default:
		fmt.Println("Unknown command")
	}
}

func HandleGetBlocks(request []byte, chain *blockchain.Blockchain) {
	var buff bytes.Buffer
	var payload p2p.GetBlocks

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blocks := chain.GetBlockHashes()

	nc.SendInv(payload.AddrFrom, "block", blocks)
}

func HandleGetData(request []byte, chain *blockchain.Blockchain) {
	var buff bytes.Buffer
	var payload p2p.GetData

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	if payload.Type == "block" {

		block := chain.GetBlockByHash(string(payload.ID[:]))

		nc.SendBlock(payload.AddrFrom, block.Serialize())
	}
	if payload.Type == "tx" {
		txID := hex.EncodeToString(payload.ID)
		for _, t := range chain.GetPendingTransactions() {
			if t.Hash == txID {
				tx := t
				nc.SendTx(payload.AddrFrom, tx.Serialize())
			}
		}

	}
}
func HandleVersion(request []byte, chain *blockchain.Blockchain) {
	var buff bytes.Buffer
	var payload p2p.Version

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	bestHeight := chain.GetBlocksAmount()
	otherHeight := payload.BestHeight

	if bestHeight < otherHeight {
		nc.SendGetBlocks(payload.AddrFrom)
	} else if bestHeight > otherHeight {
		nc.SendVersion(payload.AddrFrom, chain.GetBlocksAmount())
	} else {
		fmt.Printf("Current Blockchain is up-to-date with %s peer", payload.AddrFrom)
	}

	if !NodeIsKnown(payload.AddrFrom) {
		nc.KnownNodes = append(nc.KnownNodes, payload.AddrFrom)
	}

}

func NodeIsKnown(addr string) bool {
	for _, node := range nc.KnownNodes {
		if node == addr {
			return true
		}
	}
	return false
}

func HandleTx(request []byte, chain *blockchain.Blockchain) {
	var buff bytes.Buffer
	var payload p2p.Tx

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	txData := payload.Transaction
	tx := transactions.DeserializeTransaction(txData)
	chain.AddTransaction(*tx, "")

	fmt.Printf("%s, %d ATTENTION TRANSACTION WILL NOT BE ADDED SINCE WE'RE NOT YET SENDING SIGNATURE!!!", nc.GetNodeAddress(), len(chain.GetPendingTransactions()))

	if nc.GetNodeAddress() == nc.KnownNodes[0] {
		for _, node := range nc.GetKnownNodes() {
			if node != nc.GetNodeAddress() && node != payload.AddrFrom {
				nc.SendInv(node, "tx", [][]byte{[]byte(tx.GetHash())})
			}
		}
	}
}

func HandleInv(request []byte, chain *blockchain.Blockchain) {
	var buff bytes.Buffer
	var payload p2p.Inv

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Received inventory with %d %s\n", len(payload.Items), payload.Type)

	if payload.Type == "block" {
		blocksInTransit = payload.Items

		blockHash := payload.Items[0]
		nc.SendGetData(payload.AddrFrom, "block", blockHash)

		newInTransit := [][]byte{}

		for _, b := range blocksInTransit {
			if bytes.Compare(b, blockHash) != 0 {
				newInTransit = append(newInTransit, b)
			}
		}
		blocksInTransit = newInTransit
	}
	if payload.Type == "tx" {
		txID := payload.Items[0]

		if chain.GetPendingTransactionByHash(string(txID[:])).GetHash() == "" {
			nc.SendGetData(payload.AddrFrom, "tx", txID)
		}

	}
}

func HandleAddr(request []byte) {
	var buff bytes.Buffer
	var payload p2p.Addr
	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	nc.AppendKnownNode(payload.AddrList)
	fmt.Printf("There are %d known nodes\n", len(nc.GetKnownNodes()))
	RequestBlocks()
}

func HandleBlock(request []byte, chain *blockchain.Blockchain) {
	var buff bytes.Buffer
	var payload p2p.Block

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blockData := payload.Block
	block := blockchain.DeserializeBlock(blockData)

	fmt.Println("Received a new Block!")
	chain.AddBlock(block)

	fmt.Printf("Added block %s\n", block.GetHash())

	if len(blocksInTransit) > 0 {
		blockHash := blocksInTransit[0]
		nc.SendGetData(payload.AddrFrom, "block", blockHash)

		blocksInTransit = blocksInTransit[1:]
	} else {
		checkChain := chain.ValidateChain()
		fmt.Println(checkChain)
	}
}

func RequestBlocks() {
	for _, node := range nc.GetKnownNodes() {
		nc.SendGetBlocks(node)
	}
}

func StartServer(nodeID string) {
	nc = p2p.GetInstance()
	nc.AppendKnownNode([]string{"192.168.2.101:3000"})
	nc.SetNodeAddress("192.168.2.101:3000")
	//nc.SetNodeAddress("192.168.2.110:3000")

	ln, err := net.Listen(protocol, nc.GetNodeAddress())
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	chain := blockchain.GetInstance()
	defer chain.CloseDB()

	if nc.GetNodeAddress() != nc.GetKnownNodes()[0] {
		nc.SendVersion(nc.GetKnownNodes()[0], chain.GetBlocksAmount())
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		go HandleConnection(conn, chain)
	}
}
