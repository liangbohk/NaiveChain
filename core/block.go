package core

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	//block height
	Height int64
	//prev block hash
	PrevHash []byte
	//hash
	Hash []byte
	//transactions
	Data []byte
	//time stamp
	Timestamp int64
	//nouce attribute
	Nonce int64
}

//create new block
func NewBlock(data string, height int64, prevHash []byte) *Block {
	block := &Block{height, prevHash, nil, []byte(data), time.Now().Unix(), 0}

	//POW
	pow := NewProofOfWork(block)
	//assume zero digits number is 6
	hash, nonce := pow.Run()
	block.Hash, block.Nonce = hash, nonce

	return block
}

//generate the genesis block
func CreateGenesisBlock(data string) *Block {
	return NewBlock(data, 0, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}

//serialize the block
func (block *Block) Serialize() []byte {
	var res bytes.Buffer
	//initialize an encoder
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}

	return res.Bytes()
}

//deserialize the block bytes
func Deserialize(blockbytes []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(blockbytes))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}
