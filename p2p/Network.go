package p2p

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"goBlockchain/blockchain"
	"goBlockchain/transactions"
	"gopkg.in/vrecan/death.v3"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"runtime"
	"syscall"
)

type Addr struct {
	AddrList []string
}

type Block struct {
	AddrFrom string
	Block    []byte
}
type GetBlocks struct {
	AddrFrom string
}
type GetData struct {
	AddrFrom string
	Type     string
	ID       []byte
}
type Inv struct {
	AddrFrom string
	Type     string
	Items    [][]byte
}
type Tx struct {
	AddrFrom    string
	Transaction []byte
}
type Version struct {
	Version    int
	BestHeight int
	AddrFrom   string
}

const (
	protocol      = "tcp"
	version       = 1
	commandLength = 12
)

var (
	nodeAdress      string
	minerAddress    string
	KnownNodes      = []string{"192.168.2.101:3000"}
	blocksInTransit = [][]byte{}
	memoryPool      = make(map[string]transactions.Transaction)
)

func CmdToBytes(cmd string) []byte {
	var bytes [commandLength]byte

	for i, c := range cmd {
		bytes[i] = byte(c)
	}

	return bytes[:]
}

func BytesToCmd(bytes []byte) string {
	var cmd []byte

	for _, b := range bytes {
		if b != 0 {
			cmd = append(cmd, b)
		}
	}
	return fmt.Sprintf("%s", cmd)
}

func GobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

