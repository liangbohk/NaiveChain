package core

import (
	"bytes"
	"io"
	"log"
	"net"
)

//send command to other nodes

func sendData(to string, data []byte) {
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

func sendVersion(to string, blc *Blockchain) {
	baseHeight, _ := blc.GetBlockchainHeight()
	//baseHeight := 1
	version := &Version{NODE_VERSION, baseHeight, nodeAddress}
	thisBytes := version.Serialize()
	data := append(command2Bytes(VERSION_COMMAND), thisBytes...)
	sendData(to, data)

}

func sendGetBlocks(to string) {
	getBlocks := &GetBlocks{nodeAddress}
	thisBytes := getBlocks.Serialize()
	data := append(command2Bytes(GETBLOCKS_COMMAND), thisBytes...)
	sendData(to, data)
}

func sendInv(to string, command string, hashes [][]byte) {
	inv := &Inv{nodeAddress, BLOCK_TYPE, hashes}
	thisBytes := inv.Serialize()
	data := append(command2Bytes(command), thisBytes...)
	sendData(to, data)
}
