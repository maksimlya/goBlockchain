package p2p

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

var instance *NetworkController
var once sync.Once

type NetworkController struct {
	NodeAddress string
	KnownNodes  []string
}

func (nc *NetworkController) GetNodeAddress() string {
	return nc.NodeAddress
}
func (nc *NetworkController) GetKnownNodes() []string {
	return nc.KnownNodes
}

func (nc *NetworkController) AppendKnownNode(address []string) {
	nc.KnownNodes = append(nc.KnownNodes, address...)
}
func (nc *NetworkController) SetNodeAddress(address string) {
	nc.NodeAddress = address
}

func GetInstance() *NetworkController {
	once.Do(func() {
		instance = &NetworkController{}
	})
	return instance
}

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
	data := Block{nc.GetNodeAddress(), block}
	payload := GobEncode(data)
	request := append(CmdToBytes("block"), payload...)

	nc.SendData(addr, request)
}

func (nc *NetworkController) BroadcastBlock(block []byte) {
	data := Block{nc.GetNodeAddress(), block}
	payload := GobEncode(data)
	request := append(CmdToBytes("block"), payload...)

	for _, node := range nc.GetKnownNodes() {
		if node != nc.GetNodeAddress() {
			nc.SendData(node, request)
		}
	}

}

func (nc *NetworkController) SendInv(address, kind string, items [][]byte) {
	inventory := Inv{nc.GetNodeAddress(), kind, items}
	payload := GobEncode(inventory)
	request := append(CmdToBytes("inv"), payload...)

	nc.SendData(address, request)
}

func (nc *NetworkController) SendTx(addr string, tx []byte) {
	data := Tx{nc.GetNodeAddress(), tx}
	payload := GobEncode(data)
	request := append(CmdToBytes("tx"), payload...)

	nc.SendData(addr, request)
}

func (nc *NetworkController) SendVersion(addr string, bestHeight int) {
	payload := GobEncode(Version{version, bestHeight, nc.GetNodeAddress()})

	request := append(CmdToBytes("version"), payload...)

	nc.SendData(addr, request)
}

func (nc *NetworkController) SendGetBlocks(address string) {
	payload := GobEncode(GetBlocks{nc.GetNodeAddress()})
	request := append(CmdToBytes("getBlocks"), payload...)

	nc.SendData(address, request)
}

func (nc *NetworkController) SendGetData(address, kind string, id []byte) {
	payload := GobEncode(GetData{nc.GetNodeAddress(), kind, id})
	request := append(CmdToBytes("getData"), payload...)

	nc.SendData(address, request)
}

func (nc *NetworkController) SendAddr(address string) {
	nodes := Addr{nc.GetKnownNodes()}
	nodes.AddrList = append(nodes.AddrList, nc.GetNodeAddress())
	payload := GobEncode(nodes)
	request := append(CmdToBytes("addr"), payload...)

	nc.SendData(address, request)
}

func (nc *NetworkController) SendData(addr string, data []byte) {
	conn, err := net.Dial(protocol, addr)

	if err != nil {
		fmt.Printf("%s is not available\n", addr)
		var updatedNodes []string

		for _, node := range nc.GetKnownNodes() {
			if node != addr {
				updatedNodes = append(updatedNodes, node)
			}
		}
		nc.AppendKnownNode(updatedNodes)
		return
	}

	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
}
