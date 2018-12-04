package core

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
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

//generate a blockchain with a genesis block
func CreateBlockchainWithAGenesisBlock() *Blockchain {
	//initialize the database
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	var tailHash []byte
	err = db.Update(func(tx *bolt.Tx) error {
		//get the bucket
		bucket := tx.Bucket([]byte(tableName))
		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte(tableName))
			if err != nil {
				log.Panic(err)
			}
		}

		if bucket != nil {
			//generate a genesis block
			genesisBlock := CreateGenesisBlock("Genesis block")
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
			tailHash = genesisBlock.Hash
		}
		return nil
	})

	return &Blockchain{tailHash, db}
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

//print the info of a blockchain by iterating visit
func (blc *Blockchain) PrintChain() {
	//fmt.Println("Start to open")
	//db, err := bolt.Open(dbName, 0600, nil)
	//fmt.Println("Open")
	//if err != nil {
	//	log.Fatal(err)
	//}

	err := blc.DB.View(func(tx *bolt.Tx) error {
		fmt.Println("View")
		bucket := tx.Bucket([]byte(tableName))
		if bucket != nil {
			var curHash = blc.Tail
			for {
				//in this loop, get current block
				block := Deserialize(bucket.Get(curHash))
				fmt.Printf("Height: %d, PredHash:%x, Data:%s, TimeStamp:%s, Hash:%x, Nonce:%d\n",
					block.Height, block.PrevHash, block.Data, time.Unix(block.Timestamp, 0).Format("2006-01-02 03:04:05 PM"), block.Hash, block.Nonce)

				var hashInt big.Int
				hashInt.SetBytes(block.PrevHash)
				if hashInt.Cmp(big.NewInt(0)) == 0 {
					break
				}
				curHash = block.PrevHash
			}
		}

		return nil

	})
	if err != nil {
		log.Panic(err)
	}
}
