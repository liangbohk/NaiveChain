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

//send node version
func sendVersion(to string, blc *Blockchain) {
	baseHeight, _ := blc.GetBlockchainHeight()
	//baseHeight := 1
	version := &Version{NODE_VERSION, baseHeight, nodeAddress}
	thisBytes := version.Serialize()
	data := append(command2Bytes(VERSION_COMMAND), thisBytes...)
	sendData(to, data)

}

//request blocks
func sendGetBlocks(to string) {
	getBlocks := &GetBlocks{nodeAddress}
	thisBytes := getBlocks.Serialize()
	data := append(command2Bytes(GETBLOCKS_COMMAND), thisBytes...)
	sendData(to, data)
}

//send all block hash
func sendInv(to string, thisType string, hashes [][]byte) {
	inv := &Inv{nodeAddress, thisType, hashes}
	thisBytes := inv.Serialize()
	data := append(command2Bytes(INV_COMMAND), thisBytes...)
	sendData(to, data)
}

//request data
func sendGetData(to string, thisType string, blockHash []byte) {
	getData := &GetData{nodeAddress, thisType, blockHash}
	thisBytes := getData.Serialize()
	data := append(command2Bytes(GETDATA_COMMAND), thisBytes...)
	sendData(to, data)
}
