package core

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
)

const PROTOCOL = "tcp"
const COMMANDLENGTH = 12
const NODE_VERSION = 1

var knowNodes = []string{"localhost:3000"}

func StartServer(nodeID string, mineAddress string) {

	//ip address
	nodeAddress := fmt.Sprintf("localhost:%s", nodeID)

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
		sendVersion(nodeAddress, knowNodes[0], blc)
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
		//read data from client
		request, err := ioutil.ReadAll(conn)
		if err != nil {
			log.Panic(err)
		}
		fmt.Printf("msg received: %s\n", request)
	}
}

func sendData(to string, from string, data []byte) {
	conn, err := net.Dial("tcp", to)
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	//send data
	_, err = io.Copy(conn, bytes.NewBuffer(data))
	if err != nil {
		log.Panic(err)
	}
}

func sendVersion(from string, to string, blc *Blockchain) {
	//baseHeight := blc.getBaseHeight()
	baseHeight := 1
	version := &Version{NODE_VERSION, baseHeight, from}
	bytes := version.Serialize()
	data := append(command2Bytes("version"), bytes...)
	sendData(from, to, data)

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
