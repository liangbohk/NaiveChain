package core

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
)

const PROTOCOL = "tcp"
const COMMANDLENGTH = 12
const NODE_VERSION = 1

const VERSION_COMMAND = "version"
const ADDR_COMMAND = "addr"
const BLOCK_COMMAND = "block"
const INV_COMMAND = "inv"
const GETBLOCKS_COMMAND = "getblocks"
const GETDATA_COMMAND = "getdata"
const TX_COMMAND = "tx"

const BLOCK_TYPE = "block"
const TX_TYPE = "tx"

var knowNodes = []string{"localhost:3000"}
var nodeAddress = knowNodes[0]

func StartServer(nodeID string, mineAddress string) {

	//ip address
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)

	ln, err := net.Listen(PROTOCOL, nodeAddress)

	if err != nil {
		log.Panic(err)
	}

	defer ln.Close()

	blc := BlockchainObject(nodeID)

	//master node 3000
	//wallet node 3001
	//miner node 3002
	if nodeAddress != knowNodes[0] {
		//send to master node to request data
		sendVersion(knowNodes[0], blc)
	}

	//go func(){
	//	time.Sleep(3*time.Second)
	//	//sendMessage(knowNodes[0],nodeAddress)
	//	sendMessage(knowNodes[0],"localhost:3001")
	//}()

	//receive msg from client
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleConnection(conn, blc)
	}
}

//convert command to bytes
func command2Bytes(command string) []byte {
	var bytes [COMMANDLENGTH]byte
	for index, commandByte := range command {
		bytes[index] = byte(commandByte)
	}
	return bytes[:]
}

//convert bytes to command
func bytes2Command(bytes []byte) string {
	var command []byte
	for _, commandByte := range bytes {
		if commandByte != 0x0 {
			command = append(command, commandByte)
		}
	}
	return string(command)
}

func handleConnection(conn net.Conn, blc *Blockchain) {
	//read data from client
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("msg received: %s\n", string(request))

	command := bytes2Command(request[:COMMANDLENGTH])
	switch command {
	case VERSION_COMMAND:
		handleVersion(request, blc)
	case ADDR_COMMAND:
		handleAddr(request, blc)
	case BLOCK_COMMAND:
		handleBlock(request, blc)
	case INV_COMMAND:
		handleInv(request, blc)
	case TX_COMMAND:
		handleTx(request, blc)
	case GETBLOCKS_COMMAND:
		handleGetBlocks(request, blc)
	case GETDATA_COMMAND:
		handleGetData(request, blc)
	default:
		fmt.Println("Unknown command!")
	}
}
