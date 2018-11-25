package p2p

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
)

type NetworkController struct{}

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
	NodeAdress string
	KnownNodes = []string{"192.168.2.101:3000"}
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

func (nc *NetworkController) SendBlock(addr string, block []byte) {
	data := Block{NodeAdress, block}
	payload := GobEncode(data)
	request := append(CmdToBytes("block"), payload...)

	SendData(addr, request)
}

func (nc *NetworkController) BroadcastBlock(block []byte) {
	data := Block{NodeAdress, block}
	payload := GobEncode(data)
	request := append(CmdToBytes("block"), payload...)

	for _, node := range KnownNodes {
		SendData(node, request)
	}

}

func SendInv(address, kind string, items [][]byte) {
	inventory := Inv{NodeAdress, kind, items}
	payload := GobEncode(inventory)
	request := append(CmdToBytes("inv"), payload...)

	SendData(address, request)
}

func SendTx(addr string, tx []byte) {
	data := Tx{NodeAdress, tx}
	payload := GobEncode(data)
	request := append(CmdToBytes("tx"), payload...)

	SendData(addr, request)
}

func SendVersion(addr string, bestHeight int) {
	payload := GobEncode(Version{version, bestHeight, NodeAdress})

	request := append(CmdToBytes("version"), payload...)

	SendData(addr, request)
}

func SendGetBlocks(address string) {
	payload := GobEncode(GetBlocks{NodeAdress})
	request := append(CmdToBytes("getBlocks"), payload...)

	SendData(address, request)
}

func SendGetData(address, kind string, id []byte) {
	payload := GobEncode(GetData{NodeAdress, kind, id})
	request := append(CmdToBytes("getData"), payload...)

	SendData(address, request)
}

func SendAddr(address string) {
	nodes := Addr{KnownNodes}
	nodes.AddrList = append(nodes.AddrList, NodeAdress)
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
