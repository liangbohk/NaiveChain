package core

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
)

//const PROTOCOL = "tcp"

var knowNodes = []string{"localhost:3000"}

func StartServer(nodeID string, mineAddress string) {

	//ip address
	nodeAddress := fmt.Sprintf("localhost:%s", nodeID)

	ln, err := net.Listen("tcp", nodeAddress)

	if err != nil {
		log.Panic(err)
	}

	defer ln.Close()

	//master node 3000
	//wallet node 3001
	//miner node 3002
	if nodeAddress != knowNodes[0] {
		//send to master node to request data
		sendMessage(knowNodes[0], nodeAddress)
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

func sendMessage(to string, from string) {
	conn, err := net.Dial("tcp", to)
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	//send data
	_, err = io.Copy(conn, bytes.NewBuffer([]byte(from)))
	if err != nil {
		log.Panic(err)
	}
}
