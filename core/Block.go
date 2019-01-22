package core

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

//the block structure
type Block struct {
	//block height
	Height int64
	//prev block hash
	PrevHash []byte
	//hash
	Hash []byte
	//transactions
	Txs []*Transaction
	//time stamp
	Timestamp int64
	//nouce attribute
	Nonce int64
}

func (block *Block) Transactions2Hash() []byte {
	//var txHashArr [][]byte
	//
	//for _, tx := range block.Txs {
	//	txHashArr = append(txHashArr, tx.TxHash)
	//}
	//txsHash := sha256.Sum256(bytes.Join(txHashArr, []byte{}))
	//
	//return txsHash[:]

	var data [][]byte
	for _, tx := range block.Txs {
		data = append(data, tx.TxHash)
	}
	merkleTree := NewMerkleTree(data)
	return merkleTree.RootNode.CheckData
}

//create new block
func NewBlock(txs []*Transaction, height int64, prevHash []byte) *Block {
	block := &Block{height, prevHash, nil, txs, time.Now().Unix(), 0}

	//POW
	pow := NewProofOfWork(block)
	//assume zero digits number is 6
	hash, nonce := pow.Run()
	block.Hash, block.Nonce = hash, nonce

	return block
}

//generate the genesis block
func CreateGenesisBlock(txs []*Transaction) *Block {
	return NewBlock(txs, 0, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
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
func DeserializeBlock(blockbytes []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(blockbytes))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}