func CloseDB(chain *blockchain.Blockchain) {
	d := death.NewDeath(syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	d.WaitForDeathWithFunc(func() {
		defer os.Exit(1)
		defer runtime.Goexit()
		chain.CloseDB()
	})
}

func SendBlock(addr string, b *blockchain.Block) {
	data := Block{nodeAdress, b.Serialize()}
	payload := GobEncode(data)
	request := append(CmdToBytes("block"), payload...)

	SendData(addr, request)
}

func SendInv(address, kind string, items [][]byte) {
	inventory := Inv{nodeAdress, kind, items}
	payload := GobEncode(inventory)
	request := append(CmdToBytes("inv"), payload...)

	SendData(address, request)
}

func SendTx(addr string, tnx *transactions.Transaction) {
	data := Tx{nodeAdress, tnx.Serialize()}
	payload := GobEncode(data)
	request := append(CmdToBytes("tx"), payload...)

	SendData(addr, request)
}

func SendVersion(addr string, chain *blockchain.Blockchain) {
	bestHeight := chain.GetLastBlock().GetId() + 1
	payload := GobEncode(Version{version, bestHeight, nodeAdress})

	request := append(CmdToBytes("version"), payload...)

	SendData(addr, request)
}

func SendGetBlocks(address string) {
	payload := GobEncode(GetBlocks{nodeAdress})
	request := append(CmdToBytes("getBlocks"), payload...)

	SendData(address, request)
}

func SendGetData(address, kind string, id []byte) {
	payload := GobEncode(GetData{nodeAdress, kind, id})
	request := append(CmdToBytes("getData"), payload...)

	SendData(address, request)
}

func SendAddr(address string) {
	nodes := Addr{KnownNodes}
	nodes.AddrList = append(nodes.AddrList, nodeAdress)
	payload := GobEncode(nodes)
	request := append(CmdToBytes("addr"), payload...)

	SendData(address, request)
}

func SendData(addr string, data []byte) {
	conn, err := net.Dial(protocol, addr)

	if err != nil {
		fmt.Printf("%s is not available\n", addr)
		var updatedNodes []string

		for _, node := range KnownNodes {
			if node != addr {
				updatedNodes = append(updatedNodes, node)
			}
		}
		KnownNodes = updatedNodes
		return
	}

	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
}

func HandleConnection(conn net.Conn, chain *blockchain.Blockchain) {
	req, err := ioutil.ReadAll(conn)
	defer conn.Close()

	if err != nil {
		log.Panic(err)
	}
	command := BytesToCmd(req[:commandLength])
	fmt.Printf("Received %s command \n", command)

	switch command {
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
	case "version":
		HandleVersion(req, chain)
	default:
		fmt.Println("Unknown command")
	}
}

func HandleGetBlocks(request []byte, chain *blockchain.Blockchain) {
	var buff bytes.Buffer
	var payload GetBlocks

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blocks := chain.GetBlockHashes()

	SendInv(payload.AddrFrom, "block", blocks)
}

func HandleGetData(request []byte, chain *blockchain.Blockchain) {
	var buff bytes.Buffer
	var payload GetData

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	if payload.Type == "block" {
		blockId := binary.BigEndian.Uint64(payload.ID)
		block := chain.GetBlockById(int(blockId))

		SendBlock(payload.AddrFrom, block)
	}
	if payload.Type == "tx" {
		txID := hex.EncodeToString(payload.ID)
		for _, t := range chain.GetPendingTransactions() {
			if t.Hash == txID {
				tx := t
				SendTx(payload.AddrFrom, &tx)
			}
		}

	}
}
func HandleVersion(request []byte, chain *blockchain.Blockchain) {
	var buff bytes.Buffer
	var payload Version

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	bestHeight := chain.GetLastBlock().GetId()
	otherHeight := payload.BestHeight

	if bestHeight < otherHeight {
		SendGetBlocks(payload.AddrFrom)
	} else if bestHeight > otherHeight {
		SendVersion(payload.AddrFrom, chain)
	}

	if !NodeIsKnown(payload.AddrFrom) {
		KnownNodes = append(KnownNodes, payload.AddrFrom)
	}

}

func NodeIsKnown(addr string) bool {
	for _, node := range KnownNodes {
		if node == addr {
			return true
		}
	}
	return false
}

func HandleTx(request []byte, chain *blockchain.Blockchain) {
	var buff bytes.Buffer
	var payload Tx

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	txData := payload.Transaction
	tx := transactions.DeserializeTransaction(txData)
	memoryPool[hex.EncodeToString([]byte(tx.GetHash()))] = *tx

	fmt.Printf("%s, %d ATTENTION TRANSACTION WILL NOT BE ADDED SINCE WE'RE NOT YET SENDING SIGNATURE!!!", nodeAdress, len(memoryPool))

	if nodeAdress == KnownNodes[0] {
		for _, node := range KnownNodes {
			if node != nodeAdress && node != payload.AddrFrom {
				SendInv(node, "tx", [][]byte{[]byte(tx.GetHash())})
			}
		}
	} else {
		if len(memoryPool) >= 2 && len(minerAddress) > 0 {
			//MineTx(chain)
			fmt.Printf("Here it should MineTx(chain)")
		}
	}
}

func HandleInv(request []byte, chain *blockchain.Blockchain) {
	var buff bytes.Buffer
	var payload Inv

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
		SendGetData(payload.AddrFrom, "block", blockHash)

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

		if memoryPool[hex.EncodeToString(txID)].GetHash() == "" {
			SendGetData(payload.AddrFrom, "tx", txID)
		}

	}
}

func HandleAddr(request []byte) {
	var buff bytes.Buffer
	var payload Addr
	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	KnownNodes = append(KnownNodes, payload.AddrList...)
	fmt.Printf("There are %d known nodes\n", len(KnownNodes))
	RequestBlocks()
}

func HandleBlock(request []byte, chain *blockchain.Blockchain) {
	var buff bytes.Buffer
	var payload Block

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
		SendGetData(payload.AddrFrom, "block", blockHash)

		blocksInTransit = blocksInTransit[1:]
	} else {
		checkChain := chain.ValidateChain()
		fmt.Println(checkChain)
	}
}

func RequestBlocks() {
	for _, node := range KnownNodes {
		SendGetBlocks(node)
	}
}

func StartServer(nodeID, minerAddress string) {
	nodeAdress = fmt.Sprintf("192.168.2.101:%s", nodeID)
	minerAddress = minerAddress

	ln, err := net.Listen(protocol, nodeAdress)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	chain := blockchain.GetInstance()
	defer chain.CloseDB()
	go CloseDB(chain)

	if nodeAdress != KnownNodes[0] {
		SendVersion(KnownNodes[0], chain)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		go HandleConnection(conn, chain)
	}
}
