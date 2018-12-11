package core

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
	"time"
)

//database name
const dbName = "naivechain.db"
const tableName = "blocks"

//blockchain structure
type Blockchain struct {
	//Blocks []*Block //store ordered blocks, no need if there is a db

	Tail []byte //newest block hash
	DB   *bolt.DB
}

//check if the db exist
func DBExist() bool {
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		return false
	}
	return true
}

//generate a blockchain with a genesis block
func CreateBlockchainWithAGenesisBlock(data string) {
	//check if the db exist
	if DBExist() {
		fmt.Println("genesis block already exists!")
		os.Exit(1)
	}

	fmt.Println("creating the genesis block")
	//initialize the database
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {

		bucket, err := tx.CreateBucket([]byte(tableName))
		if err != nil {
			log.Panic(err)
		}

		if bucket != nil {
			//generate a genesis block
			genesisBlock := CreateGenesisBlock(data)
			//save genesisBlock to the table
			err = bucket.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}

			//save the newest block hash
			err = bucket.Put([]byte("tail"), genesisBlock.Hash)
			if err != nil {
				log.Panic(err)
			}
		}
		return nil
	})
}

//add a block to the blockchain
func (blc *Blockchain) AddBlockToBlockchain(data string) {
	//save the block to the database
	err := blc.DB.Update(func(tx *bolt.Tx) error {

		//get the table from the database
		table := tx.Bucket([]byte(tableName))
		//get the newest block
		lastBlock := Deserialize(table.Get(blc.Tail))
		//create a new block
		block := NewBlock(data, lastBlock.Height+1, lastBlock.Hash)
		//save the block to the database
		err := table.Put(block.Hash, block.Serialize())
		if err != nil {
			log.Panic(err)
		}
		//update the tail in the database
		err = table.Put([]byte("tail"), block.Hash)
		if err != nil {
			log.Panic(err)
		}
		//update the tail in the blockchain
		blc.Tail = block.Hash

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

}

//return a point to the iterator
func (blc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{blc.Tail, blc.DB}
}

//print the info of a blockchain by iterating visit
func (blc *Blockchain) PrintChain() {
	//initialize a blockchain iterator
	blcIter := blc.Iterator()
	for {
		block := blcIter.Next()
		if block == nil {
			return
		}
		fmt.Printf("Height: %d, PredHash:%x, Data:%s, TimeStamp:%s, Hash:%x, Nonce:%d\n",
			block.Height, block.PrevHash, block.Data, time.Unix(block.Timestamp, 0).Format("2006-01-02 03:04:05 PM"), block.Hash, block.Nonce)
	}

}

//get a blockchain object
func BlockchainObject() *Blockchain {
	//open the database
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	var tailHash []byte
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(tableName))

		if b != nil {
			tailHash = b.Get([]byte("tail"))

		} else {
			log.Fatal("no blockchain")
		}
		return nil

	})

	return &Blockchain{tailHash, db}
}
