package core

import (
	"bytes"
	"crypto/sha256"
	"strconv"
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
}

func (block *Block) SetHash() {
	//joint the componets in block struct, transcode the component to []byte form if it is not
	//timestamp
	timeStampBytes := []byte(strconv.FormatInt(block.Timestamp, 2))

	//height
	heightBytes := Int2ByteArray(block.Height)

	//join the componets
	blockBytes := bytes.Join([][]byte{heightBytes, block.PrevHash, block.Hash, block.Data, timeStampBytes}, []byte{})

	//compute hash
	hash := sha256.Sum256(blockBytes)

	block.Hash = hash[:]

}

//create new block
func NewBlock(data string, height int64, prevHash []byte) *Block {
	block := &Block{height, prevHash, nil, []byte(data), time.Now().Unix()}

	//set block hash
	block.SetHash()

	return block
}

//generate the genesis block
func CreateGenesisBlock(data string) *Block {
	return NewBlock(data, 0, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}
